package api

import (
	"context"

	"github.com/clintjedwards/basecoat/internal/models"
	"github.com/clintjedwards/basecoat/internal/storage"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetColorant returns a single colorant by key
func (api *API) GetColorant(ctx context.Context, request *proto.GetColorantRequest) (*proto.GetColorantResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.GetColorantResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.GetColorantResponse{}, status.Error(codes.FailedPrecondition, "id required")
	}

	colorantRaw, err := api.db.GetColorant(api.db, account, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.GetColorantResponse{}, status.Error(codes.NotFound, "colorant requested not found")
		}
		return &proto.GetColorantResponse{}, status.Error(codes.Internal, "failed to retrieve colorant from datacolorant")
	}

	colorant := models.ColorantMetadata{}
	colorant.FromStorage(&colorantRaw)

	return &proto.GetColorantResponse{Colorant: colorant.ToProto()}, nil
}

// ListColorants returns a list of all colorants' metadata
func (api *API) ListColorants(ctx context.Context, _ *proto.ListColorantsRequest) (*proto.ListColorantsResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.ListColorantsResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	colorantsRaw, err := api.db.ListColorants(api.db, account, 0, 0)
	if err != nil {
		return &proto.ListColorantsResponse{}, status.Error(codes.Internal, "failed to retrieve colorants from datacolorant")
	}

	protoColorants := []*proto.ColorantMetadata{}
	for _, colorantRaw := range colorantsRaw {
		var colorantMetadata models.ColorantMetadata
		colorantMetadata.FromStorage(&colorantRaw)
		protoColorants = append(protoColorants, colorantMetadata.ToProto())
	}

	return &proto.ListColorantsResponse{Colorants: protoColorants}, nil
}

// CreateColorant registers a new colorant
func (api *API) CreateColorant(ctx context.Context, request *proto.CreateColorantRequest) (*proto.CreateColorantResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.CreateColorantResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Label == "" {
		return &proto.CreateColorantResponse{}, status.Error(codes.FailedPrecondition, "colorant label required")
	}

	if request.Manufacturer == "" {
		return &proto.CreateColorantResponse{}, status.Error(codes.FailedPrecondition, "colorant manufacturer required")
	}

	colorant := models.NewColorantMetadata(account, request.Label, request.Manufacturer)

	err := api.db.InsertColorant(api.db, colorant.ToStorage())
	if err != nil {
		if err == storage.ErrEntityExists {
			return &proto.CreateColorantResponse{}, status.Error(codes.AlreadyExists, "could not save colorant; colorant already exists")
		}
		log.Error().Err(err).Msg("could not save colorant")
		return &proto.CreateColorantResponse{}, status.Error(codes.Internal, "could not save colorant")
	}

	log.Info().Str("id", colorant.ID).Str("label", colorant.Label).Str("manufacturer", colorant.Manufacturer).Msg("colorant created")

	return &proto.CreateColorantResponse{
		Colorant: colorant.ToProto(),
	}, nil
}

// UpdateColorant updates an already existing colorant
func (api *API) UpdateColorant(ctx context.Context, request *proto.UpdateColorantRequest) (*proto.UpdateColorantResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.UpdateColorantResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.UpdateColorantResponse{}, status.Error(codes.FailedPrecondition, "colorant id required")
	}

	err := api.db.UpdateColorant(api.db, account, request.Id, storage.UpdatableColorantFields{
		Label:        &request.Label,
		Manufacturer: &request.Manufacturer,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not save colorant")
		return &proto.UpdateColorantResponse{}, status.Error(codes.Internal, "could not save colorant")
	}

	log.Debug().Str("id", request.Id).Msg("colorant updated")
	return &proto.UpdateColorantResponse{}, nil
}

// AssociateColorantWithFormula records the colorant as belonging to the formula given.
func (api *API) AssociateColorantWithFormula(ctx context.Context, request *proto.AssociateColorantWithFormulaRequest) (*proto.AssociateColorantWithFormulaResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.AssociateColorantWithFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Formula == "" {
		return &proto.AssociateColorantWithFormulaResponse{}, status.Error(codes.FailedPrecondition, "formula id required")
	}

	if request.Colorant == "" {
		return &proto.AssociateColorantWithFormulaResponse{}, status.Error(codes.FailedPrecondition, "colorant id required")
	}

	if request.Amount == "" {
		return &proto.AssociateColorantWithFormulaResponse{}, status.Error(codes.FailedPrecondition, "colorant amount required")
	}

	err := api.db.AssociateColorantWithFormula(api.db, &storage.FormulaColorant{
		Account:  account,
		Formula:  request.Formula,
		Colorant: request.Colorant,
		Amount:   request.Amount,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not attach colorant to formula")
		return &proto.AssociateColorantWithFormulaResponse{}, status.Error(codes.Internal, "could not attach colorant to formula")
	}

	log.Debug().Str("formula", request.Formula).Str("colorant", request.Colorant).Msg("colorant added to formula")
	return &proto.AssociateColorantWithFormulaResponse{}, nil
}

// DisassociateColorantFromFormula removes a colorant entry from a formula.
func (api *API) DisassociateColorantFromFormula(ctx context.Context, request *proto.DisassociateColorantFromFormulaRequest) (*proto.DisassociateColorantFromFormulaResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.DisassociateColorantFromFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Formula == "" {
		return &proto.DisassociateColorantFromFormulaResponse{}, status.Error(codes.FailedPrecondition, "formula id required")
	}

	if request.Colorant == "" {
		return &proto.DisassociateColorantFromFormulaResponse{}, status.Error(codes.FailedPrecondition, "colorant id required")
	}

	err := api.db.DeleteFormulaColorant(api.db, account, request.Formula, request.Colorant)
	if err != nil {
		log.Error().Err(err).Msg("could not remove colorant from Formula")
		return &proto.DisassociateColorantFromFormulaResponse{}, status.Error(codes.Internal, "could not remove colorant from Formula")
	}

	log.Debug().Str("formula", request.Formula).Str("colorant", request.Colorant).Msg("colorant removed to formula")
	return &proto.DisassociateColorantFromFormulaResponse{}, nil
}

func (api *API) DeleteColorant(ctx context.Context, request *proto.DeleteColorantRequest) (*proto.DeleteColorantResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.DeleteColorantResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.DeleteColorantResponse{}, status.Error(codes.FailedPrecondition, "colorant id required")
	}

	err := api.db.DeleteColorant(api.db, account, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.DeleteColorantResponse{}, status.Error(codes.NotFound, "could not delete Colorant; colorant key not found")
		}

		log.Error().Err(err).Msg("could not delete Colorant")
		return &proto.DeleteColorantResponse{}, status.Error(codes.Internal, "could not delete Colorant")
	}

	return &proto.DeleteColorantResponse{}, nil
}
