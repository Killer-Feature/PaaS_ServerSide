package taskmanager

import (
	"context"
	servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"
	"sync/atomic"
)

var (
	errProcessTask = "error processing task"
)

type workerManager[TKey comparable] struct {
	ctx          context.Context
	taskChan     chan *Task[TKey]
	workers      map[ID]*worker[TKey]
	workersCount uint32
	logger       *servlog.ServLogger
}

type worker[TKey comparable] struct {
	task   *Task[TKey]
	done   chan struct{}
	logger *servlog.ServLogger
}

func newWorkerManager[TKey comparable](ctx context.Context, taskChan chan *Task[TKey], logger *servlog.ServLogger) *workerManager[TKey] {
	m := &workerManager[TKey]{
		ctx:          ctx,
		taskChan:     taskChan,
		workers:      map[ID]*worker[TKey]{},
		workersCount: 0,
		logger:       logger,
	}
	go m.run()
	return m
}

func (m *workerManager[TKey]) createWorker(task *Task[TKey], logger *servlog.ServLogger) chan struct{} {
	doneCh := make(chan struct{})
	newWorker := &worker[TKey]{
		task:   task,
		done:   doneCh,
		logger: logger,
	}
	m.workers[task.ID] = newWorker
	go newWorker.doWork(m.ctx)
	return doneCh
}

func (m *workerManager[TKey]) run() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case data := <-m.taskChan:
			atomic.AddUint32(&m.workersCount, 1)
			m.createWorker(data, m.logger)
		}
	}
}

func (w *worker[TKey]) doWork(ctx context.Context) {
	defer w.task.callback(w.task.ID)
	err := w.task.ProcessTask(w.task.ID)
	if err != nil {
		w.logger.TaskError(uint64(w.task.ID), errProcessTask+": "+err.Error())
		close(w.done)
	}

	select {
	case <-ctx.Done():
	default:
	}
}
