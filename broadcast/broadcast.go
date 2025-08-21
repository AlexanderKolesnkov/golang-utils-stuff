package broadcast

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
	"time"
)

type Broadcast struct {
	ch          chan Message
	subscribers map[string]subscriber

	mu     *sync.RWMutex
	logger *zap.Logger
}

type subscriber struct {
	ctx       context.Context
	isClosing *atomic.Int32
	ch        chan Message
}

func NewBroadcast(logger *zap.Logger) *Broadcast {
	return &Broadcast{
		ch:          make(chan Message),
		subscribers: make(map[string]subscriber),

		mu:     &sync.RWMutex{},
		logger: logger,
	}
}

func (broadcast *Broadcast) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-broadcast.ch:
			broadcast.iterateSubscribers(message)
		case <-time.After(1 * time.Millisecond):
			break
		}
	}
}

func (broadcast *Broadcast) iterateSubscribers(message Message) {
	broadcast.mu.RLock()
	defer broadcast.mu.RUnlock()

	for key, sub := range broadcast.subscribers {
		if sub.isClosing.Load() == 1 {
			continue
		}

		select {
		case <-sub.ctx.Done():
			sub.isClosing.Store(1)
			go broadcast.deleteSubscribe(key)
			break
		case sub.ch <- message:
			break
		case <-time.After(1 * time.Millisecond):
			break
		}
	}
}

func (broadcast *Broadcast) Send(symbol, timeframe string, StartTime int64, confirm bool) {
	select {
	case broadcast.ch <- Message{
		Symbol:    symbol,
		Timeframe: timeframe,
		StartTime: StartTime,
		Confirm:   confirm,
	}:
		break
	case <-time.After(1 * time.Millisecond):
		break
	}
}

func (broadcast *Broadcast) Subscribe(ctx context.Context, key string) chan Message {
	broadcast.mu.Lock()
	defer broadcast.mu.Unlock()

	broadcast.logger.Info("Subscribe", zap.String("key", key))
	ch := make(chan Message)

	broadcast.subscribers[key] = subscriber{
		ctx:       ctx,
		isClosing: &atomic.Int32{},
		ch:        ch,
	}

	return ch
}

func (broadcast *Broadcast) deleteSubscribe(key string) {
	broadcast.mu.Lock()
	defer broadcast.mu.Unlock()

	broadcast.logger.Info("Delete", zap.String("key", key))

	close(broadcast.subscribers[key].ch)
	delete(broadcast.subscribers, key)
}
