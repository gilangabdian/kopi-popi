package analytics

import (
	"context"
	"errors"
)

type Service interface {
	GetSalesReport(ctx context.Context, role string, authBranchID *int, reqBranchID *int, startDate, endDate string) ([]SalesReportResponse, error)
	GetTopProducts(ctx context.Context, role string, authBranchID *int, reqBranchID *int, startDate, endDate string, limit int) ([]TopProductResponse, error)
	GetShiftReports(ctx context.Context, role string, authBranchID *int, reqBranchID *int, startDate, endDate string) ([]ShiftReportResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func authorizeBranchAccess(role string, authBranchID *int, reqBranchID *int) (*int, error) {
	if role == "Customer" || role == "Cashier" {
		return nil, errors.New("forbidden: you do not have permission to view reports")
	}

	if role == "Manager" {
		if authBranchID == nil {
			return nil, errors.New("forbidden: manager must have an assigned branch")
		}
		if reqBranchID != nil && *reqBranchID != *authBranchID {
			return nil, errors.New("forbidden: you can only view reports for your own branch")
		}
		// Force branch filter to manager's branch
		return authBranchID, nil
	}

	if role == "Admin" || role == "ADMIN" {
		// Admin can filter by any branch or no branch (all branches)
		return reqBranchID, nil
	}

	return nil, errors.New("unauthorized: invalid role")
}

func (s *service) GetSalesReport(ctx context.Context, role string, authBranchID *int, reqBranchID *int, startDate, endDate string) ([]SalesReportResponse, error) {
	finalBranchID, err := authorizeBranchAccess(role, authBranchID, reqBranchID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetSalesReport(ctx, finalBranchID, startDate, endDate)
}

func (s *service) GetTopProducts(ctx context.Context, role string, authBranchID *int, reqBranchID *int, startDate, endDate string, limit int) ([]TopProductResponse, error) {
	finalBranchID, err := authorizeBranchAccess(role, authBranchID, reqBranchID)
	if err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetTopProducts(ctx, finalBranchID, startDate, endDate, limit)
}

func (s *service) GetShiftReports(ctx context.Context, role string, authBranchID *int, reqBranchID *int, startDate, endDate string) ([]ShiftReportResponse, error) {
	finalBranchID, err := authorizeBranchAccess(role, authBranchID, reqBranchID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetShiftReports(ctx, finalBranchID, startDate, endDate)
}
