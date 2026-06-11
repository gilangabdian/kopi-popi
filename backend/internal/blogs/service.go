package blogs

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type Service interface {
	CreateBlog(ctx context.Context, authorID string, req CreateBlogRequest) (*Blog, error)
	GetBlogs(ctx context.Context, limit int, offset int) ([]Blog, int64, error)
	GetBlogByIDOrSlug(ctx context.Context, idOrSlug string) (*Blog, error)
	UpdateBlog(ctx context.Context, id string, req UpdateBlogRequest) (*Blog, error)
	DeleteBlog(ctx context.Context, id string) error
	CalculateEstimatedReadTime(content string) int
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CalculateEstimatedReadTime menghitung estimasi waktu baca (asumsi 200 kata per menit)
func (s *service) CalculateEstimatedReadTime(content string) int {
	wordCount := len(strings.Fields(content))
	readTime := wordCount / 200
	if readTime == 0 && wordCount > 0 {
		return 1 // minimal 1 menit jika ada teks
	}
	return readTime
}

func (s *service) generateUniqueSlug(ctx context.Context, title string, currentID string) (string, error) {
	baseSlug := slug.Make(title)
	finalSlug := baseSlug
	counter := 1

	for {
		exists, err := s.repo.CheckSlugExists(ctx, finalSlug)
		if err != nil {
			return "", err
		}

		if !exists {
			break // slug is available
		}

		// If it exists, check if it belongs to the current blog (for updates)
		existingBlog, _ := s.repo.FindByIDOrSlug(ctx, finalSlug)
		if existingBlog != nil && existingBlog.ID == currentID {
			break // slug is available because it's already owned by this blog
		}

		// Slug is taken by another blog, increment counter
		finalSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return finalSlug, nil
}

func (s *service) CreateBlog(ctx context.Context, authorID string, req CreateBlogRequest) (*Blog, error) {
	newID := uuid.New().String()
	finalSlug, err := s.generateUniqueSlug(ctx, req.Title, newID)
	if err != nil {
		return nil, err
	}

	blog := &Blog{
		ID:                    newID,
		AuthorID:              authorID,
		Title:                 req.Title,
		Slug:                  finalSlug,
		Content:               req.Content,
		EstimatedReadTimeMins: s.CalculateEstimatedReadTime(req.Content),
	}

	err = s.repo.Create(ctx, blog)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *service) GetBlogs(ctx context.Context, limit int, offset int) ([]Blog, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *service) GetBlogByIDOrSlug(ctx context.Context, idOrSlug string) (*Blog, error) {
	return s.repo.FindByIDOrSlug(ctx, idOrSlug)
}

func (s *service) UpdateBlog(ctx context.Context, id string, req UpdateBlogRequest) (*Blog, error) {
	blog, err := s.repo.FindByIDOrSlug(ctx, id)
	if err != nil {
		return nil, err
	}

	finalSlug, err := s.generateUniqueSlug(ctx, req.Title, blog.ID)
	if err != nil {
		return nil, err
	}

	blog.Title = req.Title
	blog.Slug = finalSlug
	blog.Content = req.Content
	blog.EstimatedReadTimeMins = s.CalculateEstimatedReadTime(req.Content)

	err = s.repo.Update(ctx, blog)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *service) DeleteBlog(ctx context.Context, id string) error {
	_, err := s.repo.FindByIDOrSlug(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
