package gcr

import (
	"reflect"
	"testing"

	"github.com/google/go-containerregistry/pkg/v1/google"
)

func TestService_ListTags(t *testing.T) {
	type args struct {
		host      string
		repoName  string
		userAgent string
	}
	tests := []struct {
		name    string
		args    args
		want    *google.Tags
		wantErr bool
	}{
		{
			name: "",
			args: args{
				host:      "gcr.io",
				repoName:  "cloud-builders/docker",
				userAgent: "",
			},
			want: &google.Tags{
				Children:  nil,
				Manifests: nil,
				Name:      "",
				Tags:      nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := GcrAdapterImpl{host: "asia.gcr.io"}
			got, err := s.ListTags(tt.args.repoName, tt.args.userAgent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTags() got = %v, want %v", got, tt.want)
			}
		})
	}
}
