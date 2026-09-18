package redirect

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Increment(ctx context.Context, code string) (*Redirect, error) {
	var err error

	res, err := s.repo.Increment(ctx, code)
	if err != nil {
		return nil, err
	}
	return res, nil
}
