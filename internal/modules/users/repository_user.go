package user

import (
	"context"
	"errors"
	"go-starterkit/internal/database/models"
	"go-starterkit/internal/database/query"
	"go-starterkit/internal/response"

	"gorm.io/gorm"
)

type RepositoryUser interface {
	CreateUser(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id int64) (*models.User, error) 
}

type repositoryUser struct {
	db    *gorm.DB
	query *query.Query
}

func NewRepositoryUser(query *query.Query, db *gorm.DB) RepositoryUser {
	return &repositoryUser{
		query: query,
		db:    db,
	}
}

func (r *repositoryUser) CreateUser(ctx context.Context, user *models.User) error {
	return r.query.User.WithContext(ctx).Create(user)
}

func (r *repositoryUser) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u := r.query.User
	result, err := u.WithContext(ctx).Where(u.Email.Eq(email)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrUserNotFound
		}
		return nil, err
	}
	return result, nil
}

func (r *repositoryUser) FindByID(ctx context.Context, id int64) (*models.User, error) {
	u := r.query.User
	result, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}
	return result, nil
}
