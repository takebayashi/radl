package radl

import (
	"encoding/json"
	"net/url"
)

type Show interface {
	json.Marshaler
	SourceId() string
	SeriesId() string
	Title() string
	Index() int
	MediaUrl() *url.URL
}

type ShowPayload struct {
	Show    Show
	Payload []byte
}
