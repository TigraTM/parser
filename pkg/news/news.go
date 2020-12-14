package news

import (
	"context"
)

type News struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"descriptions"`
	Link        string `json:"link"`
}

type Repository interface {
	CreateNews(ctx context.Context, news News) error
	GetNews(ctx context.Context, search string) ([]News, error)
}
