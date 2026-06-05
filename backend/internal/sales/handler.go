package sales

import (
	"fmt"
	"strconv"

	"github.com/gilangages/kopi-popi/internal/payment"
	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service    Service
	paymentSvc payment.Service
}

func NewHandler(service Service, paymentSvc payment.Service) *Handler {
	return &Handler{service, paymentSvc}
}

// --- SHIFTS ---

func (h *Handler) OpenShift(c *gin.Context) {
	role := c.GetString("role")
	if role != "Cashier" {
		response.Error(c, 403, "forbidden: only cashier can open shift")
		return
	}

	userID := c.GetString("user_id")
	branchIDFloat, exists := c.Get("branch_id")
	if !exists {
		response.Error(c, 403, "forbidden: no branch assigned")
		return
	}
	branchID := int(branchIDFloat.(float64))

	var req OpenShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	shift, err := h.service.OpenShift(c.Request.Context(), branchID, userID, req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, shift)
}

func (h *Handler) CloseShift(c *gin.Context) {
	role := c.GetString("role")
	if role != "Cashier" {
		response.Error(c, 403, "forbidden: only cashier can close shift")
		return
	}

	userID := c.GetString("user_id")

	var req CloseShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	err := h.service.CloseShift(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "shift closed successfully")
}

func (h *Handler) GetMyOpenShift(c *gin.Context) {
	role := c.GetString("role")
	if role != "Cashier" {
		response.Error(c, 403, "forbidden: only cashier has shifts")
		return
	}

	userID := c.GetString("user_id")

	shift, err := h.service.GetMyOpenShift(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, shift) // can be null if no open shift
}

// --- CARTS ---

func (h *Handler) AddCartItem(c *gin.Context) {
	var req AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	// Butuh branch ID, dari mana? 
	// Kasir = dari JWT claim (mereka masukin barang ke keranjang offline)
	// Customer = dari mana? Harus dilempar dari body/query karena customer gak terikat cabang spesifik.
	// Oh, I forgot to add BranchID to AddCartItemRequest for Customer?
	// But it's okay, let's extract it from query params for Customer, or just assume it's in a header/query.
	// We'll read from Query for Customer.
	
	role := c.GetString("role")
	var customerID *string
	var branchID int

	if role == "Customer" {
		uID := c.GetString("user_id")
		customerID = &uID
		bIDStr := c.Query("branch_id")
		if bIDStr == "" {
			response.Error(c, 400, "invalid: branch_id query parameter is required for customer")
			return
		}
		b, err := strconv.Atoi(bIDStr)
		if err != nil {
			response.Error(c, 400, "invalid branch_id")
			return
		}
		branchID = b
	} else if role == "Cashier" {
		// Cashier must have initiated an offline cart previously, this endpoint needs cart_id? 
		// Ah wait, offline cart has ID. But AddCartItem in my service gets Cart from DB IF customerID != nil.
		// If Kasir, they must provide cartID directly?
		// My Service says: "if customerID != nil (Online) -> Find by Customer. Else -> return error".
		// Oh wait, for Offline Cart, we actually just use the regular "AddOrUpdateCartItem" repo without searching by customer.
		// Wait, my Service AddCartItem only handles Customer side for searching.
		response.Error(c, 400, "invalid: kasir harus nambahin langsung ke item dengan cart_id khusus. Endpoint ini khusus customer online.")
		return
	}

	err := h.service.AddCartItem(c.Request.Context(), customerID, branchID, req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, "item added to cart")
}

// InitOfflineCart
func (h *Handler) InitOfflineCart(c *gin.Context) {
	role := c.GetString("role")
	if role != "Cashier" {
		response.Error(c, 403, "forbidden")
		return
	}
	branchIDFloat, exists := c.Get("branch_id")
	if !exists {
		response.Error(c, 403, "forbidden: no branch assigned")
		return
	}
	branchID := int(branchIDFloat.(float64))

	var req InitOfflineCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	cart, err := h.service.InitOfflineCart(c.Request.Context(), branchID, req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, cart)
}

// Tambahan Khusus Kasir: AddItemToOfflineCart
func (h *Handler) AddItemToOfflineCart(c *gin.Context) {
	role := c.GetString("role")
	if role != "Cashier" {
		response.Error(c, 403, "forbidden")
		return
	}
	cartID := c.Param("id")

	var req AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}
	
	branchIDFloat, exists := c.Get("branch_id")
	if !exists {
		response.Error(c, 403, "forbidden: no branch assigned")
		return
	}
	branchID := int(branchIDFloat.(float64))

	err := h.service.AddItemToOfflineCart(c.Request.Context(), cartID, branchID, req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	
	response.Success(c, 200, "item added to offline cart")
}

// GetMyCart
func (h *Handler) GetMyCart(c *gin.Context) {
	role := c.GetString("role")
	if role != "Customer" {
		response.Error(c, 403, "forbidden")
		return
	}
	userID := c.GetString("user_id")

	cart, err := h.service.GetMyCart(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, 200, cart)
}

// Checkout
func (h *Handler) Checkout(c *gin.Context) {
	role := c.GetString("role")
	var customerID, cashierID *string
	userID := c.GetString("user_id")

	if role == "Customer" {
		customerID = &userID
	} else if role == "Cashier" {
		cashierID = &userID
	} else {
		response.Error(c, 403, "forbidden")
		return
	}

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "invalid payload: "+err.Error())
		return
	}

	trx, err := h.service.Checkout(c.Request.Context(), customerID, cashierID, req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	// Jika Online, Generate Midtrans URL
	var paymentURL string
	if req.OrderType == "Online_Pickup" || req.OrderType == "Online_Delivery" {
		if h.paymentSvc != nil {
			var custName, custEmail string
			if trx.CustomerName != nil {
				custName = *trx.CustomerName
			} else {
				custName = "Guest"
			}
			custEmail = "guest@example.com" // asumsikan default, atau fetch profile

			url, err := h.paymentSvc.CreateSnapURL(trx.ID, trx.TotalAmount, custName, custEmail)
			if err == nil {
				paymentURL = url
			} else {
				// Cetak error ke terminal agar mudah di-debug
				fmt.Println("[ERROR] Failed to create Midtrans Snap URL:", err)
			}
		}
	}

	response.Success(c, 200, gin.H{
		"transaction": trx,
		"payment_url": paymentURL,
	})
}
