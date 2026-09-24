package link

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"

	"github.com/pkab00/shortenit/internal/apperr"
	stringutils "github.com/pkab00/shortenit/pkg/string_utils"
)

type CreateResult struct {
	Response *LinkResponse
	Created  bool
	Error    error
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

func (s *Service) checkCustomCode(ctx context.Context, code string) error {
	const MAX_LEN int = 10
	RESERVED_WORDS := []string{
		"shortenit",
		"swagger",
	}

	if slices.Contains(RESERVED_WORDS, code) {
		return fmt.Errorf("custom code error: %w", apperr.ErrorCustomCodeInUse)
	}
	if len(code) > MAX_LEN || stringutils.IsEmpty(code) {
		return fmt.Errorf("custom code error: %w", apperr.ErrorInvalidCustomCode)
	}

	_, err := s.repo.ByCode(ctx, code)
	if err == nil || errors.Is(err, apperr.ErrorLinkNotFound) {
		return nil
	} else {
		return fmt.Errorf("custom code error: %w", err)
	}
}

func (s *Service) Create(ctx context.Context, req *CreateLinkRequest) (*CreateResult, error) {
	var err error

	if strings.TrimSpace(req.URL) == "" {
		return nil, apperr.ErrorInvalidCreateRequest
	}

	fixedUrl, err := fixURL(req.URL)
	if err != nil {
		return nil, err
	}

	foundByURL, _ := s.repo.ByURL(ctx, *fixedUrl)
	if foundByURL != nil {
		return &CreateResult{
			Response: foundByURL.toResponse(),
			Created:  false,
		}, nil
	}

	var customCode *string = nil
	if req.Code != nil {
		err = s.checkCustomCode(ctx, *req.Code)
		if err != nil {
			return nil, err
		}
		customCode = req.Code

		link, _ := s.repo.ByCode(ctx, *customCode)
		if link != nil {
			return nil, apperr.ErrorCustomCodeInUse
		}
	}

	link, err := s.repo.Create(ctx, *fixedUrl, customCode)
	if err != nil {
		return nil, err
	}
	return &CreateResult{
		Response: link.toResponse(),
		Created:  true,
	}, nil
}

func (s *Service) CreateMany(ctx context.Context, reqs []CreateLinkRequest) []CreateResult {
	var results []CreateResult

	for _, req := range reqs {
		res, err := s.Create(ctx, &req)
		if err == nil {
			results = append(results, *res)
		} else {
			failureRes := CreateResult{
				Response: nil,
				Created:  false,
				Error:    err,
			}
			results = append(results, failureRes)
		}
	}
	return results
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
