package stats

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Stats struct {
	gorm.Model
	LinkId     uint           `json:"link_id"`
	ClickCount uint           `json:"click_count"`
	Date       datatypes.Date `json:"date"`
}
