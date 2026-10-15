package auth

import (
	"errors"
	"go-starterkit/internal/middleware"
	"go-starterkit/internal/response"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerAuth struct {
	service ServiceAuth
}

func NewHandlerAuth(s ServiceAuth) *HandlerAuth {
	return &HandlerAuth{
		service: s,
	}
}

const refreshCookiePath = "/api/v1/auth"

func setRefreshCookie(c *gin.Context, token string, maxAge int) {
	// secure := c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
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
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// 1. Ambil User-Agent & Client IP dari Gin Context
	userAgent := c.Request.UserAgent()
	ipAddress := c.ClientIP()

	// 2. Oper ke Service
	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req.Email, req.Password, userAgent, ipAddress)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	setRefreshCookie(c, refreshToken, 604800)
	response.Success(c, http.StatusOK, "Login success", gin.H{
		"access_token": accessToken,
	})
}
func (h *HandlerAuth) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil && refreshToken != "" {

		_ = h.service.Logout(c.Request.Context(), refreshToken)
	}

	// 2. Hapus Cookie di browser client dengan menyetel MaxAge = -1
	// Argumen: (name, value, maxAge, path, domain, secure, httpOnly)
	setRefreshCookie(c, "", -1)

	response.Success(c, http.StatusOK, "Logout Success", nil)
}

func (h *HandlerAuth) Register(c *gin.Context) {
	var req RegisterRequest

	// Validasi Input Payload/Body JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Payload request Not Found", err.Error())
		return
	}

	// Service Register
	_, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, response.ErrEmailExist) {
			response.Error(c, http.StatusConflict, "Email already registered", nil)
			return
		}

		response.Error(c, http.StatusInternalServerError, "Failed Registration User", nil)
		return
	}

	// Response Sucess HTTP 201 Created
	response.Success(c, http.StatusCreated, "Registrasi Success,Please login", nil)
}
func (h *HandlerAuth) Refresh(c *gin.Context) {
	// Cookie is preferred; JSON body is supported for API clients without a cookie jar.
	refreshTokenStr, err := readRefreshToken(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}
	if refreshTokenStr == "" {
		response.Error(c, http.StatusUnauthorized, "Refresh token Not Found", nil)
		return
	}

	// 2. Minta service melakukan verifikasi & rotation
	newAccessToken, newRefreshToken, err := h.service.RefreshToken(c.Request.Context(), refreshTokenStr)
	if err != nil {
		// Hapus cookie invalid jika token kedaluwarsa/di-revoke
		setRefreshCookie(c, "", -1)
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	// 3. Set HttpOnly Cookie baru untuk Refresh Token berikutnya
	setRefreshCookie(c, newRefreshToken, 604800)

	// 4. Kirim Access Token baru ke response JSON
	response.Success(c, http.StatusOK, "Refresh Token Success", gin.H{
		"access_token": newAccessToken,
	})
}

func (h *HandlerAuth) RevokeAll(c *gin.Context) {
	// Ambil langsung dari Gin Context
	val, exists := c.Get(middleware.CtxUserID)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User ID tidak ditemukan", nil)
		return
	}

	userID, ok := val.(int64)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Tipe data User ID tidak valid", nil)
		return
	}

	// Panggil service
	if err := h.service.LogoutAllDevices(c.Request.Context(), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "Gagal mencabut seluruh sesi", err.Error())
		return
	}

	setRefreshCookie(c, "", -1)
	response.Success(c, http.StatusOK, "Berhasil logout dari seluruh perangkat", nil)
}
