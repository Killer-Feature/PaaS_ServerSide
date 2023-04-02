package taskmanager

import (
	"context"
	servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"
	"sync/atomic"
)

var (
	errDeleteTask = "error deleting task by id from task manager"
	errFindTask   = "error finding task by id from task manager"
)

type ID uint64

type Manager[TKey comparable] struct {
	tasksByID  map[ID]*Task[TKey]
	tasksByKey map[TKey]*Task[TKey]

	logger *servlog.ServLogger

	// TODO: Change on chan with backlog
	workerChan    chan *Task[TKey]
	workerManager *workerManager[TKey]

	currentIndex uint64
}

func (m *Manager[_]) deleteTask(ID ID) {
	task, ok := m.tasksByID[ID]
	if !ok {
		m.logger.TaskError(uint64(ID), errFindTask)
	}
	delete(m.tasksByID, ID)
	_, ok = m.tasksByKey[task.Key]
	if !ok {
		m.logger.TaskError(uint64(ID), errDeleteTask)
	}
	delete(m.tasksByKey, task.Key)
}

type ConnectType int

type Task[TKey comparable] struct {
	Key         TKey
	ID          ID
	ProcessTask func(taskID ID) error
	callback    func(ID ID)
}

type AuthData struct {
	Login    string
	Password string
}

func NewTaskManager[TKey comparable](ctx context.Context, logger *servlog.ServLogger) *Manager[TKey] {
	taskChan := make(chan *Task[TKey])
	return &Manager[TKey]{
		tasksByID:     map[ID]*Task[TKey]{},
		tasksByKey:    map[TKey]*Task[TKey]{},
		workerChan:    taskChan,
		workerManager: newWorkerManager(ctx, taskChan, logger),
		currentIndex:  0,
		logger:        logger,
	}
}

func (m *Manager[TKey]) AddTask(processTask func(taskId ID) error, key TKey) (ID, error) {
	id := atomic.AddUint64(&m.currentIndex, 1)
	task := &Task[TKey]{
		Key:         key,
		ID:          ID(id),
		ProcessTask: processTask,
		callback:    m.deleteTask,
	}

	m.tasksByID[task.ID] = task
	m.tasksByKey[key] = task

	go func() {
		m.workerChan <- task
	}()

	return ID(id), nil
}
