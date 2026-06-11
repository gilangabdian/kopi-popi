package catalog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gosimple/slug"
)

type Service interface {
	// Category
	GetAllCategories(ctx context.Context) ([]Category, error)
	GetCategoryByIDOrSlug(ctx context.Context, idOrSlug string) (*Category, error)
	CreateCategory(ctx context.Context, req CategoryRequest) error
	UpdateCategory(ctx context.Context, id int, req CategoryRequest) error
	DeleteCategory(ctx context.Context, id int) error

	// Material
	GetAllMaterials(ctx context.Context) ([]Material, error)
	CreateMaterial(ctx context.Context, req MaterialRequest) error
	UpdateMaterial(ctx context.Context, id int, req MaterialRequest) error
	DeleteMaterial(ctx context.Context, id int) error

	// Product
	GetAllProducts(ctx context.Context, categoryID *int, search string) ([]Product, error)
	GetProductDetail(ctx context.Context, idOrSlug string, role string, includeRecipe bool) (*Product, error)
	GetProductsBOM(ctx context.Context, productIDs []int) (map[int][]ProductBOM, error)
	CreateProduct(ctx context.Context, req ProductRequest) error
	UpdateProduct(ctx context.Context, id int, req ProductRequest) error
	DeleteProduct(ctx context.Context, id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

// -- Category --
func (s *service) generateUniqueCategorySlug(ctx context.Context, name string, currentID int) (string, error) {
	baseSlug := slug.Make(name)
	finalSlug := baseSlug
	counter := 1

	for {
		exists, err := s.repo.CheckCategorySlugExists(ctx, finalSlug)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
		existingCat, _ := s.repo.FindCategoryByIDOrSlug(ctx, finalSlug)
		if existingCat != nil && existingCat.ID == currentID {
			break
		}
		finalSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
	return finalSlug, nil
}

func (s *service) GetAllCategories(ctx context.Context) ([]Category, error) {
	categories, err := s.repo.FindAllCategories(ctx)
	if err != nil {
		return nil, err
	}
	if categories == nil {
		categories = []Category{}
	}
	return categories, nil
}

func (s *service) GetCategoryByIDOrSlug(ctx context.Context, idOrSlug string) (*Category, error) {
	return s.repo.FindCategoryByIDOrSlug(ctx, idOrSlug)
}

func (s *service) CreateCategory(ctx context.Context, req CategoryRequest) error {
	finalSlug, err := s.generateUniqueCategorySlug(ctx, req.Name, 0)
	if err != nil {
		return err
	}
	cat := &Category{Name: req.Name, Slug: finalSlug}
	return s.repo.CreateCategory(ctx, cat)
}

func (s *service) UpdateCategory(ctx context.Context, id int, req CategoryRequest) error {
	cat, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return errors.New("category not found")
	}

	finalSlug, err := s.generateUniqueCategorySlug(ctx, req.Name, cat.ID)
	if err != nil {
		return err
	}

	cat.Name = req.Name
	cat.Slug = finalSlug
	return s.repo.UpdateCategory(ctx, cat)
}

func (s *service) DeleteCategory(ctx context.Context, id int) error {
	cat, err := s.repo.FindCategoryByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return errors.New("category not found")
	}
	return s.repo.DeleteCategory(ctx, id)
}

// -- Material --
func (s *service) GetAllMaterials(ctx context.Context) ([]Material, error) {
	mats, err := s.repo.FindAllMaterials(ctx)
	if err != nil {
		return nil, err
	}
	if mats == nil {
		mats = []Material{}
	}
	return mats, nil
}

func (s *service) CreateMaterial(ctx context.Context, req MaterialRequest) error {
	mat := &Material{
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Unit:       req.Unit,
		CreatedAt:  time.Now(),
	}
	return s.repo.CreateMaterial(ctx, mat)
}

func (s *service) UpdateMaterial(ctx context.Context, id int, req MaterialRequest) error {
	mat, err := s.repo.FindMaterialByID(ctx, id)
	if err != nil {
		return err
	}
	if mat == nil {
		return errors.New("material not found")
	}
	mat.CategoryID = req.CategoryID
	mat.Name = req.Name
	mat.Unit = req.Unit
	return s.repo.UpdateMaterial(ctx, mat)
}

func (s *service) DeleteMaterial(ctx context.Context, id int) error {
	mat, err := s.repo.FindMaterialByID(ctx, id)
	if err != nil {
		return err
	}
	if mat == nil {
		return errors.New("material not found")
	}
	return s.repo.DeleteMaterial(ctx, id)
}

// -- Product --
func (s *service) generateUniqueProductSlug(ctx context.Context, name string, currentID int) (string, error) {
	baseSlug := slug.Make(name)
	finalSlug := baseSlug
	counter := 1

	for {
		exists, err := s.repo.CheckProductSlugExists(ctx, finalSlug)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
		existingProd, _ := s.repo.FindProductByIDOrSlug(ctx, finalSlug)
		if existingProd != nil && existingProd.ID == currentID {
			break
		}
		finalSlug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
	return finalSlug, nil
}

func (s *service) GetAllProducts(ctx context.Context, categoryID *int, search string) ([]Product, error) {
	products, err := s.repo.FindAllProducts(ctx, categoryID, search)
	if err != nil {
		return nil, err
	}
	if products == nil {
		products = []Product{}
	}
	// BOM dipastikan nil untuk list
	for i := range products {
		products[i].Recipe = nil
	}
	return products, nil
}

func (s *service) GetProductDetail(ctx context.Context, idOrSlug string, role string, includeRecipe bool) (*Product, error) {
	product, err := s.repo.FindProductByIDOrSlug(ctx, idOrSlug)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Filter hak akses resep (Hanya Admin / Manager yang bisa includeRecipe)
	if !includeRecipe || (role != "Admin" && role != "Manager") {
		product.Recipe = nil
	}

	return product, nil
}

func (s *service) GetProductsBOM(ctx context.Context, productIDs []int) (map[int][]ProductBOM, error) {
	result := make(map[int][]ProductBOM)
	if len(productIDs) == 0 {
		return result, nil
	}

	for _, pid := range productIDs {
		product, err := s.repo.FindProductByID(ctx, pid)
		if err != nil {
			return nil, err
		}
		if product != nil && product.Recipe != nil {
			result[pid] = product.Recipe
		} else {
			result[pid] = []ProductBOM{}
		}
	}

	return result, nil
}

func (s *service) CreateProduct(ctx context.Context, req ProductRequest) error {
	finalSlug, err := s.generateUniqueProductSlug(ctx, req.Name, 0)
	if err != nil {
		return err
	}

	product := &Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        finalSlug,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	var boms []ProductBOM
	for _, rb := range req.Recipe {
		boms = append(boms, ProductBOM{
			MaterialID:     rb.MaterialID,
			QuantityNeeded: rb.QuantityNeeded,
		})
	}

	return s.repo.CreateProductWithBOM(ctx, product, boms)
}

func (s *service) UpdateProduct(ctx context.Context, id int, req ProductRequest) error {
	product, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	finalSlug, err := s.generateUniqueProductSlug(ctx, req.Name, product.ID)
	if err != nil {
		return err
	}

	product.CategoryID = req.CategoryID
	product.Name = req.Name
	product.Slug = finalSlug
	product.Description = req.Description
	product.Price = req.Price
	product.ImageURL = req.ImageURL
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	product.UpdatedAt = time.Now()

	var boms []ProductBOM
	for _, rb := range req.Recipe {
		boms = append(boms, ProductBOM{
			MaterialID:     rb.MaterialID,
			QuantityNeeded: rb.QuantityNeeded,
		})
	}

	return s.repo.UpdateProductWithBOM(ctx, product, boms)
}

func (s *service) DeleteProduct(ctx context.Context, id int) error {
	product, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	return s.repo.SoftDeleteProduct(ctx, id)
}
