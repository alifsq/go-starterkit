package auth

import (
	"context"
	"errors"
	"time"

	"go-starterkit/internal/database/models"
	"go-starterkit/internal/database/query"
	"go-starterkit/internal/response"

	"gorm.io/gorm"
)

type RepositoryAuth interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) (*models.RefreshToken, error)
	FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID int64) error
	DeleteExpiredTokens(ctx context.Context) error
}

type repositoryAuth struct {
	db    *gorm.DB
	query *query.Query
}

func NewRepositoryAuth(query *query.Query, db *gorm.DB) RepositoryAuth {
	return &repositoryAuth{
		query: query,
		db:    db,
	}
}

func (r *repositoryAuth) Register(ctx context.Context, user *models.User) (*models.User, error) {
	u := r.query.User
	if err := u.WithContext(ctx).Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repositoryAuth) SaveRefreshToken(ctx context.Context, token *models.RefreshToken) (*models.RefreshToken, error) {
	rt := r.query.RefreshToken
	if token == nil {
		return nil, errors.New("refresh token is required")
	}
	if err := rt.WithContext(ctx).Create(token); err != nil {
		return nil, err
	}
	return token, nil
}

func (r *repositoryAuth) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	rt := r.query.RefreshToken
	tokenModel, err := rt.WithContext(ctx).Where(rt.TokenHash.Eq(tokenHash)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrUserNotFound
		}
		return nil, err
	}
	return tokenModel, nil
}

func (r *repositoryAuth) RevokeToken(ctx context.Context, tokenHash string) error {
	rt := r.query.RefreshToken
	_, err := rt.WithContext(ctx).Where(rt.TokenHash.Eq(tokenHash)).Update(rt.IsRevoked, true)
	return err
}

func (r *repositoryAuth) RevokeAllUserTokens(ctx context.Context, userID int64) error {
	rt := r.query.RefreshToken
	_, err := rt.WithContext(ctx).Where(rt.UserID.Eq(userID)).Update(rt.IsRevoked, true)
	return err
}

func (r *repositoryAuth) DeleteExpiredTokens(ctx context.Context) error {
	rt := r.query.RefreshToken
	_, err := rt.WithContext(ctx).Where(rt.ExpiresAt.Lte(time.Now())).Delete()
	return err
}
