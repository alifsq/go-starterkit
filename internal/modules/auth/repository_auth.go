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

	// JWT auth
	SaveRefreshToken(ctx context.Context, token *models.RefreshToken) (*models.RefreshToken, error)
	FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID int64) error
	DeleteExpiredTokens(ctx context.Context) error

	// Mail
	SaveEmailVerificationToken(ctx context.Context, token *models.EmailVerificationToken) error
	FindEmailVerificationToken(ctx context.Context, tokenHash string) (*models.EmailVerificationToken, error)
	MarkEmailVerificationTokenAsUsed(ctx context.Context, tokenID int64) error
	MarkUserEmailAsVerified(ctx context.Context, userID int64) error
	UpdateEmailVerificationToken(ctx context.Context, evtModel *models.EmailVerificationToken) error
	
	// Password Reset Token Methods
	SavePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error
	FindPasswordResetToken(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error)
	MarkPasswordResetTokenAsUsed(ctx context.Context, tokenID int64) error
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
	_, err := rt.WithContext(ctx).Where(rt.TokenHash.Eq(tokenHash)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("refresh token not found")
		}
		return err
	}

	_, err = rt.WithContext(ctx).Where(rt.TokenHash.Eq(tokenHash)).Update(rt.IsRevoked, true)
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

func (r *repositoryAuth) SaveEmailVerificationToken(ctx context.Context, token *models.EmailVerificationToken) error {
	evt := r.query.EmailVerificationToken
	return evt.WithContext(ctx).Create(token)
}

func (r *repositoryAuth) UpdateEmailVerificationToken(ctx context.Context, evtModel *models.EmailVerificationToken) error {
	evt := r.query.EmailVerificationToken
	_, err := evt.WithContext(ctx).Where(evt.UserID.Eq(evtModel.UserID)).Updates(evtModel)
	return err
}

func (r *repositoryAuth) FindEmailVerificationToken(ctx context.Context, tokenHash string) (*models.EmailVerificationToken, error) {
	evt := r.query.EmailVerificationToken
	return evt.WithContext(ctx).Preload(evt.User).Where(evt.TokenHash.Eq(tokenHash)).First()
}

func (r *repositoryAuth) MarkEmailVerificationTokenAsUsed(ctx context.Context, tokenID int64) error {
	evt := r.query.EmailVerificationToken
	now := time.Now()
	_, err := evt.WithContext(ctx).Where(evt.ID.Eq(tokenID)).Update(evt.UsedAt, &now)
	return err

}

func (r *repositoryAuth) MarkUserEmailAsVerified(ctx context.Context, userID int64) error {
	u := r.query.User
	now := time.Now()
	_, err := u.WithContext(ctx).Where(u.ID.Eq(userID)).UpdateSimple(u.EmailVerifiedAt.Value(now), u.IsEmailVerified.Value(true))
	return err
}

func (r *repositoryAuth) SavePasswordResetToken(ctx context.Context, token *models.PasswordResetToken) error {
	prt := r.query.PasswordResetToken
	return prt.WithContext(ctx).Create(token)
}

func (r *repositoryAuth) FindPasswordResetToken(ctx context.Context, tokenHash string) (*models.PasswordResetToken, error) {
	prt := r.query.PasswordResetToken
	return prt.WithContext(ctx).Preload(prt.User).Where(prt.TokenHash.Eq(tokenHash)).First()
}

func (r *repositoryAuth) MarkPasswordResetTokenAsUsed(ctx context.Context, tokenID int64) error {
	prt := r.query.PasswordResetToken
	now := time.Now()
	_, err := prt.WithContext(ctx).Where(prt.ID.Eq(tokenID)).Update(prt.UsedAt, &now)
	return err
}
