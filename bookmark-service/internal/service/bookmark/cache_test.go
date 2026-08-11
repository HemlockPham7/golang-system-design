package bookmark_test

import (
	"context"
	"testing"
	"time"

	"github.com/HemlockPham7/golang-system-design/internal/integration_test/data/fixtures"
	"github.com/HemlockPham7/golang-system-design/internal/model"
	mock_cache "github.com/HemlockPham7/golang-system-design/internal/repository/cache/mocks"
	"github.com/HemlockPham7/golang-system-design/internal/service/bookmark"
	mock_bookmark "github.com/HemlockPham7/golang-system-design/internal/service/bookmark/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkServiceWithCache_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupService func(ctx context.Context) *mock_bookmark.Service
		setupCache   func(ctx context.Context) *mock_cache.DB

		expectedResult *bookmark.GetBookmarksResult
		expectedError  error
	}{
		{
			name: "success - with cache",

			setupService: func(ctx context.Context) *mock_bookmark.Service {
				mockService := mock_bookmark.NewService(t)
				return mockService
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockCache := mock_cache.NewDB(t)
				mockCache.On("GetCacheData", ctx, "get_bookmarks_d7c13097-67a7-4eae-a60e-0b9b533b7bd4", "1_10").Return(
					[]byte(`{"bookmarks": [{"id": "d7c13097-67a7-4eae-a60e-0b9b533b7bd6", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z", "description": "Bookmark 1", "url": "https://www.google.com", "code": "bookmark1"}, {"id": "d7c13097-67a7-4eae-a60e-0b9b533b7bd7", "created_at": "2023-01-01T00:00:00Z", "updated_at": "2023-01-01T00:00:00Z", "description": "Bookmark 2", "url": "https://www.google.com", "code": "bookmark2"}], "total": 2}`),
					nil,
				)
				return mockCache
			},
			expectedResult: &bookmark.GetBookmarksResult{
				Bookmarks: []*model.Bookmark{
					{
						Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
						Description: "Bookmark 1",
						URL:         "https://www.google.com",
						Code:        "bookmark1",
					},
					{
						Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd7"),
						Description: "Bookmark 2",
						URL:         "https://www.google.com",
						Code:        "bookmark2",
					},
				},
				Total: 2,
			},
			expectedError: nil,
		},
		{
			name: "success - with no cache",

			setupService: func(ctx context.Context) *mock_bookmark.Service {
				mockService := mock_bookmark.NewService(t)
				mockService.On("GetBookmarks", ctx, "d7c13097-67a7-4eae-a60e-0b9b533b7bd4", 1, 10).Return(
					&bookmark.GetBookmarksResult{
						Bookmarks: []*model.Bookmark{
							{
								Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
								Description: "Bookmark 1",
								URL:         "https://www.google.com",
								Code:        "bookmark1",
							},
							{
								Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd7"),
								Description: "Bookmark 2",
								URL:         "https://www.google.com",
								Code:        "bookmark2",
							},
						},
						Total: 2,
					},
					nil,
				)
				return mockService
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockCache := mock_cache.NewDB(t)
				mockCache.On("GetCacheData", ctx, "get_bookmarks_d7c13097-67a7-4eae-a60e-0b9b533b7bd4", "1_10").Return(nil, redis.Nil)
				mockCache.On("SetCacheData",
					ctx,
					"get_bookmarks_d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
					"1_10",
					[]byte(`{"bookmarks":[{"id":"d7c13097-67a7-4eae-a60e-0b9b533b7bd6","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z","description":"Bookmark 1","url":"https://www.google.com","code":"bookmark1","user_id":""},{"id":"d7c13097-67a7-4eae-a60e-0b9b533b7bd7","created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z","description":"Bookmark 2","url":"https://www.google.com","code":"bookmark2","user_id":""}],"total":2}`),
					24*time.Hour).Return(nil)
				return mockCache
			},
			expectedResult: &bookmark.GetBookmarksResult{
				Bookmarks: []*model.Bookmark{
					{
						Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
						Description: "Bookmark 1",
						URL:         "https://www.google.com",
						Code:        "bookmark1",
					},
					{
						Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd7"),
						Description: "Bookmark 2",
						URL:         "https://www.google.com",
						Code:        "bookmark2",
					},
				},
				Total: 2,
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockService := tc.setupService(ctx)
			mockCache := tc.setupCache(ctx)

			testService := bookmark.NewBookmarkServiceWithCache(mockService, mockCache)
			res, err := testService.GetBookmarks(ctx, "d7c13097-67a7-4eae-a60e-0b9b533b7bd4", 1, 10)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

func TestBookmarkServiceWithCache_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupService func(ctx context.Context) *mock_bookmark.Service
		setupCache   func(ctx context.Context) *mock_cache.DB

		expectedResult *model.Bookmark
		expectedError  error
	}{
		{
			name: "create bookmark successfully",

			setupService: func(ctx context.Context) *mock_bookmark.Service {
				mockService := mock_bookmark.NewService(t)
				mockService.On("CreateBookmark", ctx, "abc xyz", "https://www.youtube.com", "d7c13097-a60e-67a7-4eae-0b9b533b7bd4").Return(&model.Bookmark{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
					Description: "abc xyz",
					URL:         "https://www.youtube.com",
					Code:        "abcxyz",
					UserID:      "d7c13097-a60e-67a7-4eae-0b9b533b7bd4",
				}, nil)
				return mockService
			},

			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockCache := mock_cache.NewDB(t)
				mockCache.On("DeleteCache", ctx, "get_bookmarks_d7c13097-a60e-67a7-4eae-0b9b533b7bd4").Return(nil)
				return mockCache
			},

			expectedResult: &model.Bookmark{
				Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
				Description: "abc xyz",
				URL:         "https://www.youtube.com",
				Code:        "abcxyz",
				UserID:      "d7c13097-a60e-67a7-4eae-0b9b533b7bd4",
			},
			expectedError: nil,
		},
		{
			name:           "fail to create bookmark",
			expectedResult: nil,
			expectedError:  assert.AnError,

			setupService: func(ctx context.Context) *mock_bookmark.Service {
				mockService := mock_bookmark.NewService(t)
				mockService.On("CreateBookmark", ctx, "abc xyz", "https://www.youtube.com", "d7c13097-a60e-67a7-4eae-0b9b533b7bd4").Return(nil, assert.AnError)
				return mockService
			},

			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockCache := mock_cache.NewDB(t)
				mockCache.On("DeleteCache", ctx, "get_bookmarks_d7c13097-a60e-67a7-4eae-0b9b533b7bd4").Return(assert.AnError)
				return mockCache
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockService := tc.setupService(ctx)
			mockCache := tc.setupCache(ctx)

			testService := bookmark.NewBookmarkServiceWithCache(mockService, mockCache)
			bookmarkResult, err := testService.CreateBookmark(ctx, "abc xyz", "https://www.youtube.com", "d7c13097-a60e-67a7-4eae-0b9b533b7bd4")

			assert.Equal(t, tc.expectedResult, bookmarkResult)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestBookmarkServiceWithCache_UpdateBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupService func(ctx context.Context) *mock_bookmark.Service
		setupCache   func(ctx context.Context) *mock_cache.DB

		expectedResult *model.Bookmark
		expectedError  error
	}{
		{
			name: "update bookmark successfully",

			setupService: func(ctx context.Context) *mock_bookmark.Service {
				mockService := mock_bookmark.NewService(t)
				expectedBookmark := &model.Bookmark{
					Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
					Description: "123 456",
					URL:         "https://www.facebook.com",
					Code:        "123456",
					UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
				}
				mockService.On("UpdateBookmarkByID", ctx, "123 456", "https://www.facebook.com", "d7c13097-67a7-4eae-a60e-0b9b533b7bd4", "d7c13097-67a7-4eae-a60e-0b9b533b7bd6").Return(expectedBookmark, nil)
				return mockService
			},

			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockCache := mock_cache.NewDB(t)
				mockCache.On("DeleteCache", ctx, "get_bookmarks_d7c13097-67a7-4eae-a60e-0b9b533b7bd4").Return(nil)
				return mockCache
			},

			expectedResult: &model.Bookmark{
				Base:        fixtures.GetTestBase("d7c13097-67a7-4eae-a60e-0b9b533b7bd6"),
				Description: "123 456",
				URL:         "https://www.facebook.com",
				Code:        "123456",
				UserID:      "d7c13097-67a7-4eae-a60e-0b9b533b7bd4",
			},
			expectedError: nil,
		},
		{
			name: "fail to update bookmark",

			setupService: func(ctx context.Context) *mock_bookmark.Service {
				mockService := mock_bookmark.NewService(t)
				mockService.On("UpdateBookmarkByID", ctx, "123 456", "https://www.facebook.com", "d7c13097-67a7-4eae-a60e-0b9b533b7bd4", "d7c13097-67a7-4eae-a60e-0b9b533b7bd6").Return(nil, assert.AnError)
				return mockService
			},

			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockCache := mock_cache.NewDB(t)
				mockCache.On("DeleteCache", ctx, "get_bookmarks_d7c13097-67a7-4eae-a60e-0b9b533b7bd4").Return(assert.AnError)
				return mockCache
			},

			expectedResult: nil,
			expectedError:  assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockService := tc.setupService(ctx)
			mockCache := tc.setupCache(ctx)

			testService := bookmark.NewBookmarkServiceWithCache(mockService, mockCache)
			bookmarkResult, err := testService.UpdateBookmarkByID(ctx, "123 456", "https://www.facebook.com", "d7c13097-67a7-4eae-a60e-0b9b533b7bd4", "d7c13097-67a7-4eae-a60e-0b9b533b7bd6")

			assert.Equal(t, tc.expectedResult, bookmarkResult)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestBookmarkService_DeleteBookmarkByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupService func(ctx context.Context) *mock_bookmark.Service
		setupCache   func(ctx context.Context) *mock_cache.DB

		expectedResult *model.Bookmark
		expectedError  error
	}{
		{
			name: "delete bookmark successfully",

			setupService: func(ctx context.Context) *mock_bookmark.Service {
				mockService := mock_bookmark.NewService(t)
				mockService.On("DeleteBookmarkByID", ctx, "d7c13097-67a7-4eae-a60e-0b9b533b7bd4", "d7c13097-67a7-4eae-a60e-0b9b533b7bd6").Return(nil)
				return mockService
			},

			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockCache := mock_cache.NewDB(t)
				mockCache.On("DeleteCache", ctx, "get_bookmarks_d7c13097-67a7-4eae-a60e-0b9b533b7bd4").Return(nil)
				return mockCache
			},

			expectedError: nil,
		},
		{
			name: "fail to delete bookmark",

			setupService: func(ctx context.Context) *mock_bookmark.Service {
				mockService := mock_bookmark.NewService(t)
				mockService.On("DeleteBookmarkByID", ctx, "d7c13097-67a7-4eae-a60e-0b9b533b7bd4", "d7c13097-67a7-4eae-a60e-0b9b533b7bd6").Return(assert.AnError)
				return mockService
			},

			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockCache := mock_cache.NewDB(t)
				mockCache.On("DeleteCache", ctx, "get_bookmarks_d7c13097-67a7-4eae-a60e-0b9b533b7bd4").Return(assert.AnError)
				return mockCache
			},

			expectedError: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockService := tc.setupService(ctx)
			mockCache := tc.setupCache(ctx)

			testService := bookmark.NewBookmarkServiceWithCache(mockService, mockCache)
			err := testService.DeleteBookmarkByID(ctx, "d7c13097-67a7-4eae-a60e-0b9b533b7bd4", "d7c13097-67a7-4eae-a60e-0b9b533b7bd6")

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
