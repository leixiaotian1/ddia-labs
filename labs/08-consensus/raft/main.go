package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Raft 角色定义
const (
	Follower  = "Follower"
	Candidate = "Candidate"
	Leader    = "Leader"
)

// LogEntry 日志条目
type LogEntry struct {
	Term    int
	Command string
}

// RaftNode 代表集群中的一个节点
type RaftNode struct {
	ID    int
	mu    sync.Mutex
	state string

	currentTerm int
	votedFor    int
	log         []LogEntry

	// 集群信息
	peers []*RaftNode

	// 选举计时器
	electionTimeout  time.Duration
	lastHeartbeat    time.Time
	heartbeatTimeout time.Duration

	// 模拟停止
	stopped bool
}

func NewRaftNode(id int, heartbeatTimeout time.Duration) *RaftNode {
	// 随机选举超时时间 150ms - 300ms
	timeout := time.Duration(150+rand.Intn(150)) * time.Millisecond
	return &RaftNode{
		ID:               id,
		state:            Follower,
		currentTerm:      0,
		votedFor:         -1,
		log:              make([]LogEntry, 0),
		electionTimeout:  timeout,
		heartbeatTimeout: heartbeatTimeout,
		lastHeartbeat:    time.Now(),
	}
}

func (rn *RaftNode) SetPeers(peers []*RaftNode) {
	rn.peers = peers
}

func (rn *RaftNode) Start() {
	go rn.runLoop()
}

func (rn *RaftNode) Stop() {
	rn.mu.Lock()
	defer rn.mu.Unlock()
	rn.stopped = true
}

func (rn *RaftNode) runLoop() {
	for {
		rn.mu.Lock()
		if rn.stopped {
			rn.mu.Unlock()
			return
		}
		state := rn.state
		rn.mu.Unlock()

		switch state {
		case Follower:
			rn.runFollower()
		case Candidate:
			rn.runCandidate()
		case Leader:
			rn.runLeader()
		}
	}
}

func (rn *RaftNode) runFollower() {
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()

	for {
		<-tick.C
		rn.mu.Lock()
		if rn.stopped || rn.state != Follower {
			rn.mu.Unlock()
			return
		}

		// 检查选举超时
		if time.Since(rn.lastHeartbeat) >= rn.electionTimeout {
			fmt.Printf("[Node-%d] 选举超时，转变为 Candidate (Term: %d)\n", rn.ID, rn.currentTerm+1)
			rn.state = Candidate
			rn.mu.Unlock()
			return
		}
		rn.mu.Unlock()
	}
}

func (rn *RaftNode) runCandidate() {
	rn.mu.Lock()
	rn.currentTerm++
	rn.votedFor = rn.ID
	rn.lastHeartbeat = time.Now()
	term := rn.currentTerm
	rn.mu.Unlock()

	votes := 1 // 投给自己
	var voteMu sync.Mutex
	var wg sync.WaitGroup

	for _, peer := range rn.peers {
		if peer.ID == rn.ID {
			continue
		}
		wg.Add(1)
		go func(p *RaftNode) {
			defer wg.Done()
			if p.RequestVote(term, rn.ID) {
				voteMu.Lock()
				votes++
				voteMu.Unlock()
			}
		}(peer)
	}

	// 等待投票结果或超时
	done := make(chan bool, 1)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		rn.mu.Lock()
		defer rn.mu.Unlock()
		if rn.state != Candidate || rn.currentTerm != term {
			return
		}
		if votes > len(rn.peers)/2 {
			fmt.Printf("[Node-%d] 赢得选举成为 Leader (Term: %d, 选票: %d/%d)\n", rn.ID, term, votes, len(rn.peers))
			rn.state = Leader
		} else {
			fmt.Printf("[Node-%d] 选举失败 (选票: %d/%d)，返回 Follower\n", rn.ID, votes, len(rn.peers))
			rn.state = Follower
			rn.votedFor = -1
		}
	case <-time.After(rn.electionTimeout):
		// 重新开始选举循环
		return
	}
}

func (rn *RaftNode) runLeader() {
	tick := time.NewTicker(rn.heartbeatTimeout)
	defer tick.Stop()

	for {
		<-tick.C
		rn.mu.Lock()
		if rn.stopped || rn.state != Leader {
			rn.mu.Unlock()
			return
		}
		term := rn.currentTerm
		rn.mu.Unlock()

		// 发送心跳
		for _, peer := range rn.peers {
			if peer.ID == rn.ID {
				continue
			}
			go peer.AppendEntries(term, rn.ID)
		}
	}
}

// RPC: RequestVote
func (rn *RaftNode) RequestVote(term int, candidateID int) bool {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	if rn.stopped {
		return false
	}

	if term > rn.currentTerm {
		rn.currentTerm = term
		rn.state = Follower
		rn.votedFor = -1
	}

	if term == rn.currentTerm && (rn.votedFor == -1 || rn.votedFor == candidateID) {
		rn.votedFor = candidateID
		rn.lastHeartbeat = time.Now()
		return true
	}

	return false
}

// RPC: AppendEntries (心跳和日志同步)
func (rn *RaftNode) AppendEntries(term int, leaderID int) bool {
	rn.mu.Lock()
	defer rn.mu.Unlock()

	if rn.stopped {
		return false
	}

	if term >= rn.currentTerm {
		if term > rn.currentTerm || rn.state != Follower {
			rn.state = Follower
			rn.currentTerm = term
		}
		rn.lastHeartbeat = time.Now()
		return true
	}

	return false
}

func main() {
	fmt.Println("=== Raft 共识算法演示 (Leader Election) ===")
	fmt.Println()
	fmt.Println("特点: 强一致性、通过选主和日志复制实现。")
	fmt.Println("本示例演示 Raft 的核心选主流程及节点故障后的自动恢复。")
	fmt.Println()

	rand.Seed(time.Now().UnixNano())

	// 创建 3 个节点
	nodeCount := 3
	nodes := make([]*RaftNode, nodeCount)
	for i := 0; i < nodeCount; i++ {
		nodes[i] = NewRaftNode(i, 50*time.Millisecond)
	}

	for _, n := range nodes {
		n.SetPeers(nodes)
		n.Start()
	}

	// 观察一段时间的选主
	time.Sleep(1 * time.Second)

	fmt.Println("\n--- 模拟故障: 停止当前的 Leader ---")
	var leader *RaftNode
	for _, n := range nodes {
		n.mu.Lock()
		if n.state == Leader {
			leader = n
		}
		n.mu.Unlock()
	}

	if leader != nil {
		fmt.Printf("[System] 停止节点 Node-%d (Leader)\n", leader.ID)
		leader.Stop()
	}

	// 观察重新选主
	time.Sleep(1 * time.Second)

	fmt.Println("\n--- 恢复正常: 检查新的 Leader ---")
	for _, n := range nodes {
		n.mu.Lock()
		fmt.Printf("Node-%d 状态: %s (Term: %d)\n", n.ID, n.state, n.currentTerm)
		n.mu.Unlock()
	}

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. 节点通过心跳超时发现 Leader 故障。")
	fmt.Println("2. Candidate 开始选举并请求其他节点投票。")
	fmt.Println("3. 获得大多数选票的节点成为新的 Leader。")
	fmt.Println("4. 即使有节点挂掉，只要满足大多数节点 (N/2 + 1) 存活，系统就能继续工作。")

	// 停止所有节点
	for _, n := range nodes {
		n.Stop()
	}
}
