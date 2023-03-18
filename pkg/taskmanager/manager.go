package taskmanager

import (
	"context"
	"net/netip"
	"sync"
	"sync/atomic"

	cconn "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn"
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

	ProcessTask func(cconn.ClientConn) error

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

func (m *Manager) AddTask(processTask func(conn cconn.ClientConn) error, parsedIP netip.AddrPort, authData AuthData) (ID, error) {
	id := atomic.AddUint64(&m.currentIndex, 1)
	task := &Task{
		IP:          parsedIP,
		connectType: ssh,
		ID:          ID(id),
		AuthData:    authData,
		ProcessTask: processTask,
		callback:    m.deleteTask,
	}

	m.tasksByID[task.ID] = task
	m.tasksByIP[parsedIP] = task

	go func() {
		m.workerChan <- task
	}()

	return ID(id), nil
}
