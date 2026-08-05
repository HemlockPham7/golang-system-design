package bookmark

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HemlockPham7/golang-system-design/internal/model"
	"github.com/HemlockPham7/golang-system-design/internal/repository/cache"
	"github.com/rs/zerolog/log"
)

const (
	getBookmarksCacheGroupKeyFormat = "get_bookmarks_%s"
	getBookmarksCacheKeyFormat      = "%d_%d"
	getBookmarksCacheExp            = 24 * time.Hour
)

type bookmarkServiceWithCache struct {
	s Service
	c cache.DB
}

func NewBookmarkServiceWithCache(s Service, c cache.DB) Service {
	return &bookmarkServiceWithCache{s: s, c: c}
}

func (s *bookmarkServiceWithCache) CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	cacheGroupKey := fmt.Sprintf(getBookmarksCacheGroupKeyFormat, userID)
	err := s.c.DeleteCache(ctx, cacheGroupKey)
	if err != nil {
		log.Err(err).Str("key", cacheGroupKey).Msg("Failed to delete cache")
	}

	return s.s.CreateBookmark(ctx, description, url, userID)
}

func (s *bookmarkServiceWithCache) UpdateBookmark(ctx context.Context, description, url, userID, ID string) (*model.Bookmark, error) {
	cacheGroupKey := fmt.Sprintf(getBookmarksCacheGroupKeyFormat, userID)
	err := s.c.DeleteCache(ctx, cacheGroupKey)
	if err != nil {
		log.Err(err).Str("key", cacheGroupKey).Msg("Failed to delete cache")
	}

	return s.s.UpdateBookmark(ctx, description, url, userID, ID)
}

func (s *bookmarkServiceWithCache) DeleteBookmark(ctx context.Context, userID, ID string) error {
	cacheGroupKey := fmt.Sprintf(getBookmarksCacheGroupKeyFormat, userID)
	err := s.c.DeleteCache(ctx, cacheGroupKey)
	if err != nil {
		log.Err(err).Str("key", cacheGroupKey).Msg("Failed to delete cache")
	}

	return s.s.DeleteBookmark(ctx, userID, ID)
}

func (s *bookmarkServiceWithCache) GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarksResult, error) {
	// tao cache key
	cacheGroupKey := fmt.Sprintf(getBookmarksCacheGroupKeyFormat, userID)
	cacheKey := fmt.Sprintf(getBookmarksCacheKeyFormat, page, limit)

	// get cache data
	cacheData, err := s.c.GetCacheData(ctx, cacheGroupKey, cacheKey)

	// if cache exits, return cache
	if err == nil && len(cacheData) > 0 {
		result := &GetBookmarksResult{}

		err := json.Unmarshal(cacheData, result)
		if err != nil {
			cacheErr := s.c.DeleteCache(ctx, cacheGroupKey)
			if cacheErr != nil {
				log.Err(cacheErr).Str("key", cacheGroupKey).Msg("Failed to delete cache")
			}
		} else {
			return result, nil
		}
	}

	// if not, call service
	result, err := s.s.GetBookmarks(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	// save cache
	resultBytes, err := json.Marshal(result)
	if err == nil {
		cacheErr := s.c.SetCacheData(ctx, cacheGroupKey, cacheKey, resultBytes, getBookmarksCacheExp)
		if cacheErr != nil {
			log.Err(cacheErr).Str("key", cacheGroupKey).Msg("failed to set cache")
		}
	}

	// return result
	return result, nil
}
