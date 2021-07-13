package gcr

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type GCRAdapterInterface interface {
	ListImages(ctx context.Context, fullname string, option *ListTagsOption) ([]Image, error)
}

type ListAdapterImpl struct {
}

type ListTagsOption struct {
	ProjectID     string
	getImageLimit int
}

// Image represents a image object of gks
type Image struct {
	Digest string   `json:"digest",yaml:"metadata,flow"`
	Tags   []string `json:"tags,omitempty",yaml:"spec,flow"`
}

func (c *ListAdapterImpl) ListImages(ctx context.Context, fullname string, option *ListTagsOption) ([]Image, error) {
	args := []string{
		"container", "images",
		"list-tags", fullname,
		"--format", "json",
	}

	if option != nil {
		args = append(args, "--limit", fmt.Sprintf("%d", option.getImageLimit))
	}

	cmd := exec.Command("gcloud", args...)
	log.Printf("%v", cmd.Args)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	images := []Image{}
	err = json.Unmarshal(out, &images)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func test() {
	//s, err := store.DefaultGCRCredStore()
	//if err != nil {
	//	panic(err)
	//}
	//
	//a := auth.GCRLoginAgent{
	//	AllowBrowser: true,
	//	In:           nil,
	//	Out:          nil,
	//	OpenBrowser:  nil,
	//}
	//t, err := a.PerformLogin()
	//if err != nil {
	//	panic(err)
	//}
	//s.SetGCRAuth(t)

}
