package dispatch

import (
	"context"
	"sync"
)

// Job — единица работы, которую мы будем ставить в очередь конкретного пользователя.
// Пока это просто функция. Позже мы будем “упаковывать” туда обработку DM-сообщения.
type Job func(ctx context.Context) error

// DMDispatcher — диспетчер задач для личных сообщений (DM), который гарантирует:
// 1) последовательность задач в рамках одного userID
// 2) параллельность между разными userID
type DMDispatcher struct {
	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup

	mu    sync.Mutex
	users map[int64]*userWorker
}

type userWorker struct {
	userID int64
	queue  chan Job
}

// NewDMDispatcher создаёт диспетчер и внутренний контекст жизни диспетчера.
// Этот ctx будет использоваться для остановки всех воркеров при shutdown.
func NewDMDispatcher(parent context.Context) *DMDispatcher {
	ctx, cancel := context.WithCancel(parent)
	return &DMDispatcher{
		ctx:    ctx,
		cancel: cancel,
		users:  make(map[int64]*userWorker),
	}
}

func (d *DMDispatcher) Enqueue(userID int64, job Job) {
	d.mu.Lock()
	uw, ok := d.users[userID]
	if !ok {
		uw = &userWorker{
			userID: userID,
			queue:  make(chan Job, 64), // буфер: чтобы не блокировать сразу
		}
		d.users[userID] = uw

		if !ok {
			uw = &userWorker{
				userID: userID,
				queue:  make(chan Job, 64),
			}
			d.users[userID] = uw

			d.wg.Add(1)
			go d.runUserWorker(uw)
		}
	}
	d.mu.Unlock()

	// Пока просто кладём job в очередь.
	// В следующем кусочке появится reader, который будет их выполнять.
	select {
	case <-d.ctx.Done():
		// Диспетчер уже остановлен — задачу не принимаем.
		return
	case uw.queue <- job:
		return
	}
}

func (d *DMDispatcher) runUserWorker(uw *userWorker) {
	defer d.wg.Done()

	for {
		select {
		case <-d.ctx.Done():
			return

		case job := <-uw.queue:
			// job может быть nil (если кто-то положил nil)
			if job == nil {
				continue
			}

			// Выполняем задачу. Ошибку пока никуда не отдаём — позже решим стратегию.
			_ = job(d.ctx)
		}
	}
}

// Stop — попросить диспетчер остановиться (через ctx) и дождаться завершения воркеров.
func (d *DMDispatcher) Stop() {
	d.cancel()
	d.wg.Wait()
}
