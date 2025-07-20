package stats_test

import (
	"linkshortener/internal/stats"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/datatypes"
)

type MockStatsRepositoryImpl struct {
	stats  []stats.Stats
	nextID uint
}

func NewMockStatsRepositoryImpl() *MockStatsRepositoryImpl {
	return &MockStatsRepositoryImpl{
		stats:  make([]stats.Stats, 0),
		nextID: 1,
	}
}

func (repo *MockStatsRepositoryImpl) AddClick(linkId uint) error {
	today := time.Now()
	currentDate := datatypes.Date(today)

	for i, stat := range repo.stats {
		statDay := time.Time(stat.Date).Format("2006-01-02")
		currentDay := today.Format("2006-01-02")
		if stat.LinkId == linkId && statDay == currentDay {
			repo.stats[i].ClickCount++
			return nil
		}
	}

	newStat := stats.Stats{
		LinkId:     linkId,
		ClickCount: 1,
		Date:       currentDate,
	}
	newStat.ID = repo.nextID
	repo.nextID++

	repo.stats = append(repo.stats, newStat)
	return nil
}

func (repo *MockStatsRepositoryImpl) GetStats(by string, startDate, endDate time.Time) stats.StatsResponse {
	totalClicks := 0

	// Подсчитываем общее количество кликов в диапазоне
	for _, stat := range repo.stats {
		statDate := time.Time(stat.Date)
		if !statDate.Before(startDate) && !statDate.After(endDate) {
			totalClicks += int(stat.ClickCount)
		}
	}

	return stats.StatsResponse{
		Stats: []stats.StatsPayload{
			{
				PeriodFrom: startDate.Format("2006-01-02"),
				PeriodTo:   endDate.Format("2006-01-02"),
				Clicks:     totalClicks,
			},
		},
		TotalClicks: totalClicks,
	}
}

func TestStatsRepositoryAddClickNewRecord(t *testing.T) {
	godotenv.Load()

	repo := NewMockStatsRepositoryImpl()
	linkId := uint(123)

	err := repo.AddClick(linkId)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(repo.stats) != 1 {
		t.Fatalf("Expected 1 stats record, got %d", len(repo.stats))
	}

	stat := repo.stats[0]
	if stat.LinkId != linkId {
		t.Fatalf("Expected linkId %d, got %d", linkId, stat.LinkId)
	}

	if stat.ClickCount != 1 {
		t.Fatalf("Expected ClickCount 1, got %d", stat.ClickCount)
	}

	today := time.Now().Format("2006-01-02")
	statDateStr := time.Time(stat.Date).Format("2006-01-02")
	if statDateStr != today {
		t.Fatalf("Expected today's date %s, got %s", today, statDateStr)
	}
}

func TestStatsRepositoryAddClickExistingRecord(t *testing.T) {
	godotenv.Load()

	repo := NewMockStatsRepositoryImpl()
	linkId := uint(456)

	err := repo.AddClick(linkId)
	if err != nil {
		t.Fatalf("Expected no error on first click, got %v", err)
	}

	err = repo.AddClick(linkId)
	if err != nil {
		t.Fatalf("Expected no error on second click, got %v", err)
	}

	if len(repo.stats) != 1 {
		t.Fatalf("Expected 1 stats record, got %d", len(repo.stats))
	}

	stat := repo.stats[0]
	if stat.ClickCount != 2 {
		t.Fatalf("Expected ClickCount 2, got %d", stat.ClickCount)
	}
}

func TestStatsRepositoryAddClickMultipleLinks(t *testing.T) {
	godotenv.Load()

	repo := NewMockStatsRepositoryImpl()

	linkIds := []uint{111, 222, 333, 111, 222}

	for _, linkId := range linkIds {
		err := repo.AddClick(linkId)
		if err != nil {
			t.Fatalf("Expected no error for linkId %d, got %v", linkId, err)
		}
	}

	if len(repo.stats) != 3 {
		t.Fatalf("Expected 3 stats records, got %d", len(repo.stats))
	}

	clickCounts := make(map[uint]uint)
	for _, stat := range repo.stats {
		clickCounts[stat.LinkId] = stat.ClickCount
	}

	if clickCounts[111] != 2 {
		t.Fatalf("Expected 2 clicks for link 111, got %d", clickCounts[111])
	}

	if clickCounts[222] != 2 {
		t.Fatalf("Expected 2 clicks for link 222, got %d", clickCounts[222])
	}

	if clickCounts[333] != 1 {
		t.Fatalf("Expected 1 click for link 333, got %d", clickCounts[333])
	}
}

func TestStatsRepositoryGetStatsByDay(t *testing.T) {
	godotenv.Load()

	repo := NewMockStatsRepositoryImpl()

	today := time.Now()
	todayDate := datatypes.Date(today)

	stat1 := stats.Stats{LinkId: 123, ClickCount: 5, Date: todayDate}
	stat1.ID = 1
	stat2 := stats.Stats{LinkId: 456, ClickCount: 3, Date: todayDate}
	stat2.ID = 2

	repo.stats = []stats.Stats{stat1, stat2}

	startDate := today.Truncate(24 * time.Hour)
	endDate := startDate.Add(24 * time.Hour)

	response := repo.GetStats(stats.StatsByDay, startDate, endDate)

	if response.TotalClicks != 8 {
		t.Fatalf("Expected total clicks 8, got %d", response.TotalClicks)
	}

	if len(response.Stats) != 1 {
		t.Fatalf("Expected 1 stats entry, got %d", len(response.Stats))
	}

	if response.Stats[0].Clicks != 8 {
		t.Fatalf("Expected 8 clicks in stats entry, got %d", response.Stats[0].Clicks)
	}
}

func TestStatsRepositoryGetStatsEmptyRange(t *testing.T) {
	godotenv.Load()

	repo := NewMockStatsRepositoryImpl()

	today := time.Now()
	stat := stats.Stats{LinkId: 123, ClickCount: 5, Date: datatypes.Date(today)}
	stat.ID = 1
	repo.stats = []stats.Stats{stat}

	lastWeek := today.AddDate(0, 0, -7)
	startDate := lastWeek.AddDate(0, 0, -1)
	endDate := lastWeek

	response := repo.GetStats(stats.StatsByDay, startDate, endDate)

	if response.TotalClicks != 0 {
		t.Fatalf("Expected total clicks 0, got %d", response.TotalClicks)
	}
}

func TestStatsRepositoryGetStatsByMonth(t *testing.T) {
	godotenv.Load()

	repo := NewMockStatsRepositoryImpl()

	today := time.Now()
	todayDate := datatypes.Date(today)

	stat1 := stats.Stats{LinkId: 123, ClickCount: 10, Date: todayDate}
	stat1.ID = 1
	stat2 := stats.Stats{LinkId: 456, ClickCount: 15, Date: todayDate}
	stat2.ID = 2

	repo.stats = []stats.Stats{stat1, stat2}

	startOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	response := repo.GetStats(stats.StatsByMonth, startOfMonth, endOfMonth)

	if response.TotalClicks != 25 {
		t.Fatalf("Expected total clicks 25, got %d", response.TotalClicks)
	}
}
