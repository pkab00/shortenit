package link

type LinkResponse struct {
	Code      string `json:"link_code"`
	URL       string `json:"link_body"`
	CreatedAt string `json:"created_at"`
}

type Link struct {
	ID        int
	Code      string
	URL       string
	CreatedAt string
}
