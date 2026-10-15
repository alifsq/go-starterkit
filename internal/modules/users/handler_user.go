package user

import (
	"go-starterkit/internal/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerUser struct {
	service ServiceUser
}

func NewHandlerUser(s ServiceUser) *HandlerUser {
	return &HandlerUser{
		service: s,
	}
}

func (h *HandlerUser) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	// Bind JSON & Validasi Tag
	if err := c.ShouldBindJSON(&req); err != nil {
		// Helper validator
		if formattedErrs := response.FormatValidationError(err); formattedErrs != nil {
			response.Error(c, http.StatusBadRequest, "Validasi input fail", formattedErrs)
			return
		}

		// Helper error
		response.Error(c, http.StatusBadRequest, "Invalid Request Payload", err.Error())
		return
	}

	// Validasi costum modules user


	// Service User
	res, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// response helper success
	response.Success(c, http.StatusCreated, "User Created Successfully", res)
}
