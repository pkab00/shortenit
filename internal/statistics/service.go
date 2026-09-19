package statistics

import (
	"context"
)

type Service struct {
	repo *PostgresRepository
}

func NewService(repo *PostgresRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, code string) (*StatisticsResponse, error) {
	var err error

	stat, err := s.repo.Get(ctx, code)
	if err != nil {
		return nil, err
	}

	return stat.toResponse(), nil
}
