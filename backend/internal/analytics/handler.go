package analytics

import (
	"strconv"

	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetSalesReport(c *gin.Context) {
	role := c.GetString("role")
	
	var authBranchID *int
	if bID, exists := c.Get("branch_id"); exists && bID != nil {
		val := int(bID.(float64))
		authBranchID = &val
	}

	var reqBranchID *int
	if branchQ := c.Query("branch_id"); branchQ != "" {
		if b, err := strconv.Atoi(branchQ); err == nil {
			reqBranchID = &b
		}
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	results, err := h.service.GetSalesReport(c.Request.Context(), role, authBranchID, reqBranchID, startDate, endDate)
	if err != nil {
		if len(err.Error()) > 9 && err.Error()[:9] == "forbidden" {
			response.Error(c, 403, err.Error())
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, 200, results)
}

func (h *Handler) GetTopProducts(c *gin.Context) {
	role := c.GetString("role")
	
	var authBranchID *int
	if bID, exists := c.Get("branch_id"); exists && bID != nil {
		val := int(bID.(float64))
		authBranchID = &val
	}

	var reqBranchID *int
	if branchQ := c.Query("branch_id"); branchQ != "" {
		if b, err := strconv.Atoi(branchQ); err == nil {
			reqBranchID = &b
		}
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	
	limit := 10
	if limitQ := c.Query("limit"); limitQ != "" {
		if l, err := strconv.Atoi(limitQ); err == nil {
			limit = l
		}
	}

	results, err := h.service.GetTopProducts(c.Request.Context(), role, authBranchID, reqBranchID, startDate, endDate, limit)
	if err != nil {
		if len(err.Error()) > 9 && err.Error()[:9] == "forbidden" {
			response.Error(c, 403, err.Error())
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, 200, results)
}

func (h *Handler) GetShiftReports(c *gin.Context) {
	role := c.GetString("role")
	
	var authBranchID *int
	if bID, exists := c.Get("branch_id"); exists && bID != nil {
		val := int(bID.(float64))
		authBranchID = &val
	}

	var reqBranchID *int
	if branchQ := c.Query("branch_id"); branchQ != "" {
		if b, err := strconv.Atoi(branchQ); err == nil {
			reqBranchID = &b
		}
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	results, err := h.service.GetShiftReports(c.Request.Context(), role, authBranchID, reqBranchID, startDate, endDate)
	if err != nil {
		if len(err.Error()) > 9 && err.Error()[:9] == "forbidden" {
			response.Error(c, 403, err.Error())
		} else {
			response.Error(c, 500, err.Error())
		}
		return
	}

	response.Success(c, 200, results)
}
