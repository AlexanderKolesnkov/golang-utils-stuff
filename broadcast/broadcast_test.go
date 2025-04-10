package broadcast

import (
	"context"
	"testing"
	"time"
)

func TestFullChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	broadcast := NewBroadcast()
	go broadcast.Listen(ctx)

	//channels := make([]chan Message, 0)

	broadcast.Subscribe(ctx, "first")
	broadcast.Subscribe(ctx, "second")
	broadcast.Subscribe(ctx, "third")

	go func() {
		for {
			broadcast.ch <- Message{}
			time.Sleep(time.Second)
		}
	}()

	select {}

	//for _, channel := range channels {
	//<-channel
	//}
}
