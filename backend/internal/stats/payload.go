package stats

type StatsPayload struct {
	PeriodFrom string `json:"period_from"`
	PeriodTo   string `json:"period_to"`
	Clicks     int    `json:"clicks"`
}

type StatsResponse struct {
	Stats       []StatsPayload `json:"stats"`
	TotalClicks int            `json:"total_clicks"`
}
