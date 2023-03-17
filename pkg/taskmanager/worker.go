package taskmanager

import (
	"context"
	"fmt"
	"sync/atomic"

	"KillerFeature/ServerSide/internal/client_conn"
	ssh2 "KillerFeature/ServerSide/internal/client_conn/ssh"
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

	var cc client_conn.ClientConn
	var err error
	switch w.task.connectType {
	case ssh:
		sshBuilder := ssh2.NewSSHBuilder()

		cc, err = sshBuilder.CreateCC(&client_conn.Creds{
			IP:       w.task.IP,
			Login:    w.task.AuthData.Login,
			Password: w.task.AuthData.Password,
		})
		if err != nil {
			w.done <- struct{}{}
			// return errors.Wrap(err, service.ErrorCreateConnWrap)
		}
		defer cc.Close()
	}

	// TODO: GetOSCommandLib возвращает структуру с командами для конкретной ОС, вид ОС можно узнать через SSH

	// TODO: Close() когда задеплоится

	select {
	case <-ctx.Done():
	default:
	}

	data, _ := cc.Exec("ls")

	fmt.Print(string(data))

}
