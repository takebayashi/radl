package onsen

import (
	"github.com/takebayashi/gonsen"
	"github.com/takebayashi/radl"
)

type OnsenSource struct {
}

func (s OnsenSource) Title() string {
	return "Onsen"
}

func (s OnsenSource) Id() string {
	return "onsen"
}

func (s OnsenSource) GetShow(id string) (radl.Show, error) {
	p, err := gonsen.GetProgram(id)
	if err != nil {
		return nil, err
	}
	return OnsenShow{raw: p}, nil
}
