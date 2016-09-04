package file

import (
	"fmt"
	"github.com/takebayashi/radl"
	"os"
	"path"
)

type fileStore struct {
	baseDir string
}

func NewFileStore(baseDir string) (*fileStore, error) {
	return &fileStore{baseDir: baseDir}, nil
}

func (fs *fileStore) Save(show radl.ShowPayload) error {
	if err := fs.writeToFile(fs.getFilename(show.Show), show.Payload); err != nil {
		return err
	}
	meta, err := show.Show.MarshalJSON()
	if err != nil {
		return err
	}
	if err = fs.writeToFile(fs.getMetadataFilename(show.Show), meta); err != nil {
		return err
	}
	return nil
}

func (fs fileStore) writeToFile(filename string, bytes []byte) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(bytes)
	return err
}

func (fs fileStore) getFilename(s radl.Show) string {
	name := fmt.Sprintf("%s-%s-%d%s", s.SourceId(), s.SeriesId(), s.Index(), path.Ext(s.MediaUrl().Path))
	return path.Join(fs.baseDir, name)
}

func (fs fileStore) getMetadataFilename(s radl.Show) string {
	name := fmt.Sprintf("%s-%s-%d.json", s.SourceId(), s.SeriesId(), s.Index())
	return path.Join(fs.baseDir, name)
}
