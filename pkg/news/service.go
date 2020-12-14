package news

import (
	"context"
)

type Service interface {
	CreateNews(ctx context.Context, news News) error
	GetNews(ctx context.Context, search string) ([]News, error)
}

type service struct {
	newsRepo Repository
}

func NewService(newsRepo Repository) Service {
	return &service{
		newsRepo: newsRepo,
	}
}

func (s *service) CreateNews(ctx context.Context, news News) error {
	return s.newsRepo.CreateNews(ctx, news)
}

func (s *service) GetNews(ctx context.Context, search string) ([]News, error) {
	return s.newsRepo.GetNews(ctx, search)
}