package taskmanager

import (
	"context"
	"net/netip"
	"sync"
	"sync/atomic"
)

type ID uint64

type Manager struct {
	tasksByID map[ID]*Task
	tasksByIP map[netip.AddrPort]*Task

	// TODO: Change on chan with backlog
	workerChan    chan *Task
	workerManager *workerManager

	currentIndex uint64
}

func (m *Manager) deleteTask(ID ID) {
	task, ok := m.tasksByID[ID]
	if !ok {
		//TODO: Add logging
	}
	delete(m.tasksByID, ID)
	_, ok = m.tasksByIP[task.IP]
	if !ok {
		//TODO: Add logging
	}
	delete(m.tasksByIP, task.IP)
}

type ConnectType int

const (
	ssh ConnectType = iota
	winrm
)

type Task struct {
	mu sync.RWMutex

	IP          netip.AddrPort
	ID          ID
	AuthData    AuthData
	connectType ConnectType

	callback func(ID ID)
}

type AuthData struct {
	Login    string
	Password string
}

func NewTaskManager(ctx context.Context) *Manager {
	taskChan := make(chan *Task)
	return &Manager{
		tasksByID:     map[ID]*Task{},
		tasksByIP:     map[netip.AddrPort]*Task{},
		workerChan:    taskChan,
		workerManager: newWorkerManager(ctx, taskChan),
		currentIndex:  0,
	}
}

func (m *Manager) AddTask(ip, login, password string) (ID, error) {
	parsedIP, err := netip.ParseAddrPort(ip)
	if err != nil {
		return 0, err
	}

	id := atomic.AddUint64(&m.currentIndex, 1)
	task := &Task{
		IP:          parsedIP,
		connectType: ssh,
		ID:          ID(id),
		AuthData: AuthData{
			Login:    login,
			Password: password,
		},
		callback: m.deleteTask,
	}

	m.tasksByID[task.ID] = task
	m.tasksByIP[parsedIP] = task

	go func() {
		m.workerChan <- task
	}()

	return ID(id), nil
}
