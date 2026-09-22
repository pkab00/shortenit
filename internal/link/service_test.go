package link_test

import (
	"errors"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/pkab00/shortenit/internal/apperr"
	"github.com/pkab00/shortenit/internal/link"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -destination=repository_test_mock.go -source=repository.go -package=link

func TestCreate(t *testing.T) {
	tests := []struct {
		name      string
		req       *link.CreateLinkRequest
		want      *link.CreateResult
		wantError bool
		setupMock func(repo *link.MockRepository)
	}{
		{
			name:      "empty url",
			req:       &link.CreateLinkRequest{URL: ""},
			want:      nil,
			wantError: true,
			setupMock: func(repo *link.MockRepository) {},
		},
		{
			name: "valid url",
			req:  &link.CreateLinkRequest{URL: "www.google.com"},
			want: &link.CreateResult{
				Created: true,
				Link: &link.LinkResponse{
					URL: "https://www.google.com",
				},
			},
			wantError: false,
			setupMock: func(repo *link.MockRepository) {
				repo.EXPECT().
					ByURL(t.Context(), "https://www.google.com").
					Return(nil, apperr.ErrorLinkNotFound)
				repo.EXPECT().
					Create(t.Context(), "https://www.google.com", gomock.Any()).
					Return(&link.Link{URL: "https://www.google.com"}, nil)
			},
		},
		{
			name: "existing url",
			req:  &link.CreateLinkRequest{URL: "www.google.com"},
			want: &link.CreateResult{
				Created: false,
				Link: &link.LinkResponse{
					URL: "https://www.google.com",
				},
			},
			wantError: false,
			setupMock: func(repo *link.MockRepository) {
				repo.EXPECT().
					ByURL(t.Context(), "https://www.google.com").
					Return(&link.Link{URL: "https://www.google.com"}, nil)
			},
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := link.NewMockRepository(ctrl)
			serv := link.NewService(repo)

			testcase.setupMock(repo)

			got, err := serv.Create(t.Context(), testcase.req)

			if testcase.wantError {
				if err == nil {
					t.Errorf("wanted an error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			assert.Equal(t, testcase.want, got)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		want      *link.LinkResponse
		wantError bool
		setupMock func(repo *link.MockRepository)
	}{
		{
			name:      "existing code",
			code:      "AAAAAAAAAA",
			want:      &link.LinkResponse{Code: "AAAAAAAAAA"},
			wantError: false,
			setupMock: func(repo *link.MockRepository) {
				repo.EXPECT().
					Delete(t.Context(), "AAAAAAAAAA").
					Return(&link.Link{Code: "AAAAAAAAAA"}, nil)
			},
		},
		{
			name:      "non-existing code",
			code:      "AAAAAAAAAA",
			want:      nil,
			wantError: true,
			setupMock: func(repo *link.MockRepository) {
				repo.EXPECT().
					Delete(t.Context(), "AAAAAAAAAA").
					Return(nil, errors.New("no link with such code"))
			},
		},
		{
			name:      "no code given (empty)",
			code:      "",
			want:      nil,
			wantError: true,
			setupMock: func(repo *link.MockRepository) {
				repo.EXPECT().
					Delete(t.Context(), gomock.Any()).
					Return(nil, errors.New("no link with such code"))
			},
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := link.NewMockRepository(ctrl)
			serv := link.NewService(repo)

			testcase.setupMock(repo)

			got, err := serv.Delete(t.Context(), testcase.code)

			if testcase.wantError {
				if err == nil {
					t.Errorf("wanted an error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			assert.Equal(t, testcase.want, got)
		})
	}
}

func TestByCode(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		want      *link.LinkResponse
		wantError bool
		setupMock func(repo *link.MockRepository)
	}{
		{
			name:      "existing code",
			code:      "AAAAAAAAAA",
			want:      &link.LinkResponse{Code: "AAAAAAAAAA"},
			wantError: false,
			setupMock: func(repo *link.MockRepository) {
				repo.EXPECT().
					ByCode(t.Context(), "AAAAAAAAAA").
					Return(&link.Link{Code: "AAAAAAAAAA"}, nil)
			},
		},
		{
			name:      "non-existing code",
			code:      "AAAAAAAAAA",
			want:      nil,
			wantError: true,
			setupMock: func(repo *link.MockRepository) {
				repo.EXPECT().
					ByCode(t.Context(), "AAAAAAAAAA").
					Return(nil, errors.New("no link with such code"))
			},
		},
		{
			name:      "no code given (empty)",
			code:      "",
			want:      nil,
			wantError: true,
			setupMock: func(repo *link.MockRepository) {
				repo.EXPECT().
					ByCode(t.Context(), gomock.Any()).
					Return(nil, errors.New("no link with such code"))
			},
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := link.NewMockRepository(ctrl)
			serv := link.NewService(repo)

			testcase.setupMock(repo)

			got, err := serv.ByCode(t.Context(), testcase.code)

			if testcase.wantError {
				if err == nil {
					t.Errorf("wanted an error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			assert.Equal(t, testcase.want, got)
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name      string
		want      []link.LinkResponse
		wantError bool
		setupMock func(repo *link.MockRepository)
	}{
		{
			name: "compare output arrays",
			want: []link.LinkResponse{
				{Code: "AAAAAAAA"},
				{Code: "BBBBBBBB"},
				{Code: "CCCCCCCC"},
			},
			wantError: false,
			setupMock: func(repo *link.MockRepository) {
				repo.EXPECT().
					All(t.Context()).
					Return([]link.Link{
						{Code: "AAAAAAAA"},
						{Code: "BBBBBBBB"},
						{Code: "CCCCCCCC"},
					}, nil)
			},
		},
	}

	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := link.NewMockRepository(ctrl)
			serv := link.NewService(repo)

			testcase.setupMock(repo)

			got, err := serv.All(t.Context())

			if testcase.wantError {
				if err == nil {
					t.Errorf("wanted an error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			assert.Equal(t, testcase.want, got)
		})
	}
}
