package news

import (
	"context"
)

type Service interface {
	CreateNews(ctx context.Context, news News) error
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
	return nil
}