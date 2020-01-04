package service

import (
	"context"

	"github.com/clintjedwards/toolkit/tkerrors"
	"go.uber.org/zap"

	"github.com/clintjedwards/basecoat/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetContractor returns a single contractor by key
func (bc *API) GetContractor(ctx context.Context, request *api.GetContractorRequest) (*api.GetContractorResponse, error) {

	account, present := getAccountFromContext(ctx)
	if !present {
		return &api.GetContractorResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	contractor, err := bc.storage.GetContractor(account, request.Id)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.GetContractorResponse{}, status.Error(codes.NotFound, "contractor requested not found")
		}
		return &api.GetContractorResponse{}, status.Error(codes.Internal, "failed to retrieve contractor from database")
	}

	return &api.GetContractorResponse{Contractor: contractor}, nil
}

// ListContractors returns a list of all contractors on basecoat service
func (bc *API) ListContractors(ctx context.Context, request *api.ListContractorsRequest) (*api.ListContractorsResponse, error) {

	account, present := getAccountFromContext(ctx)
	if !present {
		return &api.ListContractorsResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	contractors, err := bc.storage.GetAllContractors(account)
	if err != nil {
		return &api.ListContractorsResponse{}, status.Error(codes.Internal, "failed to retrieve contractors from database")
	}

	return &api.ListContractorsResponse{Contractors: contractors}, nil
}

// CreateContractor registers a new contractor
func (bc *API) CreateContractor(ctx context.Context, request *api.CreateContractorRequest) (*api.CreateContractorResponse, error) {

	account, present := getAccountFromContext(ctx)
	if !present {
		return &api.CreateContractorResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Company == "" {
		return &api.CreateContractorResponse{}, status.Error(codes.FailedPrecondition, "company required")
	}

	newContractor := api.Contractor{
		Company: request.Company,
		Contact: request.Contact,
		Jobs:    request.Jobs,
	}

	contractorID, err := bc.storage.AddContractor(account, &newContractor)
	if err != nil {
		if err == tkerrors.ErrEntityExists {
			return &api.CreateContractorResponse{}, status.Error(codes.AlreadyExists, "could not save contractor; contractor already exists")
		}
		zap.S().Errorw("could not save contractor", "error", err)
		return &api.CreateContractorResponse{}, status.Error(codes.Internal, "could not save contractor")
	}

	newContractor.Id = contractorID

	// Update search index
	for _, job := range newContractor.Jobs {
		go bc.search.UpdateJobIndex(account, job)
	}

	zap.S().Infow("contractor created", "contractor", newContractor)
	return &api.CreateContractorResponse{Contractor: &newContractor}, nil
}

// UpdateContractor updates an already existing contractor
func (bc *API) UpdateContractor(ctx context.Context, request *api.UpdateContractorRequest) (*api.UpdateContractorResponse, error) {

	account, present := getAccountFromContext(ctx)
	if !present {
		return &api.UpdateContractorResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &api.UpdateContractorResponse{}, status.Error(codes.FailedPrecondition, "contractor id required")
	}

	updatedContractor := api.Contractor{
		Id:      request.Id,
		Company: request.Company,
		Contact: request.Contact,
		Jobs:    request.Jobs,
	}

	err := bc.storage.UpdateContractor(account, request.Id, &updatedContractor)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.UpdateContractorResponse{}, status.Error(codes.NotFound, "could not update contractor; contractor key not found")
		}
		zap.S().Errorw("could not update contractor", "error", err)
		return &api.UpdateContractorResponse{}, status.Error(codes.Internal, "could not update contractor")
	}

	zap.S().Infow("contractor updated", "contractor", updatedContractor)
	return &api.UpdateContractorResponse{Contractor: &updatedContractor}, nil
}

// DeleteContractor removes a contractor
func (bc *API) DeleteContractor(ctx context.Context, request *api.DeleteContractorRequest) (*api.DeleteContractorResponse, error) {

	account, present := getAccountFromContext(ctx)
	if !present {
		return &api.DeleteContractorResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &api.DeleteContractorResponse{}, status.Error(codes.FailedPrecondition, "contractor id required")
	}

	err := bc.storage.DeleteContractor(account, request.Id)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.DeleteContractorResponse{}, status.Error(codes.NotFound, "could not delete contractor; contractor key not found")
		}
		zap.S().Errorw("could not delete contractor", "error", err, "contractor_id", request.Id)
		return &api.DeleteContractorResponse{}, status.Error(codes.Internal, "could not delete contractor")
	}

	zap.S().Infow("contractor deleted", "contractor_id", request.Id)
	return &api.DeleteContractorResponse{}, nil
}
