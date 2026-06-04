package branch

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	GetAllBranches(ctx context.Context, role string, includeInactive bool) ([]Branch, error)
	CreateBranch(ctx context.Context, req CreateBranchRequest) error
	UpdateBranch(ctx context.Context, id int, req UpdateBranchRequest) error
	DeleteBranch(ctx context.Context, id int) error
	UpdateOperatingHours(ctx context.Context, id int, req UpdateOperatingHoursRequest, role string) error
	ToggleAcceptingOrders(ctx context.Context, id int, req UpdateAcceptingOrdersRequest, role string, branchID *int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetAllBranches(ctx context.Context, role string, includeInactive bool) ([]Branch, error) {
	// Jika bukan Admin, paksa includeInactive menjadi false (Customer/Manager hanya lihat yang aktif)
	if role != "Admin" {
		includeInactive = false
	}
	
	branches, err := s.repo.FindAll(ctx, includeInactive)
	if err != nil {
		return nil, err
	}
	
	if branches == nil {
		branches = []Branch{}
	}
	
	return branches, nil
}

func (s *service) CreateBranch(ctx context.Context, req CreateBranchRequest) error {
	branch := &Branch{
		Name:              req.Name,
		Address:           req.Address,
		IsActive:          true,
		IsAcceptingOrders: true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	
	return s.repo.Create(ctx, branch)
}

func (s *service) UpdateBranch(ctx context.Context, id int, req UpdateBranchRequest) error {
	branch, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if branch == nil {
		return errors.New("branch not found")
	}

	if req.Name != nil {
		branch.Name = *req.Name
	}
	if req.Address != nil {
		branch.Address = *req.Address
	}
	if req.IsActive != nil {
		branch.IsActive = *req.IsActive
	}
	
	branch.UpdatedAt = time.Now()
	
	return s.repo.Update(ctx, branch)
}

func (s *service) DeleteBranch(ctx context.Context, id int) error {
	branch, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if branch == nil {
		return errors.New("branch not found")
	}

	// Soft delete
	branch.IsActive = false
	branch.UpdatedAt = time.Now()
	
	return s.repo.Update(ctx, branch)
}

func (s *service) UpdateOperatingHours(ctx context.Context, id int, req UpdateOperatingHoursRequest, role string) error {
	if role != "Admin" {
		return errors.New("forbidden: only admin can update operating hours")
	}

	branch, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if branch == nil {
		return errors.New("branch not found")
	}

	branch.OpeningTime = req.OpeningTime
	branch.ClosingTime = req.ClosingTime
	branch.UpdatedAt = time.Now()

	return s.repo.Update(ctx, branch)
}

func (s *service) ToggleAcceptingOrders(ctx context.Context, id int, req UpdateAcceptingOrdersRequest, role string, userBranchID *int) error {
	if role != "Admin" && role != "Manager" && role != "Cashier" {
		return errors.New("forbidden: unauthorized role")
	}

	// Manager and Cashier can only toggle their own branch
	if role != "Admin" {
		if userBranchID == nil || *userBranchID != id {
			return errors.New("forbidden: can only toggle accepting orders for your own branch")
		}
	}

	branch, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if branch == nil {
		return errors.New("branch not found")
	}

	branch.IsAcceptingOrders = *req.IsAcceptingOrders
	branch.UpdatedAt = time.Now()

	return s.repo.Update(ctx, branch)
}
