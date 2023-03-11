package api

import (
	"context"
	"time"

	"github.com/clintjedwards/basecoat/internal/models"
	"github.com/clintjedwards/basecoat/internal/storage"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetContractor returns a single contractor by key
func (api *API) GetContractor(ctx context.Context, request *proto.GetContractorRequest) (*proto.GetContractorResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.GetContractorResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.GetContractorResponse{}, status.Error(codes.FailedPrecondition, "id required")
	}

	contractorRaw, err := api.db.GetContractor(api.db, account, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.GetContractorResponse{}, status.Error(codes.NotFound, "contractor requested not found")
		}
		return &proto.GetContractorResponse{}, status.Error(codes.Internal, "failed to retrieve contractor from database")
	}

	contractor := models.Contractor{}
	contractor.FromStorage(&contractorRaw)

	return &proto.GetContractorResponse{Contractor: contractor.ToProto()}, nil
}

// ListContractors returns a list of all contractors's metadata.
func (api *API) ListContractors(ctx context.Context, _ *proto.ListContractorsRequest) (*proto.ListContractorsResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.ListContractorsResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	contractorsRaw, err := api.db.ListContractors(api.db, account, 0, 0)
	if err != nil {
		return &proto.ListContractorsResponse{}, status.Error(codes.Internal, "failed to retrieve contractors from database")
	}

	protoContractors := []*proto.Contractor{}
	for _, contractorRaw := range contractorsRaw {
		var contractor models.Contractor
		contractor.FromStorage(&contractorRaw)
		protoContractors = append(protoContractors, contractor.ToProto())
	}

	return &proto.ListContractorsResponse{Contractors: protoContractors}, nil
}

// CreateContractor registers a new contractor
func (api *API) CreateContractor(ctx context.Context, request *proto.CreateContractorRequest) (*proto.CreateContractorResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.CreateContractorResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Company == "" {
		return &proto.CreateContractorResponse{}, status.Error(codes.FailedPrecondition, "contractor company required")
	}

	contractor := models.NewContractor(account, request.Company)
	contractor.Contact = request.Contact

	err := api.db.InsertContractor(api.db, contractor.ToStorage())
	if err != nil {
		if err == storage.ErrEntityExists {
			return &proto.CreateContractorResponse{}, status.Error(codes.AlreadyExists, "could not save contractor; contractor already exists")
		}
		log.Error().Err(err).Msg("could not save contractor")
		return &proto.CreateContractorResponse{}, status.Error(codes.Internal, "could not save contractor")
	}

	log.Info().Str("id", contractor.ID).Str("company", contractor.Company).Msg("contractor created")

	return &proto.CreateContractorResponse{
		Contractor: contractor.ToProto(),
	}, nil
}

// UpdateContractor updates an already existing contractor
func (api *API) UpdateContractor(ctx context.Context, request *proto.UpdateContractorRequest) (*proto.UpdateContractorResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.UpdateContractorResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.UpdateContractorResponse{}, status.Error(codes.FailedPrecondition, "contractor id required")
	}

	err := api.db.UpdateContractor(api.db, account, request.Id, storage.UpdatableContractorFields{
		Company:  request.Company,
		Contact:  request.Contact,
		Modified: ptr(time.Now().UnixMilli()),
	})
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.UpdateContractorResponse{}, status.Error(codes.NotFound, "contractor requested not found")
		}
		log.Error().Err(err).Msg("could not save contractor")
		return &proto.UpdateContractorResponse{}, status.Error(codes.Internal, "could not save contractor")
	}

	log.Debug().Str("id", request.Id).Msg("contractor updated")
	return &proto.UpdateContractorResponse{}, nil
}

func (api *API) DeleteContractor(ctx context.Context, request *proto.DeleteContractorRequest) (*proto.DeleteContractorResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.DeleteContractorResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.DeleteContractorResponse{}, status.Error(codes.FailedPrecondition, "contractor id required")
	}

	err := api.db.DeleteContractor(api.db, account, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.DeleteContractorResponse{}, status.Error(codes.NotFound, "could not delete contractor; contractor key not found")
		}

		log.Error().Err(err).Msg("could not delete contractor")
		return &proto.DeleteContractorResponse{}, status.Error(codes.Internal, "could not delete contractor")
	}

	return &proto.DeleteContractorResponse{}, nil
}
