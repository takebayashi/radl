package gcs

import (
	"fmt"
	"github.com/takebayashi/radl"
	"golang.org/x/net/context"
	"google.golang.org/cloud/storage"
	"path"
)

type gcsStorage struct {
	bucket string
	ctx    context.Context
	client *storage.Client
}

func NewGCSStore(bucket string) (*gcsStorage, error) {
	s := &gcsStorage{}
	s.bucket = bucket
	s.ctx = context.Background()
	cli, err := storage.NewClient(s.ctx)
	if err != nil {
		return nil, err
	}
	s.client = cli
	return s, nil
}

func (fs *gcsStorage) Save(show radl.ShowPayload) error {
	err := fs.upload(fs.getFilename(show.Show), show.Payload)
	if err != nil {
		return err
	}
	meta, err := show.Show.MarshalJSON()
	if err != nil {
		return err
	}
	err = fs.upload(fs.getMetadataFilename(show.Show), meta)
	return nil
}

func (fs *gcsStorage) upload(key string, body []byte) error {
	w := fs.client.Bucket(fs.bucket).Object(key).NewWriter(fs.ctx)
	defer w.Close()
	if _, err := w.Write(body); err != nil {
		return err
	}
	return nil
}

func (fs gcsStorage) getFilename(s radl.Show) string {
	return fmt.Sprintf("%s-%s-%d%s", s.SourceId(), s.SeriesId(), s.Index(), path.Ext(s.MediaUrl().Path))
}

func (fs gcsStorage) getMetadataFilename(s radl.Show) string {
	return fmt.Sprintf("%s-%s-%d.json", s.SourceId(), s.SeriesId(), s.Index())
}
