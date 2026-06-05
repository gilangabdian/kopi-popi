package notification

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gilangages/kopi-popi/pkg/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetMyNotifications(c *gin.Context) {
	userID := c.GetString("user_id") // Dari JWT Middleware

	notifications, err := h.service.GetMyNotifications(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, notifications)
}

func (h *Handler) MarkAsRead(c *gin.Context) {
	userID := c.GetString("user_id") // Dari JWT Middleware
	notifID := c.Param("id")

	if notifID == "" {
		response.Error(c, http.StatusBadRequest, "Notification ID is required")
		return
	}

	err := h.service.MarkAsRead(c.Request.Context(), userID, notifID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, nil)
}
