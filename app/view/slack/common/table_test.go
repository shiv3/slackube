package common

import "testing"

func Test_codeBlock(t *testing.T) {
	type args struct {
		header []string
		in     [][]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				header: []string{"A", "B", "C"},
				in: [][]string{
					{"AAAA", "BBBBB", "CCCCCCC"},
					{"AAAA", "BBBBB", "CCCCCCC"},
					{"AAAA", "BBBBB", "CCCCCCC"},
				},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := codeBlock(tt.args.header, tt.args.in); got != tt.want {
				t.Errorf("codeBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
