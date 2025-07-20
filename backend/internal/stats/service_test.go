package stats_test

import (
	"linkshortener/internal/stats"
	"linkshortener/pkg/event"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/datatypes"
)

type MockEventBus struct {
	channel chan event.Event
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		channel: make(chan event.Event, 10),
	}
}

func (m *MockEventBus) Publish(evt event.Event) {
	m.channel <- evt
}

func (m *MockEventBus) Subscribe() <-chan event.Event {
	return m.channel
}

type MockStatsRepository struct {
	addClickCalls []uint
	stats         map[uint][]stats.Stats
}

func NewMockStatsRepository() *MockStatsRepository {
	return &MockStatsRepository{
		addClickCalls: make([]uint, 0),
		stats:         make(map[uint][]stats.Stats),
	}
}

func (m *MockStatsRepository) AddClick(linkId uint) error {
	m.addClickCalls = append(m.addClickCalls, linkId)

	today := time.Now()
	currentDate := datatypes.Date(today)

	found := false
	for i, stat := range m.stats[linkId] {
		statDay := time.Time(stat.Date).Format("2006-01-02")
		currentDay := today.Format("2006-01-02")
		if statDay == currentDay {
			m.stats[linkId][i].ClickCount++
			found = true
			break
		}
	}

	if !found {
		if m.stats[linkId] == nil {
			m.stats[linkId] = make([]stats.Stats, 0)
		}
		m.stats[linkId] = append(m.stats[linkId], stats.Stats{
			LinkId:     linkId,
			ClickCount: 1,
			Date:       currentDate,
		})
	}

	return nil
}

func setupStatsService() (*stats.StatsService, *MockEventBus, *MockStatsRepository) {
	mockEventBus := NewMockEventBus()
	mockStatsRepo := NewMockStatsRepository()

	deps := &stats.StatsServiceDeps{
		EventBus:        mockEventBus,
		StatsRepository: mockStatsRepo,
	}

	statsService := stats.NewStatsService(deps)
	return statsService, mockEventBus, mockStatsRepo
}

func TestStatsServiceAddClick(t *testing.T) {
	godotenv.Load()

	statsService, mockEventBus, mockStatsRepo := setupStatsService()

	go func() {
		statsService.AddClick()
	}()

	linkId := uint(123)
	mockEventBus.Publish(event.Event{
		Type: event.LinkClicked,
		Data: linkId,
	})

	time.Sleep(100 * time.Millisecond)

	if len(mockStatsRepo.addClickCalls) != 1 {
		t.Fatalf("Expected 1 AddClick call, got %d", len(mockStatsRepo.addClickCalls))
	}

	if mockStatsRepo.addClickCalls[0] != linkId {
		t.Fatalf("Expected linkId %d, got %d", linkId, mockStatsRepo.addClickCalls[0])
	}
}

func TestStatsServiceMultipleClicks(t *testing.T) {
	godotenv.Load()

	statsService, mockEventBus, mockStatsRepo := setupStatsService()

	go func() {
		statsService.AddClick()
	}()

	linkIds := []uint{123, 456, 123}

	for _, linkId := range linkIds {
		mockEventBus.Publish(event.Event{
			Type: event.LinkClicked,
			Data: linkId,
		})
	}

	time.Sleep(200 * time.Millisecond)

	if len(mockStatsRepo.addClickCalls) != 3 {
		t.Fatalf("Expected 3 AddClick calls, got %d", len(mockStatsRepo.addClickCalls))
	}

	if len(mockStatsRepo.stats[123]) != 1 {
		t.Fatalf("Expected 1 stats entry for link 123, got %d", len(mockStatsRepo.stats[123]))
	}

	if mockStatsRepo.stats[123][0].ClickCount != 2 {
		t.Fatalf("Expected 2 clicks for link 123, got %d", mockStatsRepo.stats[123][0].ClickCount)
	}

	if len(mockStatsRepo.stats[456]) != 1 {
		t.Fatalf("Expected 1 stats entry for link 456, got %d", len(mockStatsRepo.stats[456]))
	}

	if mockStatsRepo.stats[456][0].ClickCount != 1 {
		t.Fatalf("Expected 1 click for link 456, got %d", mockStatsRepo.stats[456][0].ClickCount)
	}
}

func TestStatsServiceIgnoreNonLinkClickedEvents(t *testing.T) {
	godotenv.Load()

	statsService, mockEventBus, mockStatsRepo := setupStatsService()

	go func() {
		statsService.AddClick()
	}()

	mockEventBus.Publish(event.Event{
		Type: "other.event",
		Data: uint(123),
	})

	mockEventBus.Publish(event.Event{
		Type: event.LinkClicked,
		Data: uint(456),
	})

	time.Sleep(100 * time.Millisecond)

	if len(mockStatsRepo.addClickCalls) != 1 {
		t.Fatalf("Expected 1 AddClick call, got %d", len(mockStatsRepo.addClickCalls))
	}

	if mockStatsRepo.addClickCalls[0] != 456 {
		t.Fatalf("Expected linkId 456, got %d", mockStatsRepo.addClickCalls[0])
	}
}
