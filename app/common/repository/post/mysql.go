package post

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.kolesa-team.org/backend/go-example/app/common/entity"
	"go.uber.org/multierr"
)

var _ Repository = MySQLRepository{}

type MySQLRepository struct {
	db     *sqlx.DB
	tracer opentracing.Tracer
}

func NewMySQLRepository(
	db *sqlx.DB,
	tracer opentracing.Tracer,
) *MySQLRepository {
	return &MySQLRepository{db: db, tracer: tracer}
}

func (m MySQLRepository) Create(ctx context.Context, post *entity.Post) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx, m.tracer,
		"app::common::repository::post::mysql::create",
	)
	defer span.Finish()

	span.LogFields(log.Object("post", *post))

	result, err := m.db.ExecContext(
		ctx, `
			INSERT INTO posts (id, title, content, created_at)
			VALUES (?, ?, ?, ?)
		`,
		post.ID,
		post.Title,
		post.Content,
		post.CreatedAt,
	)

	if err != nil {
		return multierr.Append(
			errors.New("ошибка при создании поста"),
			err,
		)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return multierr.Append(
			errors.New("ошибка при получении ID созданного поста"),
			err,
		)
	}

	post.ID = int(id)

	return nil
}

func (m MySQLRepository) Get(ctx context.Context, id int) (*entity.Post, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx, m.tracer,
		"app::common::repository::post::mysql::get",
	)
	defer span.Finish()

	span.LogFields(log.Int("id", id))

	post := entity.Post{}

	err := m.db.GetContext(
		ctx, &post,
		`
			SELECT id, title, content, created_at
			FROM posts
			WHERE id = ? LIMIT 1
		`,
		id,
	)

	if err != nil {
		return nil, multierr.Append(
			errors.New("ошибка при получении поста"),
			err,
		)
	}

	return &post, nil
}

func (m MySQLRepository) List(ctx context.Context, offset, limit int) ([]*entity.Post, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx, m.tracer,
		"app::common::repository::post::mysql::list",
	)
	defer span.Finish()

	span.LogFields(
		log.Int("offset", offset),
		log.Int("limit", limit),
	)

	posts := make([]*entity.Post, 0, limit)

	err := m.db.SelectContext(
		ctx, &posts,
		`
			SELECT id, title, content, created_at
			FROM posts
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
		`,
		limit,
		offset,
	)

	if err != nil {
		return nil, multierr.Append(
			errors.New("ошибка при получении списка постов"),
			err,
		)
	}

	return posts, nil
}

func (m MySQLRepository) Count(ctx context.Context) (int, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx, m.tracer,
		"app::common::repository::post::mysql::count",
	)
	defer span.Finish()

	var count int

	err := m.db.GetContext(
		ctx, &count,
		`
			SELECT COUNT(id)
			FROM posts
		`,
	)

	if err != nil {
		return 0, multierr.Append(
			errors.New("ошибка при получении количества постов"),
			err,
		)
	}

	return count, nil
}
