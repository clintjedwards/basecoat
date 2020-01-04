package service

import (
	"context"
	"fmt"
	"time"

	"github.com/clintjedwards/toolkit/tkerrors"
	"go.uber.org/zap"

	"github.com/clintjedwards/basecoat/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetJob returns a single job by key
func (bc *API) GetJob(ctx context.Context, request *api.GetJobRequest) (*api.GetJobResponse, error) {

	account, present := getAccountFromContext(ctx)
	if !present {
		return &api.GetJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	job, err := bc.storage.GetJob(account, request.Id)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.GetJobResponse{}, status.Error(codes.NotFound, "job requested not found")
		}
		return &api.GetJobResponse{}, status.Error(codes.Internal, "failed to retrieve job from database")
	}

	return &api.GetJobResponse{Job: job}, nil
}

// SearchJobs takes in a search term and returns jobs that might match
func (bc *API) SearchJobs(ctx context.Context, request *api.SearchJobsRequest) (*api.SearchJobsResponse, error) {

	account, present := getAccountFromContext(ctx)
	if !present {
		return &api.SearchJobsResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Term == "" {
		return &api.SearchJobsResponse{}, status.Error(codes.FailedPrecondition, "search term required")
	}

	searchResults, err := bc.search.SearchJobs(account, request.Term)
	if err != nil {
		return &api.SearchJobsResponse{}, status.Error(codes.Internal, fmt.Sprintf("a search error occurred: %v", err))
	}

	return &api.SearchJobsResponse{Results: searchResults}, nil
}

// ListJobs returns a list of all jobs on basecoat service
func (bc *API) ListJobs(ctx context.Context, request *api.ListJobsRequest) (*api.ListJobsResponse, error) {

	account, present := getAccountFromContext(ctx)
	if !present {
		return &api.ListJobsResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	jobs, err := bc.storage.GetAllJobs(account)
	if err != nil {
		return &api.ListJobsResponse{}, status.Error(codes.Internal, "failed to retrieve jobs from database")
	}

	return &api.ListJobsResponse{Jobs: jobs}, nil
}

// CreateJob registers a new job
func (bc *API) CreateJob(ctx context.Context, request *api.CreateJobRequest) (*api.CreateJobResponse, error) {

	newJob := api.Job{
		Name:         request.Name,
		Address:      request.Address,
		Notes:        request.Notes,
		Created:      time.Now().Unix(),
		Modified:     time.Now().Unix(),
		Formulas:     request.Formulas,
		ContractorId: request.ContractorId,
		Contact:      request.Contact,
	}

	account, present := getAccountFromContext(ctx)
	if !present {
		return &api.CreateJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if newJob.Name == "" {
		return &api.CreateJobResponse{}, status.Error(codes.FailedPrecondition, "name required")
	}

	jobID, err := bc.storage.AddJob(account, &newJob)
	if err != nil {
		if err == tkerrors.ErrEntityExists {
			return &api.CreateJobResponse{}, status.Error(codes.AlreadyExists, "could not save job; job already exists")
		}
		zap.S().Errorw("could not save job", "error", err)
		return &api.CreateJobResponse{}, status.Error(codes.Internal, "could not save job")
	}

	newJob.Id = jobID

	go bc.search.UpdateJobIndex(account, newJob.Id)

	zap.S().Infow("job created", "job", newJob)
	return &api.CreateJobResponse{Job: &newJob}, nil
}

// UpdateJob updates an already existing job
func (bc *API) UpdateJob(ctx context.Context, request *api.UpdateJobRequest) (*api.UpdateJobResponse, error) {

	account, present := getAccountFromContext(ctx)
	if !present {
		return &api.UpdateJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &api.UpdateJobResponse{}, status.Error(codes.FailedPrecondition, "job id required")
	}

	currentJob, err := bc.storage.GetJob(account, request.Id)
	if err != nil {
		return &api.UpdateJobResponse{}, status.Error(codes.FailedPrecondition, "could not get current job")
	}

	updatedJob := api.Job{
		Id:           request.Id,
		Name:         request.Name,
		Address:      request.Address,
		Notes:        request.Notes,
		Modified:     time.Now().Unix(),
		Created:      currentJob.Created,
		Formulas:     request.Formulas,
		ContractorId: request.ContractorId,
		Contact:      request.Contact,
	}

	err = bc.storage.UpdateJob(account, request.Id, &updatedJob)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.UpdateJobResponse{}, status.Error(codes.NotFound, "could not update job; job key not found")
		}
		zap.S().Errorw("could not update job", "error", err)
		return &api.UpdateJobResponse{}, status.Error(codes.Internal, "could not update job")
	}

	go bc.search.UpdateJobIndex(account, updatedJob.Id)

	zap.S().Infow("job updated", "job", updatedJob)
	return &api.UpdateJobResponse{Job: &updatedJob}, nil
}

// DeleteJob removes a job
func (bc *API) DeleteJob(ctx context.Context, request *api.DeleteJobRequest) (*api.DeleteJobResponse, error) {

	account, present := getAccountFromContext(ctx)
	if !present {
		return &api.DeleteJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &api.DeleteJobResponse{}, status.Error(codes.FailedPrecondition, "job id required")
	}

	err := bc.storage.DeleteJob(account, request.Id)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.DeleteJobResponse{}, status.Error(codes.NotFound, "could not delete job; job key not found")
		}
		zap.S().Errorw("could not delete job", "error", err, "job_id", request.Id)
		return &api.DeleteJobResponse{}, status.Error(codes.Internal, "could not delete job")
	}

	go bc.search.DeleteJobIndex(account, request.Id)

	zap.S().Infow("job deleted", "job_id", request.Id)
	return &api.DeleteJobResponse{}, nil
}
