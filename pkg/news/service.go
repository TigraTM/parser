package news

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Service interface {
	CreateNews(ctx context.Context, news News) error
	GetNews(ctx context.Context, search string) ([]News, error)
}

type service struct {
	newsRepo Repository
	log      *logrus.Logger
}

func NewService(newsRepo Repository, log *logrus.Logger) Service {
	return &service{
		newsRepo: newsRepo,
		log:      log,
	}
}

func (s *service) CreateNews(ctx context.Context, news News) error {
	return s.newsRepo.CreateNews(ctx, news)
}

func (s *service) GetNews(ctx context.Context, search string) ([]News, error) {
	return s.newsRepo.GetNews(ctx, search)
}
