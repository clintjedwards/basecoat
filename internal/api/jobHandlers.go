package api

import (
	"context"
	"time"

	"github.com/clintjedwards/basecoat/internal/models"
	"github.com/clintjedwards/basecoat/internal/storage"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetJob returns a single job by key
func (api *API) GetJob(ctx context.Context, request *proto.GetJobRequest) (*proto.GetJobResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.GetJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.GetJobResponse{}, status.Error(codes.FailedPrecondition, "id required")
	}

	jobRaw, err := api.db.GetJob(api.db, account, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.GetJobResponse{}, status.Error(codes.NotFound, "job requested not found")
		}
		return &proto.GetJobResponse{}, status.Error(codes.Internal, "failed to retrieve job from database")
	}

	job := models.Job{}
	job.FromStorage(&jobRaw)

	return &proto.GetJobResponse{Job: job.ToProto()}, nil
}

// ListJobs returns a list of all jobs's metadata.
func (api *API) ListJobs(ctx context.Context, request *proto.ListJobsRequest) (*proto.ListJobsResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.ListJobsResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	searchResults := []string{}

	if request.Filter != "" {
		results, err := api.search.SearchFormulas(account, request.Filter)
		if err != nil {
			log.Error().Err(err).Msg("a search error occurred")
		}
		searchResults = results
		log.Debug().Strs("returned_results", searchResults).Str("query", request.Filter).Msg("filtered jobs on user's request")
	}

	jobsRaw, err := api.db.ListJobs(api.db, account, int(request.Offset), int(request.Limit))
	if err != nil {
		return &proto.ListJobsResponse{}, status.Error(codes.Internal, "failed to retrieve jobs from database")
	}

	protoJobs := []*proto.Job{}
	for _, jobRaw := range jobsRaw {
		// If the user actually enters a filter we want to use that filter.
		if request.Filter != "" {
			// If the job doesn't exist in the filter then we should skip it.
			if !slices.Contains(searchResults, jobRaw.ID) {
				continue
			}
		}

		var job models.Job
		job.FromStorage(&jobRaw)
		protoJobs = append(protoJobs, job.ToProto())
	}

	return &proto.ListJobsResponse{Jobs: protoJobs}, nil
}

// CreateJob registers a new job
func (api *API) CreateJob(ctx context.Context, request *proto.CreateJobRequest) (*proto.CreateJobResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.CreateJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Name == "" {
		return &proto.CreateJobResponse{}, status.Error(codes.FailedPrecondition, "job name required")
	}

	job := models.NewJob(account, request.ContractorId, request.Name)
	job.Contact = request.ContactId

	if request.Address != nil {
		address := models.Address{}
		address.FromProto(request.Address)
		job.Address = address
	}

	err := api.db.InsertJob(api.db, job.ToStorage())
	if err != nil {
		if err == storage.ErrEntityExists {
			return &proto.CreateJobResponse{}, status.Error(codes.AlreadyExists, "could not save job; job already exists")
		}
		log.Error().Err(err).Msg("could not save job")
		return &proto.CreateJobResponse{}, status.Error(codes.Internal, "could not save job")
	}

	log.Info().Str("id", job.ID).Str("name", job.Name).Msg("job created")

	return &proto.CreateJobResponse{
		Job: job.ToProto(),
	}, nil
}

// AssociateFormulaWithJob records the formula given as associated with the job given.
func (api *API) AssociateFormulaWithJob(ctx context.Context, request *proto.AssociateFormulaWithJobRequest) (*proto.AssociateFormulaWithJobResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.AssociateFormulaWithJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Job == "" {
		return &proto.AssociateFormulaWithJobResponse{}, status.Error(codes.FailedPrecondition, "job id required")
	}

	if request.Formula == "" {
		return &proto.AssociateFormulaWithJobResponse{}, status.Error(codes.FailedPrecondition, "formula id required")
	}

	err := api.db.AssociateFormulaWithJob(api.db, &storage.FormulaJob{
		Account: account,
		Job:     request.Job,
		Formula: request.Formula,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not associate formula with job")
		return &proto.AssociateFormulaWithJobResponse{}, status.Error(codes.Internal, "could not associate formula with job")
	}

	log.Debug().Str("formula", request.Formula).Str("job", request.Job).Msg("associated formula with job")
	return &proto.AssociateFormulaWithJobResponse{}, nil
}

// DisassociateFormulaFromJob records the formula given as associated with the job given.
func (api *API) DisassociateFormulaFromJob(ctx context.Context, request *proto.DisassociateFormulaFromJobRequest) (*proto.DisassociateFormulaFromJobResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.DisassociateFormulaFromJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Job == "" {
		return &proto.DisassociateFormulaFromJobResponse{}, status.Error(codes.FailedPrecondition, "job id required")
	}

	if request.Formula == "" {
		return &proto.DisassociateFormulaFromJobResponse{}, status.Error(codes.FailedPrecondition, "formula id required")
	}

	err := api.db.DeleteJobFormula(api.db, account, request.Job, request.Formula)
	if err != nil {
		log.Error().Err(err).Msg("could not disassociate formula from job")
		return &proto.DisassociateFormulaFromJobResponse{}, status.Error(codes.Internal, "could not disassociate formula from job")
	}

	log.Debug().Str("formula", request.Formula).Str("job", request.Job).Msg("disassociated formula from job")
	return &proto.DisassociateFormulaFromJobResponse{}, nil
}

// UpdateJob updates an already existing job
func (api *API) UpdateJob(ctx context.Context, request *proto.UpdateJobRequest) (*proto.UpdateJobResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.UpdateJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.UpdateJobResponse{}, status.Error(codes.FailedPrecondition, "job id required")
	}

	var jsonAddress *string
	if request.Address != nil {
		address := models.Address{}
		address.FromProto(request.Address)
		jsonAddress = ptr(address.ToJSON())
	}

	err := api.db.UpdateJob(api.db, account, request.Id, storage.UpdatableJobFields{
		Name:     request.Name,
		Address:  jsonAddress,
		Notes:    request.Notes,
		Contact:  request.ContactId,
		Modified: ptr(time.Now().UnixMilli()),
	})
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.UpdateJobResponse{}, status.Error(codes.NotFound, "job requested not found")
		}
		log.Error().Err(err).Msg("could not save job")
		return &proto.UpdateJobResponse{}, status.Error(codes.Internal, "could not save job")
	}

	log.Debug().Str("id", request.Id).Msg("job updated")
	return &proto.UpdateJobResponse{}, nil
}

func (api *API) DeleteJob(ctx context.Context, request *proto.DeleteJobRequest) (*proto.DeleteJobResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.DeleteJobResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.DeleteJobResponse{}, status.Error(codes.FailedPrecondition, "job id required")
	}

	err := api.db.DeleteJob(api.db, account, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.DeleteJobResponse{}, status.Error(codes.NotFound, "could not delete job; job key not found")
		}

		log.Error().Err(err).Msg("could not delete job")
		return &proto.DeleteJobResponse{}, status.Error(codes.Internal, "could not delete job")
	}

	return &proto.DeleteJobResponse{}, nil
}
