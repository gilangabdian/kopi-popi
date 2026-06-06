package main

import (
	config "github.com/gilangages/kopi-popi/configs"
	"github.com/gilangages/kopi-popi/internal/auth"
	"github.com/gilangages/kopi-popi/internal/branch"
	"github.com/gilangages/kopi-popi/internal/catalog"
	"github.com/gilangages/kopi-popi/internal/inventory"
	"github.com/gilangages/kopi-popi/internal/media"
	"github.com/gilangages/kopi-popi/internal/notification"
	"github.com/gilangages/kopi-popi/internal/payment"
	"github.com/gilangages/kopi-popi/internal/sales"
	"github.com/gilangages/kopi-popi/internal/user"
	"github.com/gilangages/kopi-popi/pkg/middleware"
	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inisialisasi Koneksi ke Database
	// Memastikan database menyala sebelum route dijalankan
	db := config.ConnectDB()
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	// 2. Setup Framework Gin (Router)
	r := gin.Default()

	// 3. Setup Global Middleware CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 4. Register Health Check Endpoint (Public)
	r.GET("/", func(c *gin.Context) {
		response.Success(c, 200, gin.H{
			"message": "Welcome to Kopi-Popi API!",
			"version": "1.0.0",
		})
	})

	// 5. Inisialisasi Domain Users
	usersRepo := user.NewRepository(db)
	usersService := user.NewService(usersRepo)
	usersHandler := user.NewHandler(usersService)

	// 5a. Inisialisasi Domain Notifications (Harus awal karena banyak yang butuh)
	notifRepo := notification.NewRepository(db)
	notifService := notification.NewService(notifRepo)
	notifHandler := notification.NewHandler(notifService)

	// 5b. Inisialisasi Domain Auth
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, notifService)
	authHandler := auth.NewHandler(authService)

	// 5c. Inisialisasi Domain Branch
	branchesRepo := branch.NewRepository(db)
	branchesService := branch.NewService(branchesRepo)
	branchesHandler := branch.NewHandler(branchesService)

	// 5d. Inisialisasi Domain Catalog
	catalogRepo := catalog.NewRepository(db)
	catalogService := catalog.NewService(catalogRepo)
	catalogHandler := catalog.NewHandler(catalogService)

	// 5e. Inisialisasi Domain Media
	mediaService := media.NewService()
	mediaHandler := media.NewHandler(mediaService)

	// 5f. Inisialisasi Domain Inventory
	inventoryRepo := inventory.NewRepository(db)
	inventoryService := inventory.NewService(inventoryRepo, notifService)
	inventoryHandler := inventory.NewHandler(inventoryService)

	// 5g. Inisialisasi Domain Payment
	paymentRepo := payment.NewRepository(db)
	paymentService := payment.NewService(paymentRepo, notifService)
	paymentHandler := payment.NewHandler(paymentService)

	// 5h. Inisialisasi Domain Sales
	salesRepo := sales.NewRepository(db)
	salesService := sales.NewService(salesRepo, branchesService, catalogService, inventoryService, notifService, usersService)
	salesHandler := sales.NewHandler(salesService, paymentService)

	// 6. Daftarkan router per-domain (Public)
	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/verify-email", authHandler.VerifyEmail)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/forgot-password", authHandler.ForgotPassword)
		authRoutes.POST("/reset-password", authHandler.ResetPassword)
		authRoutes.DELETE("/logout", authHandler.Logout)
	}

	// Webhook Midtrans (Public)
	r.POST("/payment/midtrans/webhook", paymentHandler.MidtransWebhook)

	// 7. Expose Static Folder (Supaya gambar bisa diakses publik)
	r.Static("/uploads", "./uploads")

	// 8. Daftarkan router per-domain (Protected by JWT)
	protectedRoutes := r.Group("/")
	protectedRoutes.Use(middleware.RequireAuth())
	{
		// Upload File
		protectedRoutes.POST("/uploads", mediaHandler.UploadFile)

		// Users Profile
		protectedRoutes.GET("/users/me", usersHandler.GetMyProfile)
		protectedRoutes.PATCH("/users/me", usersHandler.UpdateMyProfile)
		protectedRoutes.DELETE("/users/me/profile-picture", usersHandler.DeleteProfilePicture)
		protectedRoutes.PATCH("/users/me/password", usersHandler.UpdateMyPassword)
		protectedRoutes.POST("/users/me/request-email-otp", usersHandler.RequestEmailOTP)
		protectedRoutes.PUT("/users/me/email", usersHandler.VerifyEmailOTP)

		// Users Management(ADMIN)
		protectedRoutes.GET("/users", usersHandler.GetEmployees)
		protectedRoutes.POST("/users/managers", usersHandler.CreateManager)
		protectedRoutes.POST("/users/cashiers", usersHandler.CreateCashier)
		protectedRoutes.PATCH("/users/:id/disable", usersHandler.DisableEmployee)

		// Branches Management (ADMIN & MANAGER/CASHIER)
		protectedRoutes.POST("/branches", branchesHandler.CreateBranch)
		protectedRoutes.PUT("/branches/:id", branchesHandler.UpdateBranch)
		protectedRoutes.DELETE("/branches/:id", branchesHandler.DeleteBranch)
		protectedRoutes.PATCH("/branches/:id/operating-hours", branchesHandler.UpdateOperatingHours)
		protectedRoutes.PATCH("/branches/:id/accepting-orders", branchesHandler.ToggleAcceptingOrders)

		// Catalogues Management (ADMIN & MANAGER)
		protectedRoutes.POST("/categories", catalogHandler.CreateCategory)
		protectedRoutes.PUT("/categories/:id", catalogHandler.UpdateCategory)
		protectedRoutes.DELETE("/categories/:id", catalogHandler.DeleteCategory)
		
		protectedRoutes.GET("/materials", catalogHandler.GetAllMaterials)
		protectedRoutes.POST("/materials", catalogHandler.CreateMaterial)
		protectedRoutes.PUT("/materials/:id", catalogHandler.UpdateMaterial)
		protectedRoutes.DELETE("/materials/:id", catalogHandler.DeleteMaterial)

		protectedRoutes.POST("/products", catalogHandler.CreateProduct)
		protectedRoutes.PUT("/products/:id", catalogHandler.UpdateProduct)
		protectedRoutes.DELETE("/products/:id", catalogHandler.DeleteProduct)

		// Inventory Management (ADMIN & MANAGER)
		protectedRoutes.GET("/inventories/branches/:branch_id", inventoryHandler.GetBranchStock)
		protectedRoutes.GET("/inventories/movements", inventoryHandler.GetInventoryMovements)

		// Stok Gudang Pusat (Admin Only)
		protectedRoutes.POST("/inventories/central-warehouse/incoming", inventoryHandler.ReceiveIncomingStock)
		
		// Alokasi Stok Cabang (Admin Only)
		protectedRoutes.POST("/inventories/branches/:branch_id/allocate", inventoryHandler.AllocateStock)

		// Restock Requests
		protectedRoutes.GET("/inventories/requests", inventoryHandler.GetRestockRequests)
		protectedRoutes.POST("/inventories/restocks", inventoryHandler.CreateRestockRequest)
		protectedRoutes.PATCH("/inventories/restocks/:id/status", inventoryHandler.UpdateRestockStatus)

		// Sales & POS Management
		protectedRoutes.POST("/shifts/open", salesHandler.OpenShift)
		protectedRoutes.POST("/shifts/close", salesHandler.CloseShift)
		protectedRoutes.GET("/shifts/me", salesHandler.GetMyOpenShift)

		protectedRoutes.POST("/carts/offline", salesHandler.InitOfflineCart)
		protectedRoutes.POST("/carts/items", salesHandler.AddCartItem)
		protectedRoutes.POST("/carts/:id/items", salesHandler.AddItemToOfflineCart)
		protectedRoutes.GET("/carts/me", salesHandler.GetMyCart)

		protectedRoutes.POST("/checkout", salesHandler.Checkout)

		// Notifications Management
		protectedRoutes.GET("/notifications", notifHandler.GetMyNotifications)
		protectedRoutes.PUT("/notifications/:id/read", notifHandler.MarkAsRead)
	}

	// 8. Daftarkan router dengan Optional Auth (untuk public route yang behavior-nya berubah jika login)
	optionalAuthRoutes := r.Group("/")
	optionalAuthRoutes.Use(middleware.OptionalAuth())
	{
		// Branches Public (Bisa dilihat tanpa login, tapi Admin bisa request include_inactive)
		optionalAuthRoutes.GET("/branches", branchesHandler.GetAllBranches)
		
		// Catalogues Public (Tapi detail products punya resep khusus Admin)
		optionalAuthRoutes.GET("/categories", catalogHandler.GetAllCategories)
		optionalAuthRoutes.GET("/products", catalogHandler.GetAllProducts)
		optionalAuthRoutes.GET("/products/:id", catalogHandler.GetProductDetail)
	}

	// 9. Jalankan Server di port 8080
	r.Run(":8080")
}
