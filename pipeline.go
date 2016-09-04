package radl

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Plan struct {
	SourceId string
	ShowId   string
}

func ExecutePlans(plans []Plan, sources map[string]Source, state State, store Store) error {
	for _, plan := range plans {
		show, err := sources[plan.SourceId].GetShow(plan.ShowId)
		if err != nil {
			log.Printf("[ERROR] failed to get show: %+v; %+v", plan, err)
			return err
		}
		label := fmt.Sprintf("%s-%s-%d", show.SourceId(), show.SeriesId(), show.Index())
		if state.IsNew(show) && show.MediaUrl() != nil {
			log.Printf("download: %s; from %s", label, show.MediaUrl().String())
			res, err := http.Get(show.MediaUrl().String())
			if err != nil {
				log.Printf("[ERROR] failed to download: %s; %+v", label, err)
				return err
			}
			defer res.Body.Close()
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Printf("[ERROR] failed to download: %s; %+v", label, err)
				return err
			}
			sp := ShowPayload{Show: show, Payload: b}
			if err = store.Save(sp); err != nil {
				log.Printf("[ERROR] failed to store: %s; %+v", label, err)
				return err
			}
			log.Printf("downloaded: %s", label)
			if err = state.Update(show); err != nil {
				log.Printf("[ERROR] failed to update status: %s; %+v", label, err)
				return err
			}
		} else {
			log.Printf("skipped: %s", label)
		}
	}
	return nil
}
