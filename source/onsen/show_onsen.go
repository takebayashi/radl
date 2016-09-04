package onsen

import (
	"encoding/json"
	"github.com/takebayashi/gonsen"
	"net/url"
)

type OnsenShow struct {
	raw gonsen.Program
}

func (s OnsenShow) SourceId() string {
	return "onsen"
}

func (s OnsenShow) SeriesId() string {
	return s.raw.Slug
}

func (s OnsenShow) Title() string {
	return s.raw.Title
}

func (s OnsenShow) Index() int {
	return s.raw.Index
}

func (s OnsenShow) MediaUrl() *url.URL {
	u, err := url.Parse(s.raw.MediaUrl)
	if err != nil || s.raw.MediaUrl == "" {
		return nil
	}
	return u
}

func (s OnsenShow) MarshalJSON() ([]byte, error) {
	st := struct {
		SourceId string
		SeriesId string
		Title    string
		Index    int
		Url      string
	}{
		SourceId: s.SourceId(),
		SeriesId: s.SeriesId(),
		Title:    s.Title(),
		Index:    s.Index(),
		Url:      s.MediaUrl().String(),
	}
	return json.MarshalIndent(st, "", "  ")
}
