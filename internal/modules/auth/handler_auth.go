package auth

import (
	"errors"
	"go-starterkit/internal/package/constants"
	"go-starterkit/internal/package/logger"
	"go-starterkit/internal/response"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go-starterkit/internal/middleware"
	"go.uber.org/zap"
)

type HandlerAuth struct {
	service     ServiceAuth
	rateLimiter *middleware.RateLimitMiddleware
	log         logger.Logger
}

func NewHandlerAuth(s ServiceAuth, log logger.Logger, rateLimiter *middleware.RateLimitMiddleware) *HandlerAuth {
	return &HandlerAuth{
		service:     s,
		rateLimiter: rateLimiter,
		log:         log,
	}
}

const refreshCookiePath = "/api/v1/auth"

func setRefreshCookie(c *gin.Context, token string, maxAge int) {
	c.SetCookie("refresh_token", token, maxAge, refreshCookiePath, "", false, true)
}

func readRefreshToken(c *gin.Context) (string, error) {
	if token, err := c.Cookie("refresh_token"); err == nil && token != "" {
		return token, nil
	}

	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return req.RefreshToken, nil
}

func (h *HandlerAuth) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if formattedErrs := response.FormatValidationError(err); formattedErrs != nil {
			response.Error(c, http.StatusBadRequest, "Validasi input gagal", formattedErrs)
			return
		}
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req.Email, req.Password, userAgent, ipAddress)
	if err != nil {
		limit, limitErr := h.rateLimiter.CheckLimit(c.Request.Context(), "login_failure", strings.ToLower(strings.TrimSpace(req.Email)))
		if limitErr != nil {
			c.Error(limitErr)
			response.Error(c, http.StatusServiceUnavailable, "Rate limiting is temporarily unavailable", nil)
			return
		}
		if !limit.Allowed {
			retryAfter := limit.RetryAfter
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", strconv.FormatInt(retryAfter, 10))
			response.Error(c, http.StatusTooManyRequests, "Too many failed login attempts. Please try again later.", nil)
			return
		}
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	setRefreshCookie(c, refreshToken, 604800)
	response.Success(c, http.StatusOK, "Login success", gin.H{
		"access_token": accessToken,
	})
}

func (h *HandlerAuth) Logout(c *gin.Context) {
	refreshToken, err := readRefreshToken(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}
	if strings.TrimSpace(refreshToken) == "" {
		response.Error(c, http.StatusUnauthorized, "Refresh token not found", nil)
		return
	}
	if err := h.service.Logout(c.Request.Context(), refreshToken); err != nil {
		h.log.Error("Failed to revoke refresh token", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Logout failed", nil)
		return
	}

	setRefreshCookie(c, "", -1)
	response.Success(c, http.StatusOK, "Logout Success", nil)
}

func (h *HandlerAuth) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if formattedErrs := response.FormatValidationError(err); formattedErrs != nil {
			response.Error(c, http.StatusBadRequest, "Validasi input gagal", formattedErrs)
			return
		}
		response.Error(c, http.StatusBadRequest, "Payload request Not Found", err.Error())
		return
	}

	err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, response.ErrEmailExist) {
			response.Error(c, http.StatusConflict, "Email already registered", nil)
			return
		}
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Registrasi Success,Please login", nil)
}

func (h *HandlerAuth) Refresh(c *gin.Context) {
	refreshTokenStr, err := readRefreshToken(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}
	if refreshTokenStr == "" {
		response.Error(c, http.StatusUnauthorized, "Refresh token Not Found", nil)
		return
	}

	newAccessToken, newRefreshToken, err := h.service.RefreshToken(c.Request.Context(), refreshTokenStr)
	if err != nil {
		setRefreshCookie(c, "", -1)
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	setRefreshCookie(c, newRefreshToken, 604800)
	response.Success(c, http.StatusOK, "Refresh Token Success", gin.H{
		"access_token": newAccessToken,
	})
}

func (h *HandlerAuth) RevokeAll(c *gin.Context) {
	val, exists := c.Get(constants.CtxUserID.String())
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User ID tidak ditemukan", nil)
		return
	}

	userID, ok := val.(int64)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Tipe data User ID tidak valid", nil)
		return
	}

	if err := h.service.LogoutAllDevices(c.Request.Context(), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mencabut seluruh sesi", err.Error())
		return
	}

	setRefreshCookie(c, "", -1)
	response.Success(c, http.StatusOK, "Berhasil logout dari seluruh perangkat", nil)
}

func (h *HandlerAuth) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		response.Error(c, http.StatusBadRequest, "Verification Token Is Required", nil)
		return
	}

	if err := h.service.VerifyEmail(c, token); err != nil {
		h.log.Warn("Email Verification Failed", zap.Error(err))
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Email Verified Successfully You Can Now Login", nil)
}

func (h *HandlerAuth) ResendVerifyEmail(c *gin.Context) {

}

func (h *HandlerAuth) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid Forgot Password Request Failed", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "Invalid Request Payload", err.Error())
		return
	}
	limit, err := h.rateLimiter.CheckLimit(c.Request.Context(), "email_send_account", strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		c.Error(err)
		response.Error(c, http.StatusServiceUnavailable, "Rate limiting is temporarily unavailable", nil)
		return
	}
	if !limit.Allowed {
		response.Error(c, http.StatusOK, "If your email is registered, we have sent a password reset link to your email.", nil)
		return
	}
	if err := h.service.ForgotPassword(c, req.Email); err != nil {
		h.log.Error("Forgot Password Process Failed", zap.String("email", req.Email))
		return
	}
	response.Error(c, http.StatusOK, "If your email is registered, we have sent a password reset link to your email.", nil)
}

func (h *HandlerAuth) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid reset password request payload", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "Invalid Request Payload", err.Error())
		return
	}
	if err := h.service.ResetPassword(c, req.Token, req.NewPassword); err != nil {
		h.log.Warn("Reset password failed", zap.Error(err))
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Password Reset SuccessFully , you acn now Login", nil)
}

func (h *HandlerAuth) ResendVerificationEmail(c *gin.Context) {
	rawUserID, exist := c.Get("user_id")
	if !exist {
		response.Error(c, http.StatusUnauthorized, "Unauthorize : User Id Not Found", nil)
		return
	}

	// assert to int64
	userID, ok := rawUserID.(int64)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "User Format Not Valid", nil)
		return
	}

	// service
	err := h.service.ResendVerificationEmail(c, userID)
	if err != nil {
		h.log.Error("Resend Verification Email Failed", zap.Error(err))
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Send Email Verification Success", nil)
}
