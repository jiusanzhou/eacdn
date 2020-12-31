package service

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/gorilla/mux"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

// HandleHealth ...
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

// HandleStat ...
func HandleStat(w http.ResponseWriter, r *http.Request) {
	// TODO: return stat of the process, maybe can be replaced by metrics
	// TODO: add debug handler
}

func (s *Service) installHandlers(r *mux.Router) {

	// install headth handler
	r.HandleFunc("/_healthz", HandleHealth)

	// install pprof for debug from http default serve mux
	if s.Config.Debug {
		log.Println("I Start the pprof handler")

		debugr := r.PathPrefix("/debug").Subrouter()
		// TODO: add other endpoint for debug

		debugr.NotFoundHandler = http.DefaultServeMux
	}

	// install the stat handler for summary
	r.HandleFunc("/stat", HandleStat)

	// install the metrics handler with prometheus
	// TODO: install metrics to promhttp
	// r.Handle("/metrics", promhttp.Handler())

	// install ui
	// r.NotFoundHandler = ui.NewHandler(ui.Prefix(s.Config.RootPath))

	apiv1 := r.PathPrefix("/api/v1/").Subrouter()
	// TODO: install handler under the apiv1

	_ = apiv1
}

func (s *Service) startHTTP() error {
	r := mux.NewRouter()

	// call the handler installer
	s.installHandlers(r)

	log.Printf("I Listen EaCDN service on %s\n", s.Config.Addr)

	// listen and serve the http service
	// blocking at here
	return http.ListenAndServe(s.Config.Addr, r)
}
