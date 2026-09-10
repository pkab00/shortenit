package statistics

import (
	"context"

	"github.com/pkab00/shortenit/pkg/encode"
)

type Service struct {
	repo *PostgresRepository
}

func NewService(repo *PostgresRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, code string) (*StatisticsResponse, error) {
	var err error

	id, err := encode.NewDecoder().Decode(code)
	if err != nil {
		return nil, err
	}

	stat, err := s.repo.Get(ctx, *id)
	if err != nil {
		return nil, err
	}

	return stat.toResponse(), nil
}
