package service

import (
	"context"
	"fmt"
	"time"

	"github.com/clintjedwards/toolkit/tkerrors"

	"github.com/clintjedwards/basecoat/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetJob returns a single job by key
func (bc *API) GetJob(context context.Context, request *api.GetJobRequest) (*api.GetJobResponse, error) {

	account, present := getAccountFromContext(context)
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
func (bc *API) SearchJobs(context context.Context, request *api.SearchJobsRequest) (*api.SearchJobsResponse, error) {

	account, present := getAccountFromContext(context)
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
func (bc *API) ListJobs(context context.Context, request *api.ListJobsRequest) (*api.ListJobsResponse, error) {

	account, present := getAccountFromContext(context)
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
func (bc *API) CreateJob(context context.Context, request *api.CreateJobRequest) (*api.CreateJobResponse, error) {

	newJob := api.Job{
		Name:     request.Name,
		Street:   request.Street,
		Street2:  request.Street2,
		City:     request.City,
		State:    request.State,
		Zipcode:  request.Zipcode,
		Notes:    request.Notes,
		Created:  time.Now().Unix(),
		Modified: time.Now().Unix(),
		Formulas: request.Formulas,
		Contact:  request.Contact,
	}

	account, present := getAccountFromContext(context)
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
		bc.log.Errorw("could not save job", "error", err)
		return &api.CreateJobResponse{}, status.Error(codes.Internal, "could not save job")
	}

	newJob.Id = jobID

	bc.search.UpdateJobIndex(account, newJob)

	bc.log.Infow("job created", "job", newJob)
	return &api.CreateJobResponse{Id: newJob.Id}, nil
}

// UpdateJob updates an already existing job
func (bc *API) UpdateJob(context context.Context, request *api.UpdateJobRequest) (*api.UpdateJobResponse, error) {

	account, present := getAccountFromContext(context)
	if !present {
		return &api.UpdateJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &api.UpdateJobResponse{}, status.Error(codes.FailedPrecondition, "job id required")
	}

	currentJob, _ := bc.storage.GetJob(account, request.Id)

	updatedJob := api.Job{
		Id:       request.Id,
		Name:     request.Name,
		Street:   request.Street,
		Street2:  request.Street2,
		City:     request.City,
		State:    request.State,
		Zipcode:  request.Zipcode,
		Notes:    request.Notes,
		Modified: time.Now().Unix(),
		Created:  currentJob.Created,
		Formulas: request.Formulas,
		Contact:  request.Contact,
	}

	err := bc.storage.UpdateJob(account, request.Id, &updatedJob)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.UpdateJobResponse{}, status.Error(codes.NotFound, "could not update job; job key not found")
		}
		bc.log.Errorw("could not update job", "error", err)
		return &api.UpdateJobResponse{}, status.Error(codes.Internal, "could not update job")
	}

	bc.search.UpdateJobIndex(account, updatedJob)

	bc.log.Infow("job updated", "job", updatedJob)
	return &api.UpdateJobResponse{}, nil
}

// DeleteJob removes a job
func (bc *API) DeleteJob(context context.Context, request *api.DeleteJobRequest) (*api.DeleteJobResponse, error) {

	account, present := getAccountFromContext(context)
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
		bc.log.Errorw("could not delete job", "error", err, "job_id", request.Id)
		return &api.DeleteJobResponse{}, status.Error(codes.Internal, "could not delete job")
	}

	bc.search.DeleteJobIndex(account, request.Id)

	bc.log.Infow("job deleted", "job_id", request.Id)
	return &api.DeleteJobResponse{}, nil
}
