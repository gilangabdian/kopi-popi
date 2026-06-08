package promo

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gilangages/kopi-popi/pkg/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) CreatePromo(c *gin.Context) {
	role := c.GetString("role")
	if role != "Admin" {
		response.Error(c, 403, "forbidden: only admin can create promos")
		return
	}

	var req PromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	err := h.service.CreatePromo(c.Request.Context(), req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 201, "promo created successfully")
}

func (h *Handler) UpdatePromo(c *gin.Context) {
	role := c.GetString("role")
	if role != "Admin" {
		response.Error(c, 403, "forbidden: only admin can update promos")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, 400, "invalid promo id")
		return
	}

	var req PromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	err = h.service.UpdatePromo(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "promo updated successfully")
}

func (h *Handler) GetPromos(c *gin.Context) {
	role := c.GetString("role")

	promos, err := h.service.GetPromos(c.Request.Context(), role)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, promos)
}

func (h *Handler) ValidatePromo(c *gin.Context) {
	var req ValidatePromoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	resp, err := h.service.ValidatePromo(c.Request.Context(), req.Code, req.TotalAmount)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, resp)
}
