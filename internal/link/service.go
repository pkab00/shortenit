package link

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, url string) (*LinkResponse, error) {
	link, err := s.repo.Create(ctx, url)
	return link.toResponse(), err
}

func (s *Service) Delete(ctx context.Context, url string) (*LinkResponse, error) {
	link, err := s.repo.Delete(ctx, url)
	return link.toResponse(), err
}

func (s *Service) All(ctx context.Context) ([]LinkResponse, error) {
	var res []LinkResponse

	links, err := s.repo.All(ctx)
	if err != nil {
		return nil, err
	}

	for _, link := range links {
		res = append(res, *link.toResponse())
	}
	return res, nil
}
