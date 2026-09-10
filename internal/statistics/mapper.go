package statistics

import "github.com/pkab00/shortenit/pkg/encode"

func (s *Statistics) toResponse() *StatisticsResponse {
	var rCount int

	if s.Redirects == nil {
		rCount = 0
	} else {
		rCount = *s.Redirects
	}

	code := encode.NewEncoder().Encode(uint64(s.ID))

	return &StatisticsResponse{
		Code:      code,
		URL:       s.URL,
		CreatedAt: s.CreatedAt,
		Redirects: rCount,
	}
}
