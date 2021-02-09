package worker

import (
	"log"
	"time"

	"github.com/dmitry-vovk/csv-chg-go/api"
)

// Run is the main worker loop
func (w *Worker) Run() {
	t := time.NewTicker(w.interval)
out:
	for {
		select {
		case <-w.doneC:
			break out
		case id := <-w.deleteC:
			delete(w.uuids, id)
		case <-t.C:
			for id := range w.uuids {
				w.limitC <- struct{}{}
				w.wg.Add(1)
				go w.process(id)
			}
			w.wg.Wait()
		}
	}
	t.Stop()
	close(w.stoppedC)
}

// process takes a UUID and runs API queries against it
func (w *Worker) process(id compact) {
	uuid := id.String()
	if item, err := w.client.GetItem(uuid); err != nil {
		if err == api.ErrBadRequest {
			log.Printf("API indicated UUID %q not found, removing", uuid)
			go func() { w.deleteC <- id }()
		} else {
			log.Printf("API error: %s", err)
		}
	} else if item.UUID != uuid {
		log.Printf("APi returned wrong item, expected %q, got %q", uuid, item.UUID)
	} else if item.Quantity < 5 {
		if err = w.client.PostAlert(uuid); err != nil {
			if err == api.ErrBadRequest {
				log.Printf("API indicated UUID %q not found, removing", uuid)
				go func() { w.deleteC <- id }()
			} else {
				log.Printf("API error: %s", err)
			}
		}
	}
	<-w.limitC
	w.wg.Done()
}
