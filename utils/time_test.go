package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeSince(t *testing.T) {
	type args struct {
		start time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				start: time.Now(),
			},
			want: "",
		},
	}
	for _, tt := range tests {
		time.Sleep(time.Millisecond * 300)
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeSince(tt.args.start); got != tt.want {
				fmt.Println(len(got))
				t.Errorf("TimeSince() = [%v], want %v", got, tt.want)
			}
		})
	}
}
