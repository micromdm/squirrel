package s3

import (
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	client     *s3.S3
	bucketName string
}

func New(bucketName string) (*S3, error) {
	sess := session.Must(session.NewSession())
	client := s3.New(sess)
	return &S3{client: client, bucketName: bucketName}, nil
}

func (s *S3) HealthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := s.Healthz(); err != nil {
			http.Error(w, "health check failed", http.StatusInternalServerError)
			return
		}
	}
}

func (s *S3) Healthz() error {
	// TODO
	return nil
}

func (s *S3) FileHandler(w http.ResponseWriter, r *http.Request) {
	upath := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))
	results, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(upath),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer results.Body.Close()
	io.Copy(w, results.Body)
}
