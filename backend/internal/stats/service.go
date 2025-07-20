package stats

import (
	"linkshortener/pkg/di"
	"linkshortener/pkg/event"
	"log"
)

type StatsServiceDeps struct {
	EventBus        di.IEventBus
	StatsRepository di.IStatsRepository
}

type StatsService struct {
	deps *StatsServiceDeps
}

func NewStatsService(deps *StatsServiceDeps) *StatsService {
	return &StatsService{deps: deps}
}

func (s *StatsService) AddClick() {
	for msg := range s.deps.EventBus.Subscribe() {
		if msg.Type == event.LinkClicked {
			id, ok := msg.Data.(uint)
			if !ok {
				log.Fatalln("Bad LinkClicked Data: ", msg.Data)
				continue
			}
			s.deps.StatsRepository.AddClick(id)
		}
	}
}
