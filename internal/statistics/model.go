package statistics

type Statistics struct {
	ID        int
	Code      string
	URL       string
	CreatedAt string
	Redirects *int
}

type StatisticsResponse struct {
	Code      string `json:"code"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	Redirects int    `json:"redirects_count"`
}
