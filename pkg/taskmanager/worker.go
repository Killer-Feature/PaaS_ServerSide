package taskmanager

import (
	cconn "KillerFeature/ServerSide/pkg/client_conn"
	ssh2 "KillerFeature/ServerSide/pkg/client_conn/ssh"
	"context"
	"sync/atomic"
)

type workerManager struct {
	ctx          context.Context
	taskChan     chan *Task
	workers      map[ID]*worker
	workersCount uint32
}

type worker struct {
	task *Task
	done chan struct{}
}

func newWorkerManager(ctx context.Context, taskChan chan *Task) *workerManager {
	m := &workerManager{
		ctx:          ctx,
		taskChan:     taskChan,
		workers:      map[ID]*worker{},
		workersCount: 0,
	}
	go m.run()
	return m
}

func (m *workerManager) createWorker(task *Task) chan struct{} {
	doneCh := make(chan struct{})
	newWorker := &worker{
		task: task,
		done: doneCh,
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
			m.createWorker(data)
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
			//	TODO: log error & send to ws
		}
		defer func(cc cconn.ClientConn) {
			err := cc.Close()
			if err != nil {
				//	TODO: log error & send to ws
			}
		}(cc)

		err = w.task.ProcessTask(cc)

		if err != nil {
			//	TODO: log error & send to ws
			w.done <- struct{}{}
		}
	}

	select {
	case <-ctx.Done():
	default:
	}
}
