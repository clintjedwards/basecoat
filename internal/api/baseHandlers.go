package api

import (
	"context"

	"github.com/clintjedwards/basecoat/internal/models"
	"github.com/clintjedwards/basecoat/internal/storage"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetBase returns a single base by key
func (api *API) GetBase(ctx context.Context, request *proto.GetBaseRequest) (*proto.GetBaseResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.GetBaseResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.GetBaseResponse{}, status.Error(codes.FailedPrecondition, "id required")
	}

	base := models.Base{}

	err := storage.InsideTx(api.db.DB, func(tx *sqlx.Tx) error {
		baseRaw, err := api.db.GetBase(tx, account, request.Id)
		if err != nil {
			return err
		}

		baseMetadata := models.BaseMetadata{}
		baseMetadata.FromStorage(&baseRaw)
		base.Metadata = baseMetadata

		formulaBases, err := api.db.ListBaseFormulas(tx, account, request.Id)
		if err != nil {
			return err
		}

		formulaSet := map[string]struct{}{}

		for _, formulaBase := range formulaBases {
			formulaSet[formulaBase.Formula] = struct{}{}
		}

		base.FormulaIDs = maps.Keys(formulaSet)

		return nil
	})
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.GetBaseResponse{}, status.Error(codes.NotFound, "base requested not found")
		}
		return &proto.GetBaseResponse{}, status.Error(codes.Internal, "failed to retrieve base from database")
	}

	return &proto.GetBaseResponse{Base: base.ToProto()}, nil
}

// ListBases returns a list of all bases' metadata
func (api *API) ListBases(ctx context.Context, _ *proto.ListBasesRequest) (*proto.ListBasesResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.ListBasesResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	basesRaw, err := api.db.ListBases(api.db, account, 0, 0)
	if err != nil {
		return &proto.ListBasesResponse{}, status.Error(codes.Internal, "failed to retrieve bases from database")
	}

	protoBases := []*proto.BaseMetadata{}
	for _, baseRaw := range basesRaw {
		var baseMetadata models.BaseMetadata
		baseMetadata.FromStorage(&baseRaw)
		protoBases = append(protoBases, baseMetadata.ToProto())
	}

	return &proto.ListBasesResponse{Bases: protoBases}, nil
}

// CreateBase registers a new base
func (api *API) CreateBase(ctx context.Context, request *proto.CreateBaseRequest) (*proto.CreateBaseResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.CreateBaseResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Label == "" {
		return &proto.CreateBaseResponse{}, status.Error(codes.FailedPrecondition, "base label required")
	}

	if request.Manufacturer == "" {
		return &proto.CreateBaseResponse{}, status.Error(codes.FailedPrecondition, "base manufacturer required")
	}

	base := models.NewBaseMetadata(account, request.Label, request.Manufacturer)

	err := api.db.InsertBase(api.db, base.ToStorage())
	if err != nil {
		if err == storage.ErrEntityExists {
			return &proto.CreateBaseResponse{}, status.Error(codes.AlreadyExists, "could not save base; base already exists")
		}
		log.Error().Err(err).Msg("could not save base")
		return &proto.CreateBaseResponse{}, status.Error(codes.Internal, "could not save base")
	}

	log.Info().Str("id", base.ID).Str("label", base.Label).Str("manufacturer", base.Manufacturer).Msg("base created")

	return &proto.CreateBaseResponse{
		Base: base.ToProto(),
	}, nil
}

// UpdateBase updates an already existing base
func (api *API) UpdateBase(ctx context.Context, request *proto.UpdateBaseRequest) (*proto.UpdateBaseResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.UpdateBaseResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.UpdateBaseResponse{}, status.Error(codes.FailedPrecondition, "base id required")
	}

	err := api.db.UpdateBase(api.db, account, request.Id, storage.UpdatableBaseFields{
		Label:        &request.Label,
		Manufacturer: &request.Manufacturer,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not save base")
		return &proto.UpdateBaseResponse{}, status.Error(codes.Internal, "could not save base")
	}

	log.Debug().Str("id", request.Id).Msg("base updated")
	return &proto.UpdateBaseResponse{}, nil
}

// AssociateBaseWithFormula records the base as belonging to the formula given.
func (api *API) AssociateBaseWithFormula(ctx context.Context, request *proto.AssociateBaseWithFormulaRequest) (*proto.AssociateBaseWithFormulaResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.AssociateBaseWithFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Formula == "" {
		return &proto.AssociateBaseWithFormulaResponse{}, status.Error(codes.FailedPrecondition, "formula id required")
	}

	if request.Base == "" {
		return &proto.AssociateBaseWithFormulaResponse{}, status.Error(codes.FailedPrecondition, "base id required")
	}

	if request.Amount == "" {
		return &proto.AssociateBaseWithFormulaResponse{}, status.Error(codes.FailedPrecondition, "base amount required")
	}

	err := api.db.AssociateBaseWithFormula(api.db, &storage.FormulaBase{
		Account: account,
		Formula: request.Formula,
		Base:    request.Base,
		Amount:  request.Amount,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not attach base to formula")
		return &proto.AssociateBaseWithFormulaResponse{}, status.Error(codes.Internal, "could not attach base to formula")
	}

	log.Debug().Str("formula", request.Formula).Str("base", request.Base).Msg("base added to formula")
	return &proto.AssociateBaseWithFormulaResponse{}, nil
}

// DisassociateBaseFromFormula removes a base entry from a formula.
func (api *API) DisassociateBaseFromFormula(ctx context.Context, request *proto.DisassociateBaseFromFormulaRequest) (*proto.DisassociateBaseFromFormulaResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.DisassociateBaseFromFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Formula == "" {
		return &proto.DisassociateBaseFromFormulaResponse{}, status.Error(codes.FailedPrecondition, "formula id required")
	}

	if request.Base == "" {
		return &proto.DisassociateBaseFromFormulaResponse{}, status.Error(codes.FailedPrecondition, "base id required")
	}

	err := api.db.DeleteFormulaBase(api.db, account, request.Formula, request.Base)
	if err != nil {
		log.Error().Err(err).Msg("could not remove base from Formula")
		return &proto.DisassociateBaseFromFormulaResponse{}, status.Error(codes.Internal, "could not remove base from Formula")
	}

	log.Debug().Str("formula", request.Formula).Str("base", request.Base).Msg("base removed to formula")
	return &proto.DisassociateBaseFromFormulaResponse{}, nil
}

func (api *API) DeleteBase(ctx context.Context, request *proto.DeleteBaseRequest) (*proto.DeleteBaseResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.DeleteBaseResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.DeleteBaseResponse{}, status.Error(codes.FailedPrecondition, "base id required")
	}

	err := api.db.DeleteBase(api.db, account, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.DeleteBaseResponse{}, status.Error(codes.NotFound, "could not delete Base; base key not found")
		}

		log.Error().Err(err).Msg("could not delete Base")
		return &proto.DeleteBaseResponse{}, status.Error(codes.Internal, "could not delete Base")
	}

	return &proto.DeleteBaseResponse{}, nil
}
