package model

type Stat struct {
	ClicksTotal int     `json:"clicks_total"`
	Clicks      []Click `json:"clicks"`
}

type Click struct {
	ID        string `json:"id" db:"id"`
	Url       string `json:"url" db:"url"`
	Time      string `json:"time" db:"time"`
	UserAgent string `json:"user_agent" db:"user_agent"`
}
