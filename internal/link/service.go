package link

import (
	"context"
	"log"
	"net/url"

	"github.com/pkab00/shortenit/pkg/encode"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func fixURL(url_str string) (*string, error) {
	u, err := url.Parse(url_str)
	if err != nil {
		log.Println("error parsing URL: ", err)
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}
	res := u.String()

	return &res, nil
}

func (s *Service) Create(ctx context.Context, url string) (*LinkResponse, error) {
	var err error

	fixedUrl, err := fixURL(url)
	if err != nil {
		return nil, err
	}

	link, err := s.repo.Create(ctx, *fixedUrl)
	if err != nil {
		return nil, err
	}
	return link.toResponse(), nil
}

func (s *Service) Delete(ctx context.Context, code string) (*LinkResponse, error) {
	var err error

	id, err := encode.NewDecoder().Decode(code)
	if err != nil {
		return nil, err
	}

	link, err := s.repo.Delete(ctx, *id)
	if err != nil {
		return nil, err
	}
	return link.toResponse(), nil
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

func (s *Service) ByCode(ctx context.Context, code string) (*LinkResponse, error) {
	var err error

	id, err := encode.NewDecoder().Decode(code)
	if err != nil {
		return nil, err
	}

	link, err := s.repo.ByID(ctx, *id)
	if err != nil {
		return nil, err
	}

	return link.toResponse(), nil
}
