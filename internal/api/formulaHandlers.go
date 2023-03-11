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

// GetFormula returns a single formula by key
func (api *API) GetFormula(ctx context.Context, request *proto.GetFormulaRequest) (*proto.GetFormulaResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.GetFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.GetFormulaResponse{}, status.Error(codes.FailedPrecondition, "id required")
	}

	formulaRaw, err := api.db.GetFormula(api.db, account, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.GetFormulaResponse{}, status.Error(codes.NotFound, "formula requested not found")
		}
		return &proto.GetFormulaResponse{}, status.Error(codes.Internal, "failed to retrieve formula from database")
	}

	formula := models.FormulaMetadata{}
	formula.FromStorage(&formulaRaw)

	return &proto.GetFormulaResponse{Formula: formula.ToProto()}, nil
}

// ListFormulas returns a list of all formulas's metadata.
func (api *API) ListFormulas(ctx context.Context, request *proto.ListFormulasRequest) (*proto.ListFormulasResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.ListFormulasResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	searchResults := []string{}

	if request.Filter != "" {
		results, err := api.search.SearchFormulas(account, request.Filter)
		if err != nil {
			log.Error().Err(err).Msg("a search error occurred")
		}
		searchResults = results
		log.Debug().Strs("returned_results", searchResults).Str("query", request.Filter).Msg("filtered formulas on user's request")
	}

	formulasRaw, err := api.db.ListFormulas(api.db, account, 0, 0)
	if err != nil {
		return &proto.ListFormulasResponse{}, status.Error(codes.Internal, "failed to retrieve formulas from database")
	}

	protoFormulas := []*proto.FormulaMetadata{}
	for _, formulaRaw := range formulasRaw {
		// If the user actually enters a filter we want to use that filter.
		if request.Filter != "" {
			// If the formula doesn't exist in the filter then we should skip it.
			if !slices.Contains(searchResults, formulaRaw.ID) {
				continue
			}
		}

		var formulaMetadata models.FormulaMetadata
		formulaMetadata.FromStorage(&formulaRaw)
		protoFormulas = append(protoFormulas, formulaMetadata.ToProto())
	}

	return &proto.ListFormulasResponse{Formulas: protoFormulas}, nil
}

// CreateFormula registers a new formula
func (api *API) CreateFormula(ctx context.Context, request *proto.CreateFormulaRequest) (*proto.CreateFormulaResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.CreateFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Name == "" {
		return &proto.CreateFormulaResponse{}, status.Error(codes.FailedPrecondition, "formula name required")
	}

	formula := models.NewFormulaMetadata(account, request.Name)
	formula.Number = request.Number
	formula.Notes = request.Notes

	err := api.db.InsertFormula(api.db, formula.ToStorage())
	if err != nil {
		if err == storage.ErrEntityExists {
			return &proto.CreateFormulaResponse{}, status.Error(codes.AlreadyExists, "could not save formula; formula already exists")
		}
		log.Error().Err(err).Msg("could not save formula")
		return &proto.CreateFormulaResponse{}, status.Error(codes.Internal, "could not save formula")
	}

	log.Info().Str("id", formula.ID).Str("name", formula.Name).Msg("formula created")

	return &proto.CreateFormulaResponse{
		Formula: formula.ToProto(),
	}, nil
}

// UpdateFormula updates an already existing formula
func (api *API) UpdateFormula(ctx context.Context, request *proto.UpdateFormulaRequest) (*proto.UpdateFormulaResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.UpdateFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.UpdateFormulaResponse{}, status.Error(codes.FailedPrecondition, "formula id required")
	}

	err := api.db.UpdateFormula(api.db, account, request.Id, storage.UpdatableFormulaFields{
		Name:     &request.Name,
		Number:   &request.Number,
		Notes:    &request.Notes,
		Modified: ptr(time.Now().UnixMilli()),
	})
	if err != nil {
		if err == storage.ErrEntityExists {
			return &proto.UpdateFormulaResponse{}, status.Error(codes.AlreadyExists, "could not save formula; formula already exists")
		}
		log.Error().Err(err).Msg("could not save formula")
		return &proto.UpdateFormulaResponse{}, status.Error(codes.Internal, "could not save formula")
	}

	log.Debug().Str("id", request.Id).Msg("formula updated")
	return &proto.UpdateFormulaResponse{}, nil
}

func (api *API) DeleteFormula(ctx context.Context, request *proto.DeleteFormulaRequest) (*proto.DeleteFormulaResponse, error) {
	account, present := getAccountFromContext(ctx)
	if !present {
		return &proto.DeleteFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &proto.DeleteFormulaResponse{}, status.Error(codes.FailedPrecondition, "formula id required")
	}

	err := api.db.DeleteFormula(api.db, account, request.Id)
	if err != nil {
		if err == storage.ErrEntityNotFound {
			return &proto.DeleteFormulaResponse{}, status.Error(codes.NotFound, "could not delete formula; formula key not found")
		}

		log.Error().Err(err).Msg("could not delete formula")
		return &proto.DeleteFormulaResponse{}, status.Error(codes.Internal, "could not delete formula")
	}

	return &proto.DeleteFormulaResponse{}, nil
}
