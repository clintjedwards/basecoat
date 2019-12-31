package service

import (
	"context"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/logger"
	"github.com/clintjedwards/toolkit/tkerrors"
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

// // SearchJobs takes in a search term and returns jobs that might match
// func (basecoat *API) SearchJobs(context context.Context, request *api.SearchJobsRequest) (*api.SearchJobsResponse, error) {

// 	account, present := getAccountFromContext(context)
// 	if !present {
// 		return &api.SearchJobsResponse{}, status.Error(codes.FailedPrecondition, "account required")
// 	}

// 	if request.Term == "" {
// 		return &api.SearchJobsResponse{}, status.Error(codes.FailedPrecondition, "search term required")
// 	}

// 	searchResults, err := basecoat.search.SearchJobs(account, request.Term)
// 	if err != nil {
// 		return &api.SearchJobsResponse{}, status.Error(codes.Internal, fmt.Sprintf("a search error occurred: %v", err))
// 	}

// 	return &api.SearchJobsResponse{Results: searchResults}, nil
// }

// ListContractors returns a list of all contractors on basecoat service
func (bc *API) ListContractors(context context.Context, request *api.ListContractorsRequest) (*api.ListContractorsResponse, error) {

	account, present := getAccountFromContext(context)
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
		return &api.CreateContractorResponse{}, status.Error(codes.FailedPrecondition, "company name required")
	}

	newContractor := api.Contractor{
		Company: request.Company,
		Email:   request.Email,
		Phone:   request.Phone,
		Contact: request.Contact,
		Jobs:    []*api.Job{},
	}

	for _, newJob := range request.Jobs {
		newContractor.Jobs = append(newContractor.Jobs, &api.Job{
			Name:     newJob.Name,
			Notes:    newJob.Notes,
			Formulas: newJob.Formulas,
			Address:  newJob.Address,
		})
	}

	err := bc.storage.CreateContractor(account, &newContractor)
	if err != nil {
		if err == tkerrors.ErrEntityExists {
			return &api.CreateContractorResponse{}, status.Error(codes.AlreadyExists, "could not save contractor; contractor already exists")
		}
		bc.log.Errorw("could not save contractor", "error", err)
		return &api.CreateContractorResponse{}, status.Error(codes.Internal, "could not save contractor")
	}

	//basecoat.search.UpdateContractorIndex(account, newContractor)

	logger.Log().Infow("contractor created", "contractor", newContractor)
	return &api.CreateContractorResponse{Contractor: &newContractor}, nil
}

// UpdateContractor updates an already existing Contractor
func (bc *API) UpdateContractor(context context.Context, request *api.UpdateContractorRequest) (*api.UpdateContractorResponse, error) {

	account, present := getAccountFromContext(context)
	if !present {
		return &api.UpdateContractorResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &api.UpdateContractorResponse{}, status.Error(codes.FailedPrecondition, "contractor id required")
	}

	updatedContractor := api.Contractor{
		Id:      request.Id,
		Company: request.Company,
		Email:   request.Email,
		Phone:   request.Phone,
		Contact: request.Contact,
		Jobs:    request.Jobs,
	}

	err := bc.storage.UpdateContractor(account, request.Id, &updatedContractor)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.UpdateContractorResponse{}, status.Error(codes.NotFound, "could not update contractor; contractor key not found")
		}
		bc.log.Errorw("could not update contractor", "error", err)
		return &api.UpdateContractorResponse{}, status.Error(codes.Internal, "could not update Contractor")
	}

	//basecoat.search.UpdateContractorIndex(account, updatedContractor)

	bc.log.Infow("contractor updated", "contractor", updatedContractor)
	return &api.UpdateContractorResponse{Contractor: &updatedContractor}, nil
}

// // DeleteJob removes a job
// func (basecoat *API) DeleteJob(context context.Context, request *api.DeleteJobRequest) (*api.DeleteJobResponse, error) {

// 	account, present := getAccountFromContext(context)
// 	if !present {
// 		return &api.DeleteJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
// 	}

// 	if request.Id == "" {
// 		return &api.DeleteJobResponse{}, status.Error(codes.FailedPrecondition, "job id required")
// 	}

// 	err := basecoat.storage.DeleteJob(account, request.Id)
// 	if err != nil {
// 		if err == tkerrors.ErrEntityNotFound {
// 			return &api.DeleteJobResponse{}, status.Error(codes.NotFound, "could not delete job; job key not found")
// 		}
// 		logger.Log().Errorw("could not delete job", "error", err, "job_id", request.Id)
// 		return &api.DeleteJobResponse{}, status.Error(codes.Internal, "could not delete job")
// 	}

// 	basecoat.search.DeleteJobIndex(account, request.Id)

// 	logger.Log().Infow("job deleted", "job_id", request.Id)
// 	return &api.DeleteJobResponse{}, nil
// }
