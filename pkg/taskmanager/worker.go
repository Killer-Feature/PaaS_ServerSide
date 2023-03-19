package taskmanager

import (
	"context"
	cconn "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn"
	ssh2 "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn/ssh"
	servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"
	"sync/atomic"
)

var (
	errCreateCC    = "error creating client connection"
	errCloseCC     = "error closing client connection"
	errProcessTask = "error processing task"
)

type workerManager struct {
	ctx          context.Context
	taskChan     chan *Task
	workers      map[ID]*worker
	workersCount uint32
	logger       *servlog.ServLogger
}

type worker struct {
	task   *Task
	done   chan struct{}
	logger *servlog.ServLogger
}

func newWorkerManager(ctx context.Context, taskChan chan *Task, logger *servlog.ServLogger) *workerManager {
	m := &workerManager{
		ctx:          ctx,
		taskChan:     taskChan,
		workers:      map[ID]*worker{},
		workersCount: 0,
		logger:       logger,
	}
	go m.run()
	return m
}

func (m *workerManager) createWorker(task *Task, logger *servlog.ServLogger) chan struct{} {
	doneCh := make(chan struct{})
	newWorker := &worker{
		task:   task,
		done:   doneCh,
		logger: logger,
	}
	m.workers[task.ID] = newWorker
	go newWorker.doWork(m.ctx)
	return doneCh
}

func (m *workerManager) run() {
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

func (w *worker) doWork(ctx context.Context) {
	defer w.task.callback(w.task.ID)

	switch w.task.connectType {
	case ssh:
		sshBuilder := ssh2.NewSSHBuilder()
		cc, err := sshBuilder.CreateCC(w.task.IP, w.task.AuthData.Login, w.task.AuthData.Password)
		if err != nil {
			w.logger.TaskError(uint64(w.task.ID), errCreateCC+": "+err.Error())
		}
		defer func(cc cconn.ClientConn) {
			err := cc.Close()
			if err != nil {
				w.logger.TaskError(uint64(w.task.ID), errCloseCC+": "+err.Error())
			}
		}(cc)

		err = w.task.ProcessTask(cc)

		if err != nil {
			w.logger.TaskError(uint64(w.task.ID), errProcessTask+": "+err.Error())
			w.done <- struct{}{}
		}
	}

	select {
	case <-ctx.Done():
	default:
	}
}
