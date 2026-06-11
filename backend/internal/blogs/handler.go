package blogs

import (
	"strconv"

	"github.com/gilangages/kopi-popi/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) GetBlogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	blogs, total, err := h.service.GetBlogs(c.Request.Context(), limit, offset)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "Success",
		"data":    blogs,
		"meta": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

func (h *Handler) GetBlogByID(c *gin.Context) {
	id := c.Param("id")
	blog, err := h.service.GetBlogByIDOrSlug(c.Request.Context(), id)
	if err != nil {
		response.Error(c, 404, err.Error())
		return
	}

	response.Success(c, 200, blog)
}

func (h *Handler) CreateBlog(c *gin.Context) {
	var req CreateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	authorID := c.GetString("user_id") // From auth middleware

	blog, err := h.service.CreateBlog(c.Request.Context(), authorID, req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, 201, blog)
}

func (h *Handler) UpdateBlog(c *gin.Context) {
	id := c.Param("id")
	var req UpdateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	blog, err := h.service.UpdateBlog(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, 200, blog)
}

func (h *Handler) DeleteBlog(c *gin.Context) {
	id := c.Param("id")
	err := h.service.DeleteBlog(c.Request.Context(), id)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, 200, nil)
}
