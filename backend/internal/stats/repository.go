package stats

import (
	"linkshortener/pkg/db"
	"time"

	"gorm.io/datatypes"
)

const (
	StatsByDay   = "day"
	StatsByMonth = "month"
)

type StatsRepository struct {
	db *db.Db
}

func NewStatsRepository(db *db.Db) *StatsRepository {
	return &StatsRepository{db: db}
}

func (repo *StatsRepository) AddClick(linkId uint) error {
	var stats Stats
	currentDate := datatypes.Date(time.Now())
	repo.db.DB.Find(&stats, "link_id = ? AND date = ?", linkId, currentDate)
	if stats.ID == 0 {
		repo.db.DB.Create(&Stats{
			LinkId:     linkId,
			ClickCount: 1,
			Date:       currentDate,
		})
	} else {
		stats.ClickCount++
		repo.db.DB.Save(&stats)
	}

	return nil
}

func (repo *StatsRepository) GetStats(by string, startDate, endDate time.Time) StatsResponse {
	var stats []StatsPayload
	var totalClicks int

	repo.db.DB.Table("stats").
		Select("sum(click_count)").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Scan(&totalClicks)

	switch by {
	case StatsByDay:
		repo.db.DB.Table("stats").
			Select("to_char(date, 'YYYY-MM-DD') as period_from, to_char(date, 'YYYY-MM-DD') as period_to, sum(click_count) as clicks").
			Where("date BETWEEN ? AND ?", startDate, endDate).
			Group("to_char(date, 'YYYY-MM-DD')").
			Order("to_char(date, 'YYYY-MM-DD')").
			Scan(&stats)
	case StatsByMonth:
		repo.db.DB.Table("stats").
			Select("to_char(date_trunc('month', date), 'YYYY-MM-DD') as period_from, to_char(date_trunc('month', date) + interval '1 month' - interval '1 day', 'YYYY-MM-DD') as period_to, sum(click_count) as clicks").
			Where("date BETWEEN ? AND ?", startDate, endDate).
			Group("date_trunc('month', date)").
			Order("date_trunc('month', date)").
			Scan(&stats)
	}

	return StatsResponse{
		Stats: []StatsPayload{
			{
				PeriodFrom: startDate.Format("2006-01-02"),
				PeriodTo:   endDate.Format("2006-01-02"),
				Clicks:     totalClicks,
			},
		},
		TotalClicks: totalClicks,
	}
}
