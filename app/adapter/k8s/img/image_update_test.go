package img

import "testing"

func Test_getImageName(t *testing.T) {
	type args struct {
		imageName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				imageName: "gcr.io/test/test:latest",
			},
			want: "gcr.io/test/test",
		},
		{
			name: "test1",
			args: args{
				imageName: "gcr.io/test/test@sha256:aaaa",
			},
			want: "gcr.io/test/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getImageName(tt.args.imageName); got != tt.want {
				t.Errorf("getImageName() = %v, want %v", got, tt.want)
			}
		})
	}
}
