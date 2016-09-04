package file

import (
	"fmt"
	"github.com/takebayashi/radl"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type fileState struct {
	filename string
	statuses fileStateEntries
}

type fileStateEntries []*fileStateEntry

func (es fileStateEntries) MarshalText() (text []byte, err error) {
	ms := []string{}
	for _, e := range es {
		m, merr := e.MarshalText()
		if merr != nil {
			err = merr
			return
		}
		ms = append(ms, string(m))
	}
	text = []byte(strings.Join(ms, "\n"))
	return
}

func (es *fileStateEntries) UnmarshalText(text []byte) error {
	for _, s := range strings.Split(string(text), "\n") {
		e := &fileStateEntry{}
		if err := e.UnmarshalText([]byte(s)); err != nil {
			return err
		}
		*es = append(*es, e)
	}
	return nil
}

func (es *fileStateEntries) Add(e fileStateEntry) {
	*es = append(*es, &e)
}

type fileStateEntry struct {
	SourceId string
	ShowId   string
	Index    int
	Done     bool
}

func (e *fileStateEntry) MarshalText() (text []byte, err error) {
	text = []byte(fmt.Sprintf("%s,%s,%d,%t", e.SourceId, e.ShowId, e.Index, e.Done))
	err = nil
	return
}

func (e *fileStateEntry) UnmarshalText(text []byte) error {
	chunks := strings.Split(string(text), ",")
	if len(chunks) != 4 {
		return fmt.Errorf("expected 4 fields, but %d fields passed", len(chunks))
	}
	e.SourceId = chunks[0]
	e.ShowId = chunks[1]
	var err error
	e.Index, err = strconv.Atoi(chunks[2])
	if err != nil {
		return err
	}
	e.Done = chunks[3] == "true"
	return nil
}

func NewFileState(filename string) (*fileState, error) {
	fs := &fileState{}
	fs.filename = filename
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	fs.statuses.UnmarshalText(b)
	return fs, nil
}

func (fs *fileState) IsNew(show radl.Show) bool {
	for _, e := range fs.statuses {
		if e.SourceId == show.SourceId() &&
			e.ShowId == show.SeriesId() &&
			e.Index == show.Index() &&
			e.Done {
			return false
		}
	}
	return true
}

func (fs *fileState) Update(show radl.Show) error {
	for _, e := range fs.statuses {
		if e.SourceId == show.SourceId() &&
			e.ShowId == show.SeriesId() &&
			e.Index == show.Index() {
			e.Done = true
			return nil
		}
	}
	fs.statuses.Add(fileStateEntry{
		SourceId: show.SourceId(),
		ShowId:   show.SeriesId(),
		Index:    show.Index(),
		Done:     true,
	})
	os.Truncate(fs.filename, 0)
	f, err := os.OpenFile(fs.filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := fs.statuses.MarshalText()
	if err != nil {
		return err
	}
	f.Write(b)
	return nil
}
