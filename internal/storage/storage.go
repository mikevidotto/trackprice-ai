package storage

import (
	"context"

	"github.com/mikevidotto/blogfolio/internal/article"
)

type Storage interface {
	Create(ctx context.Context, a article.Article) (article.Article, error)
	Read(ctx context.Context) ([]article.Article, error)
	Update(ctx context.Context, a article.Article) (article.Article, error)
	Delete(ctx context.Context, id int64) error
}
