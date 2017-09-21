package sqsd

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type StatServer struct {
	Srv *http.Server
}

func NewStatServer(tr *JobTracker, p int) *StatServer {
	handler := StatHandler{Tracker: tr}
	return &StatServer{
		Srv: &http.Server{
			Addr:    ":" + strconv.Itoa(p),
			Handler: handler.BuildServeMux(),
		},
	}
}

func (s *StatServer) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("stat server start.")

	syncWait := &sync.WaitGroup{}
	syncWait.Add(1)
	go func() {
		defer syncWait.Done()
		if err := s.Srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	syncWait.Add(1)
	go func() {
		defer syncWait.Done()
		for {
			select {
			case <-ctx.Done():
				if err := s.Srv.Shutdown(ctx); err != nil {
					log.Fatal(err)
				}
				return
			}
		}
	}()
	syncWait.Wait()

	log.Println("stat server closed.")
}
