package main

import (
	"flag"
	"github.com/takebayashi/gonsen"
	"github.com/takebayashi/radl"
	"github.com/takebayashi/radl/source/onsen"
	fstate "github.com/takebayashi/radl/state/file"
	"github.com/takebayashi/radl/store/gcs"
	"log"
)

func preparePlans() []radl.Plan {
	plans := []radl.Plan{}
	names, err := gonsen.GetProgramNames()
	if err != nil {
		return plans
	}
	for _, name := range names {
		p := radl.Plan{
			SourceId: "onsen",
			ShowId:   name,
		}
		plans = append(plans, p)
	}
	return plans
}

func main() {
	var bucket = flag.String("bucket", "", "GCS bucket name")
	var stateFile = flag.String("state", "radl.state", "RADL state file")
	flag.Parse()
	if *bucket == "" {
		log.Fatal("-bucket argument is required")
	}
	plans := preparePlans()
	sources := map[string]radl.Source{
		"onsen": onsen.OnsenSource{},
	}
	state, err := fstate.NewFileState(*stateFile)
	if err != nil {
		panic(err)
	}
	store, err := gcs.NewGCSStore(*bucket)
	if err != nil {
		panic(err)
	}
	if err := radl.ExecutePlans(plans, sources, state, store); err != nil {
		log.Fatal(err)
	}
}
