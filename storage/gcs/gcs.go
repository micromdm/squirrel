package gcs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

type GCS struct {
	client *storage.Client
	bkt    *storage.BucketHandle
}

func New(bucketName, serviceAccountFile string) (*GCS, error) {
	ctx := context.Background()
	o := []option.ClientOption{
		option.WithServiceAccountFile(serviceAccountFile),
		option.WithScopes(storage.ScopeReadOnly),
	}
	client, err := storage.NewClient(ctx, o...)
	if err != nil {
		return nil, errors.Wrap(err, "creating google cloud storage client")
	}
	bkt := client.Bucket(bucketName)
	_, err = bkt.Attrs(ctx)
	if err != nil && err == storage.ErrBucketNotExist {
		return nil, fmt.Errorf("gcs bucket %s does not exist", bucketName)
	} else if err != nil {
		return nil, errors.Wrap(err, "getting bucket attributes")
	}
	return &GCS{client: client, bkt: bkt}, nil
}

func (g *GCS) HealthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := g.Healthz(); err != nil {
			http.Error(w, "health check failed", http.StatusInternalServerError)
			return
		}
	}
}

func (g *GCS) Healthz() error {
	_, err := g.bkt.Attrs(context.Background())
	return err
}

func (g *GCS) FileHandler(w http.ResponseWriter, r *http.Request) {
	upath := strings.TrimPrefix(r.URL.Path, "/")
	o := g.bkt.Object(path.Clean(upath))

	ctx := context.Background()
	f, err := o.NewReader(ctx)
	if err == storage.ErrObjectNotExist {
		errString := fmt.Sprintf("file %s not found: %s", upath, err)
		http.Error(w, errString, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}
