package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Paxos 角色: Proposer, Acceptor
// 为了简化，我们只演示单轮 Paxos 的核心投票逻辑

type PrepareResponse struct {
	Accepted  bool
	PrevTerm  int
	PrevValue string
}

type AcceptResponse struct {
	Accepted bool
}

type Acceptor struct {
	ID           int
	PromisedTerm int
	AcceptedTerm int
	AcceptedVal  string
}

func (a *Acceptor) Prepare(term int) PrepareResponse {
	if term > a.PromisedTerm {
		fmt.Printf("[Acceptor-%d] 承诺 (Promise) Term: %d\n", a.ID, term)
		a.PromisedTerm = term
		return PrepareResponse{Accepted: true, PrevTerm: a.AcceptedTerm, PrevValue: a.AcceptedVal}
	}
	return PrepareResponse{Accepted: false}
}

func (a *Acceptor) Accept(term int, value string) AcceptResponse {
	if term >= a.PromisedTerm {
		fmt.Printf("[Acceptor-%d] 接受 (Accept) Term: %d, Value: %s\n", a.ID, term, value)
		a.PromisedTerm = term
		a.AcceptedTerm = term
		a.AcceptedVal = value
		return AcceptResponse{Accepted: true}
	}
	return AcceptResponse{Accepted: false}
}

func main() {
	fmt.Println("=== 简化 Paxos 算法演示 (Basic Paxos) ===")
	fmt.Println()
	fmt.Println("Paxos 是所有共识算法的祖师爷。")
	fmt.Println("本示例演示 Basic Paxos 的两个阶段: Prepare 和 Accept。")
	fmt.Println()

	rand.Seed(time.Now().UnixNano())

	// 准备 3 个 Acceptor
	acceptors := []*Acceptor{
		{ID: 1}, {ID: 2}, {ID: 3},
	}

	// 模拟一个 Proposer 提议值 "Data-A"
	proposalTerm := 1
	proposalValue := "Data-A"

	fmt.Printf("[Proposer] 第一阶段: 发送 Prepare(Term: %d)\n", proposalTerm)
	prepareCount := 0
	for _, a := range acceptors {
		resp := a.Prepare(proposalTerm)
		if resp.Accepted {
			prepareCount++
			// 在实际 Paxos 中，如果返回了之前接受过的值，Proposer 必须改用该值
			if resp.PrevValue != "" {
				proposalValue = resp.PrevValue
			}
		}
	}

	if prepareCount > len(acceptors)/2 {
		fmt.Printf("[Proposer] 获得大多数承诺 (%d/%d)，进入第二阶段: 发送 Accept(Term: %d, Value: %s)\n",
			prepareCount, len(acceptors), proposalTerm, proposalValue)

		acceptCount := 0
		for _, a := range acceptors {
			resp := a.Accept(proposalTerm, proposalValue)
			if resp.Accepted {
				acceptCount++
			}
		}

		if acceptCount > len(acceptors)/2 {
			fmt.Printf("[Proposer] 提议成功达成共识！最终值: %s\n", proposalValue)
		} else {
			fmt.Println("[Proposer] 提议失败 (未获得大多数 Accept)")
		}
	} else {
		fmt.Println("[Proposer] 提议失败 (未获得大多数 Prepare 承诺)")
	}

	fmt.Println()
	fmt.Println("=== 总结 ===")
	fmt.Println("1. Prepare 阶段: 预留一个 Term，防止旧的提议通过。")
	fmt.Println("2. Accept 阶段: 实际写入值。")
	fmt.Println("3. Paxos 的核心是基于数学证明的安全性，但由于其实现复杂，后来出现了 Raft 等更工程化的算法。")
}
