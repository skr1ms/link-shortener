package stats

import (
	"linkshortener/config"
	"linkshortener/pkg/res"
	"net/http"
	"time"
)

type StatsHandlerDeps struct {
	Config          *config.Config
	StatsRepository *StatsRepository
}

type StatsHandler struct {
	deps *StatsHandlerDeps
}

func NewStatsHandler(router *http.ServeMux, deps *StatsHandlerDeps) {
	statsHandler := &StatsHandler{
		deps: deps,
	}
	router.HandleFunc("GET /stats", statsHandler.GetStats())
}

func (handler *StatsHandler) GetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startDate, err := time.Parse("2006-01-02", r.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, "from and to are required", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "invalid end date", http.StatusBadRequest)
			return
		}

		if startDate.After(endDate) {
			http.Error(w, "start date must be before end date", http.StatusBadRequest)
			return
		}

		by := r.URL.Query().Get("by")
		if by != "day" && by != "month" {
			http.Error(w, "invalid by", http.StatusBadRequest)
			return
		}

		stats := handler.deps.StatsRepository.GetStats(by, startDate, endDate)

		res.Response(w, http.StatusOK, stats)
	}
}
