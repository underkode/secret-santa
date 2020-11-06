package actions

import (
	"testing"
)

func Test_playOut(t *testing.T) {
	type args struct {
		participantCount int
	}
	tests := []struct {
		name string
		args args
	}{
		{"five participants", args{participantCount: 5}},
		{"four participants", args{participantCount: 4}},
		{"three participants", args{participantCount: 3}},
		{"two participants", args{participantCount: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := playOut(tt.args.participantCount)

			if len(got) != tt.args.participantCount {
				t.Errorf("playOut() = %v, len(playOut()) = %v, want %v", got, len(got), tt.args.participantCount)
			}

			if err != nil {
				t.Errorf("err = %v, want nil", err)
			}

			for index, it := range got {
				if index == it {
					t.Errorf("result[%d] == %d", index, it)
				}
			}
		})
	}
}
