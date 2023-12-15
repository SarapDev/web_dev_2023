package comment

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
	sqlx *sqlx.DB,
	tracer opentracing.Tracer,
) MySQLRepository {
	return MySQLRepository{db: sqlx, tracer: tracer}
}

func (m MySQLRepository) Create(ctx context.Context, comment *entity.Comment) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx, m.tracer,
		"app::common::repository::comment::mysql::create",
	)
	defer span.Finish()

	span.LogFields(log.Object("comment", *comment))

	result, err := m.db.ExecContext(
		ctx, `
			INSERT INTO comments (id, post_id, author, content, created_at)
			VALUES (?, ?, ?, ?, ?)
		`,
		comment.ID,
		comment.PostID,
		comment.Author,
		comment.Content,
		comment.CreatedAt,
	)

	if err != nil {
		return multierr.Append(
			errors.New("ошибка при создании комментария"),
			err,
		)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return multierr.Append(
			errors.New("ошибка при получении ID созданного комментария"),
			err,
		)
	}

	comment.ID = int(id)

	return nil
}

func (m MySQLRepository) ListByPostID(ctx context.Context, postID int) ([]*entity.Comment, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(
		ctx, m.tracer,
		"app::common::repository::comment::mysql::list_by_post_id",
	)
	defer span.Finish()

	span.LogFields(log.Int("post_id", postID))

	var comments []*entity.Comment

	err := m.db.SelectContext(
		ctx, &comments,
		`
			SELECT id, post_id, author, content, created_at
			FROM comments
			WHERE post_id = ? order by created_at desc
		`,
		postID,
	)

	if err != nil {
		return nil, multierr.Append(
			errors.New("ошибка при получении списка комментариев"),
			err,
		)
	}

	return comments, nil
}
