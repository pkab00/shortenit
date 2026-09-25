package link_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/pkab00/shortenit/internal/apperr"
	"github.com/pkab00/shortenit/internal/link"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -destination=repository_test_mock.go -source=repository.go -package=link

const (
	googleURL      = "www.google.com"
	googleFixedURL = "https://www.google.com"
	googleCode     = "google"

	yandexURL      = "www.yandex.com"
	yandexFixedURL = "https://www.yandex.com"
	yandexCode     = "yandex"
)

func request(url string) *link.CreateLinkRequest {
	return &link.CreateLinkRequest{URL: url}
}

func requestWithCode(url, code string) *link.CreateLinkRequest {
	return &link.CreateLinkRequest{
		URL:  url,
		Code: &code,
	}
}

func newService(t *testing.T) (*link.Service, *link.MockRepository) {
	t.Helper()

	repo := link.NewMockRepository(gomock.NewController(t))
	return link.NewService(repo), repo
}

func assertError(t *testing.T, err error, wantError bool) {
	t.Helper()

	if (err != nil) != wantError {
		t.Errorf("error = %v, wantError = %v", err, wantError)
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name      string
		req       *link.CreateLinkRequest
		want      *link.CreateResult
		wantError bool
		setupMock func(ctx context.Context, repo *link.MockRepository)
	}{
		{
			name:      "empty URL",
			req:       request(""),
			wantError: true,
		},
		{
			name: "new URL",
			req:  request(googleURL),
			want: &link.CreateResult{
				Created: true,
				Response: &link.LinkResponse{
					URL: googleFixedURL,
				},
			},
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByURL(ctx, googleFixedURL).
					Return(nil, apperr.ErrorLinkNotFound)
				repo.EXPECT().
					Create(ctx, googleFixedURL, gomock.Any()).
					Return(&link.Link{URL: googleFixedURL}, nil)
			},
		},
		{
			name: "existing URL",
			req:  request(googleURL),
			want: &link.CreateResult{
				Response: &link.LinkResponse{
					URL: googleFixedURL,
				},
			},
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByURL(ctx, googleFixedURL).
					Return(&link.Link{URL: googleFixedURL}, nil)
			},
		},
		{
			name:      "empty custom code",
			req:       requestWithCode(googleURL, ""),
			wantError: true,
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByURL(ctx, gomock.Any()).
					Return(nil, apperr.ErrorLinkNotFound)
			},
		},
		{
			name:      "long custom code",
			req:       requestWithCode(googleURL, "gooooooooooooooooogle"),
			wantError: true,
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByURL(ctx, gomock.Any()).
					Return(nil, apperr.ErrorLinkNotFound)
			},
		},
		{
			name:      "reserved custom code",
			req:       requestWithCode(googleURL, "shortenit"),
			wantError: true,
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByURL(ctx, gomock.Any()).
					Return(nil, apperr.ErrorLinkNotFound)
			},
		},
		{
			name: "valid custom code",
			req:  requestWithCode(googleURL, googleCode),
			want: &link.CreateResult{
				Created: true,
				Response: &link.LinkResponse{
					URL:  googleFixedURL,
					Code: googleCode,
				},
			},
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByURL(ctx, gomock.Any()).
					Return(nil, apperr.ErrorLinkNotFound)
				repo.EXPECT().
					ByCode(ctx, googleCode).
					Return(nil, apperr.ErrorLinkNotFound)
				repo.EXPECT().
					Create(ctx, googleFixedURL, gomock.Any()).
					Return(&link.Link{
						URL:  googleFixedURL,
						Code: googleCode,
					}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			service, repo := newService(t)

			if tt.setupMock != nil {
				tt.setupMock(ctx, repo)
			}

			got, err := service.Create(ctx, tt.req)

			assertError(t, err, tt.wantError)
			if tt.wantError {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateMany(t *testing.T) {
	tests := []struct {
		name      string
		reqs      []link.CreateLinkRequest
		want      []link.CreateResult
		setupMock func(ctx context.Context, repo *link.MockRepository)
	}{
		{
			name: "all succeed",
			reqs: []link.CreateLinkRequest{
				*requestWithCode(googleURL, googleCode),
				*requestWithCode(yandexURL, yandexCode),
			},
			want: []link.CreateResult{
				{
					Response: &link.LinkResponse{
						URL:  googleFixedURL,
						Code: googleCode,
					},
					Created: true,
				},
				{
					Response: &link.LinkResponse{
						URL:  yandexFixedURL,
						Code: yandexCode,
					},
					Created: true,
				},
			},
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByURL(ctx, gomock.Any()).
					Return(nil, apperr.ErrorLinkNotFound).
					AnyTimes()
				repo.EXPECT().
					ByCode(ctx, gomock.Any()).
					Return(nil, apperr.ErrorLinkNotFound).
					AnyTimes()
				repo.EXPECT().
					Create(ctx, googleFixedURL, gomock.Any()).
					Return(&link.Link{URL: googleFixedURL, Code: googleCode}, nil)
				repo.EXPECT().
					Create(ctx, yandexFixedURL, gomock.Any()).
					Return(&link.Link{URL: yandexFixedURL, Code: yandexCode}, nil)
			},
		},
		{
			name: "all failed",
			reqs: []link.CreateLinkRequest{
				*requestWithCode(googleURL, googleCode),
				*requestWithCode(yandexURL, yandexCode),
			},
			want: []link.CreateResult{
				{
					Response: nil,
					Created:  false,
					Error:    apperr.ErrorInvalidCustomCode,
				},
				{
					Response: nil,
					Created:  false,
					Error:    apperr.ErrorInvalidCustomCode,
				},
			},
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByURL(ctx, gomock.Any()).
					Return(nil, apperr.ErrorLinkNotFound).
					AnyTimes()
				repo.EXPECT().
					ByCode(ctx, gomock.Any()).
					Return(nil, apperr.ErrorLinkNotFound).
					AnyTimes()
				repo.EXPECT().
					Create(ctx, gomock.Any(), gomock.Any()).
					Return(nil, apperr.ErrorInvalidCustomCode).
					AnyTimes()
			},
		},
		{
			name: "only one succeed",
			reqs: []link.CreateLinkRequest{
				*requestWithCode(googleURL, googleCode),
				*requestWithCode(yandexURL, yandexCode),
			},
			want: []link.CreateResult{
				{
					Response: &link.LinkResponse{
						URL:  googleFixedURL,
						Code: googleCode,
					},
					Created: true,
				},
				{
					Response: nil,
					Created:  false,
					Error:    apperr.ErrorInvalidCustomCode,
				},
			},
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByURL(ctx, gomock.Any()).
					Return(nil, apperr.ErrorLinkNotFound).
					AnyTimes()
				repo.EXPECT().
					ByCode(ctx, gomock.Any()).
					Return(nil, apperr.ErrorLinkNotFound).
					AnyTimes()
				repo.EXPECT().
					Create(ctx, googleFixedURL, gomock.Any()).
					Return(&link.Link{URL: googleFixedURL, Code: googleCode}, nil)
				repo.EXPECT().
					Create(ctx, yandexFixedURL, gomock.Any()).
					Return(nil, apperr.ErrorInvalidCustomCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			service, repo := newService(t)
			tt.setupMock(ctx, repo)

			got := service.CreateMany(ctx, tt.reqs)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		want      *link.LinkResponse
		wantError bool
		setupMock func(ctx context.Context, repo *link.MockRepository)
	}{
		{
			name: "existing code",
			code: "AAAAAAAAAA",
			want: &link.LinkResponse{Code: "AAAAAAAAAA"},
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					Delete(ctx, "AAAAAAAAAA").
					Return(&link.Link{Code: "AAAAAAAAAA"}, nil)
			},
		},
		{
			name:      "non-existing code",
			code:      "AAAAAAAAAA",
			wantError: true,
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					Delete(ctx, "AAAAAAAAAA").
					Return(nil, errors.New("no link with such code"))
			},
		},
		{
			name:      "empty code",
			code:      "",
			wantError: true,
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					Delete(ctx, gomock.Any()).
					Return(nil, errors.New("no link with such code"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			service, repo := newService(t)
			tt.setupMock(ctx, repo)

			got, err := service.Delete(ctx, tt.code)

			assertError(t, err, tt.wantError)
			if tt.wantError {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestByCode(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		want      *link.LinkResponse
		wantError bool
		setupMock func(ctx context.Context, repo *link.MockRepository)
	}{
		{
			name: "existing code",
			code: "AAAAAAAAAA",
			want: &link.LinkResponse{Code: "AAAAAAAAAA"},
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByCode(ctx, "AAAAAAAAAA").
					Return(&link.Link{Code: "AAAAAAAAAA"}, nil)
			},
		},
		{
			name:      "non-existing code",
			code:      "AAAAAAAAAA",
			wantError: true,
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByCode(ctx, "AAAAAAAAAA").
					Return(nil, errors.New("no link with such code"))
			},
		},
		{
			name:      "empty code",
			code:      "",
			wantError: true,
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					ByCode(ctx, gomock.Any()).
					Return(nil, errors.New("no link with such code"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			service, repo := newService(t)
			tt.setupMock(ctx, repo)

			got, err := service.ByCode(ctx, tt.code)

			assertError(t, err, tt.wantError)
			if tt.wantError {
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name      string
		want      []link.LinkResponse
		setupMock func(ctx context.Context, repo *link.MockRepository)
	}{
		{
			name: "compare output arrays",
			want: []link.LinkResponse{
				{Code: "AAAAAAAA"},
				{Code: "BBBBBBBB"},
				{Code: "CCCCCCCC"},
			},
			setupMock: func(ctx context.Context, repo *link.MockRepository) {
				repo.EXPECT().
					All(ctx).
					Return([]link.Link{
						{Code: "AAAAAAAA"},
						{Code: "BBBBBBBB"},
						{Code: "CCCCCCCC"},
					}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			service, repo := newService(t)
			tt.setupMock(ctx, repo)

			got, err := service.All(ctx)

			assertError(t, err, false)
			assert.Equal(t, tt.want, got)
		})
	}
}
