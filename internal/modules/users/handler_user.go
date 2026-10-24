package user

import (
	"go-starterkit/internal/modules/users/dto"
	"go-starterkit/internal/modules/users/mapper"
	"go-starterkit/internal/package/constants"
	"go-starterkit/internal/response"
	"net/http"
	"strconv"

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
	var req dto.CreateUserRequest

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

func (h *HandlerUser) GetUser(c *gin.Context) {
	userIDVal, exist := c.Get(constants.CtxUserID.String())
	if !exist {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	userID, ok := userIDVal.(int64)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "error : Invalid user ID type in context", nil)
		return
	}

	profile, err := h.service.GetMyProfile(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "error: User profile not found", err.Error())
		return
	}

	myProfile := mapper.ToMyProfileResponse(profile)

	response.Success(c, http.StatusOK, "Success Get Profile", myProfile)
}

func (h *HandlerUser) GetAllUsersPaginate(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	users, total, err := h.service.GetAllUsersPaginate(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "error : Failed to fetch paginated users", err)
		return
	}
	// total page size
	totalPages := (int(total) + pageSize - 1) / pageSize

	// Meta Object
	meta := gin.H{
		"current_page": page,
		"page_size":    pageSize,
		"total_data":   total,
		"total_pages":  totalPages,
	}

	response.SuccessWithMeta(c, http.StatusOK, "Success Get All Users", users, meta)
}
