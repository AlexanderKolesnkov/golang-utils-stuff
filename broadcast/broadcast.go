package broadcast

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Broadcast struct {
	ch          chan Message
	subscribers map[string]subscriber

	mu *sync.RWMutex
}

type subscriber struct {
	ctx       context.Context
	isClosing *atomic.Int32
	ch        chan Message
}

func NewBroadcast() *Broadcast {
	return &Broadcast{
		ch:          make(chan Message),
		subscribers: make(map[string]subscriber),

		mu: &sync.RWMutex{},
	}
}

func (b *Broadcast) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-b.ch:
			b.iterateSubscribers(message)
		case <-time.After(1 * time.Millisecond):
			break
		}
	}
}

func (b *Broadcast) iterateSubscribers(message Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for key, sub := range b.subscribers {
		if sub.isClosing.Load() == 1 {
			continue
		}

		select {
		case <-sub.ctx.Done():
			sub.isClosing.Store(1)
			go b.deleteSubscribe(key)
			break
		case sub.ch <- message:
			break
		case <-time.After(1 * time.Millisecond):
			break
		}
	}
}

func (b *Broadcast) Send(symbol, timeframe string, StartTime int64, confirm bool) {
	select {
	case b.ch <- Message{
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

func (b *Broadcast) Subscribe(ctx context.Context, key string) chan Message {
	b.mu.Lock()
	defer b.mu.Unlock()

	fmt.Println("Subscribe", key)
	ch := make(chan Message)

	b.subscribers[key] = subscriber{
		ctx:       ctx,
		isClosing: &atomic.Int32{},
		ch:        ch,
	}

	return ch
}

func (b *Broadcast) deleteSubscribe(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	fmt.Println("Delete", key)

	close(b.subscribers[key].ch)
	delete(b.subscribers, key)
}
