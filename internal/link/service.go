package link

import (
	"context"
	"log"
	"net/url"
	"strings"

	"github.com/pkab00/shortenit/internal/apperr"
)

type CreateResult struct {
	Link    *LinkResponse
	Created bool
}

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

func (s *Service) Create(ctx context.Context, url string) (*CreateResult, error) {
	var err error

	if strings.TrimSpace(url) == "" {
		return nil, apperr.ErrorInvalidCreateRequest
	}

	fixedUrl, err := fixURL(url)
	if err != nil {
		return nil, err
	}

	foundByURL, _ := s.repo.ByURL(ctx, *fixedUrl)
	if foundByURL != nil {
		return &CreateResult{
			Link:    foundByURL.toResponse(),
			Created: false,
		}, nil
	}

	link, err := s.repo.Create(ctx, *fixedUrl)
	if err != nil {
		return nil, err
	}
	return &CreateResult{
		Link:    link.toResponse(),
		Created: true,
	}, nil
}

func (s *Service) Delete(ctx context.Context, code string) (*LinkResponse, error) {
	var err error

	link, err := s.repo.Delete(ctx, code)
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

	link, err := s.repo.ByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return link.toResponse(), nil
}
