package statistics

func (s *Statistics) toResponse() *StatisticsResponse {
	var rCount int

	if s.Redirects == nil {
		rCount = 0
	} else {
		rCount = *s.Redirects
	}

	return &StatisticsResponse{
		Code:      s.Code,
		URL:       s.URL,
		CreatedAt: s.CreatedAt,
		Redirects: rCount,
	}
}
