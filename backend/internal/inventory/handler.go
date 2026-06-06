package inventory

import (
	"net/http"
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

func (h *Handler) GetBranchStock(c *gin.Context) {
	branchID, err := strconv.Atoi(c.Param("branch_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid branch ID")
		return
	}

	role := c.GetString("role")
	var reqBranchID *int
	if bID, exists := c.Get("branch_id"); exists && bID != nil {
		val := int(bID.(float64)) // JWT standard decodes numbers as float64
		reqBranchID = &val
	}

	stocks, err := h.service.GetBranchStock(branchID, role, reqBranchID)
	if err != nil {
		if err.Error() == "forbidden: can only access your own branch stock" {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, stocks)
}

func (h *Handler) GetInventoryMovements(c *gin.Context) {
	branchIDStr := c.Query("branch_id")
	if branchIDStr == "" {
		response.Error(c, http.StatusBadRequest, "branch_id query parameter is required")
		return
	}
	
	branchID, err := strconv.Atoi(branchIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid branch_id query parameter")
		return
	}

	role := c.GetString("role")
	var reqBranchID *int
	if bID, exists := c.Get("branch_id"); exists && bID != nil {
		val := int(bID.(float64))
		reqBranchID = &val
	}

	movements, err := h.service.GetInventoryMovements(branchID, role, reqBranchID)
	if err != nil {
		if err.Error() == "forbidden: can only access your own branch movements" {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, movements)
}

func (h *Handler) GetRestockRequests(c *gin.Context) {
	role := c.GetString("role")
	var reqBranchID *int
	if bID, exists := c.Get("branch_id"); exists && bID != nil {
		val := int(bID.(float64))
		reqBranchID = &val
	}

	requests, err := h.service.GetRestockRequests(role, reqBranchID)
	if err != nil {
		response.Error(c, http.StatusForbidden, err.Error())
		return
	}

	response.Success(c, http.StatusOK, requests)
}

type CreateRestockRequestPayload struct {
	Reason string `json:"reason"`
	Items  []struct {
		MaterialID        int     `json:"material_id" binding:"required"`
		QuantityRequested float64 `json:"quantity_requested" binding:"required,gt=0"`
	} `json:"items" binding:"required,min=1"`
}

func (h *Handler) CreateRestockRequest(c *gin.Context) {
	var payload CreateRestockRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID := c.GetString("user_id")
	role := c.GetString("role")
	
	var reqBranchID *int
	if bID, exists := c.Get("branch_id"); exists && bID != nil {
		val := int(bID.(float64))
		reqBranchID = &val
	}

	req := RestockRequest{
		RequestedBy: userID,
		Reason:      payload.Reason,
	}
	for _, item := range payload.Items {
		req.Items = append(req.Items, RestockItem{
			MaterialID:        item.MaterialID,
			QuantityRequested: item.QuantityRequested,
		})
	}

	err := h.service.CreateRestockRequest(&req, role, reqBranchID)
	if err != nil {
		if err.Error() == "forbidden: only manager can create restock requests" || err.Error() == "forbidden: branch ID missing" {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, req)
}

type UpdateRestockStatusPayload struct {
	Status          string  `json:"status" binding:"required,oneof=Approved Rejected Delivered"`
	RejectionReason *string `json:"rejection_reason,omitempty"`
}

func (h *Handler) UpdateRestockStatus(c *gin.Context) {
	id := c.Param("id")
	var payload UpdateRestockStatusPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	role := c.GetString("role")
	var reqBranchID *int
	if bID, exists := c.Get("branch_id"); exists && bID != nil {
		val := int(bID.(float64))
		reqBranchID = &val
	}

	err := h.service.UpdateRestockStatus(id, payload.Status, payload.RejectionReason, role, reqBranchID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "not found: restock request not found" {
			statusCode = http.StatusNotFound
		} else if len(err.Error()) > 9 && err.Error()[:9] == "forbidden" {
			statusCode = http.StatusForbidden
		} else if len(err.Error()) > 8 && err.Error()[:8] == "conflict" {
			statusCode = http.StatusConflict
		}
		
		response.Error(c, statusCode, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "status updated successfully"})
}

func (h *Handler) ReceiveIncomingStock(c *gin.Context) {
	var payload ReceiveIncomingStockPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	role := c.GetString("role")
	if err := h.service.ReceiveIncomingStock(payload, role); err != nil {
		if err.Error() == "forbidden: only admin can receive incoming stock" {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "stok berhasil diterima dan dimasukkan ke gudang pusat"})
}

func (h *Handler) AllocateStock(c *gin.Context) {
	branchIDStr := c.Param("branch_id")
	branchID, err := strconv.Atoi(branchIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "branch_id tidak valid")
		return
	}

	var payload AllocateStockPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	role := c.GetString("role")
	if err := h.service.AllocateStock(branchID, payload, role); err != nil {
		if len(err.Error()) > 9 && err.Error()[:9] == "forbidden" {
			response.Error(c, http.StatusForbidden, err.Error())
			return
		} else if len(err.Error()) > 8 && err.Error()[:8] == "conflict" {
			response.Error(c, http.StatusConflict, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "stok berhasil dialokasikan ke cabang"})
}
