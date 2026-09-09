package redirect

import (
	"context"

	"github.com/pkab00/shortenit/pkg/encode"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Increment(ctx context.Context, code string) (*Redirect, error) {
	var err error

	id, err := encode.NewDecoder().Decode(code)
	if err != nil {
		return nil, err
	}

	res, err := s.repo.Increment(ctx, *id)
	if err != nil {
		return nil, err
	}
	return res, nil
}
