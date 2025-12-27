package storage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type DiskStorage struct {
	file   *os.File
	mu     sync.Mutex
	offset int64
}

func NewDiskStorage(path string) (*DiskStorage, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	stat, _ := f.Stat()
	return &DiskStorage{
		file:   f,
		offset: stat.Size(),
	}, nil
}

func (s *DiskStorage) Write(key, value string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	currOffset := s.offset
	// 使用 | 作为 Key 和 Value 的分隔符，减少冲突
	data := fmt.Sprintf("%s|%s\n", key, value)
	n, err := s.file.WriteString(data)
	if err != nil {
		return 0, err
	}

	s.offset += int64(n)
	return currOffset, nil
}

func (s *DiskStorage) ReadAt(offset int64) (string, string, error) {
	f, err := os.Open(s.file.Name())
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	_, err = f.Seek(offset, 0)
	if err != nil {
		return "", "", err
	}

	// 使用 Scanner 读取一行，更稳健
	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 2)
		if len(parts) == 2 {
			return parts[0], parts[1], nil
		}
	}
	return "", "", fmt.Errorf("read failed at offset %d", offset)
}

func (s *DiskStorage) Close() {
	s.file.Sync()
	s.file.Close()
}
