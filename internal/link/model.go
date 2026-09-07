package link

type LinkResponse struct {
	Code      string `json:"link_code"`
	Body      string `json:"link_body"`
	CreatedAt string `json:"created_at"`
}

type Link struct {
	ID        int    `json:"link_id"`
	Body      string `json:"link_body"`
	CreatedAt string `json:"created_at"`
}
