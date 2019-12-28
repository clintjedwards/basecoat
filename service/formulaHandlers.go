package service

import (
	"context"
	"fmt"
	"time"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/toolkit/logger"
	"github.com/clintjedwards/toolkit/tkerrors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetFormula returns a single formula by key
func (basecoat *API) GetFormula(context context.Context, request *api.GetFormulaRequest) (*api.GetFormulaResponse, error) {

	account, present := getAccountFromContext(context)
	if !present {
		return &api.GetFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &api.GetFormulaResponse{}, status.Error(codes.FailedPrecondition, "formula id required")
	}

	formula, err := basecoat.storage.GetFormula(account, request.Id)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.GetFormulaResponse{}, status.Error(codes.NotFound, "formula requested not found")
		}
		return &api.GetFormulaResponse{}, status.Error(codes.Internal, "failed to retrieve formula from database")
	}

	return &api.GetFormulaResponse{Formula: formula}, nil
}

// SearchFormulas takes in a search term and returns formulas that might match
func (basecoat *API) SearchFormulas(context context.Context, request *api.SearchFormulasRequest) (*api.SearchFormulasResponse, error) {

	account, present := getAccountFromContext(context)
	if !present {
		return &api.SearchFormulasResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Term == "" {
		return &api.SearchFormulasResponse{}, status.Error(codes.FailedPrecondition, "search term required")
	}

	searchResults, err := basecoat.search.SearchFormulas(account, request.Term)
	if err != nil {
		return &api.SearchFormulasResponse{}, status.Error(codes.Internal, fmt.Sprintf("a search error occurred: %v", err))
	}

	return &api.SearchFormulasResponse{Results: searchResults}, nil
}

// ListFormulas returns a list of all formulas on basecoat service
func (basecoat *API) ListFormulas(context context.Context, request *api.ListFormulasRequest) (*api.ListFormulasResponse, error) {

	account, present := getAccountFromContext(context)
	if !present {
		return &api.ListFormulasResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	formulas, err := basecoat.storage.GetAllFormulas(account)
	if err != nil {
		return &api.ListFormulasResponse{}, status.Error(codes.Internal, "failed to retrieve formulas from database")
	}

	return &api.ListFormulasResponse{Formulas: formulas}, nil
}

// CreateFormula registers a new formula
func (basecoat *API) CreateFormula(context context.Context, request *api.CreateFormulaRequest) (*api.CreateFormulaResponse, error) {

	account, present := getAccountFromContext(context)
	if !present {
		return &api.CreateFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	newFormula := api.Formula{
		Name:      request.Name,
		Number:    request.Number,
		Notes:     request.Notes,
		Created:   time.Now().Unix(),
		Modified:  time.Now().Unix(),
		Jobs:      request.Jobs,
		Bases:     request.Bases,
		Colorants: request.Colorants,
	}

	if newFormula.Name == "" {
		return &api.CreateFormulaResponse{}, status.Error(codes.FailedPrecondition, "name required")
	}

	for _, base := range newFormula.Bases {
		if base.Name == "" {
			return &api.CreateFormulaResponse{}, status.Error(codes.FailedPrecondition, "base name required")
		}
	}

	for _, colorant := range newFormula.Colorants {
		if colorant.Name == "" {
			return &api.CreateFormulaResponse{}, status.Error(codes.FailedPrecondition, "colorant name required")
		}
	}

	formulaID, err := basecoat.storage.AddFormula(account, &newFormula)
	if err != nil {
		if err == tkerrors.ErrEntityExists {
			return &api.CreateFormulaResponse{}, status.Error(codes.AlreadyExists, "could not save formula; formula already exists")
		}
		logger.Log().Errorw("could not save formula", "error", err)
		return &api.CreateFormulaResponse{}, status.Error(codes.Internal, "could not save formula")
	}

	formula, err := basecoat.storage.GetFormula(account, formulaID)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.CreateFormulaResponse{}, status.Error(codes.NotFound, "could not retrieve formula after saving")
		}
		return &api.CreateFormulaResponse{}, status.Error(codes.Internal, "could not retrieve formula after saving")
	}

	basecoat.search.UpdateFormulaIndex(account, *formula)

	logger.Log().Infow("formula created", "formula", *formula)
	return &api.CreateFormulaResponse{Id: formula.Id}, nil
}

// UpdateFormula updates an already existing formula
func (basecoat *API) UpdateFormula(context context.Context, request *api.UpdateFormulaRequest) (*api.UpdateFormulaResponse, error) {

	account, present := getAccountFromContext(context)
	if !present {
		return &api.UpdateFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &api.UpdateFormulaResponse{}, status.Error(codes.FailedPrecondition, "formula id required")
	}

	currentFormula, _ := basecoat.storage.GetFormula(account, request.Id)

	updatedFormula := api.Formula{
		Id:        request.Id,
		Name:      request.Name,
		Number:    request.Number,
		Notes:     request.Notes,
		Created:   currentFormula.Created,
		Modified:  time.Now().Unix(),
		Jobs:      request.Jobs,
		Bases:     request.Bases,
		Colorants: request.Colorants,
	}

	err := basecoat.storage.UpdateFormula(account, request.Id, &updatedFormula)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.UpdateFormulaResponse{}, status.Error(codes.NotFound, "could not update formula; formula key not found")
		}
		logger.Log().Errorw("could not update formula", "error", err)
		return &api.UpdateFormulaResponse{}, status.Error(codes.Internal, "could not update formula")
	}

	basecoat.search.UpdateFormulaIndex(account, updatedFormula)

	logger.Log().Infow("formula updated", "formula", updatedFormula)
	return &api.UpdateFormulaResponse{}, nil
}

// DeleteFormula removes a formula
func (basecoat *API) DeleteFormula(context context.Context, request *api.DeleteFormulaRequest) (*api.DeleteFormulaResponse, error) {

	account, present := getAccountFromContext(context)
	if !present {
		return &api.DeleteFormulaResponse{}, status.Error(codes.FailedPrecondition, "account required")
	}

	if request.Id == "" {
		return &api.DeleteFormulaResponse{}, status.Error(codes.FailedPrecondition, "formula id required")
	}

	err := basecoat.storage.DeleteFormula(account, request.Id)
	if err != nil {
		if err == tkerrors.ErrEntityNotFound {
			return &api.DeleteFormulaResponse{}, status.Error(codes.NotFound, "could not delete formula; formula key not found")
		}
		logger.Log().Errorw("could not delete formula", "error", err)
		return &api.DeleteFormulaResponse{}, status.Error(codes.Internal, "could not delete formula")
	}

	basecoat.search.DeleteFormulaIndex(account, request.Id)

	logger.Log().Infow("formula deleted", "id", request.Id)
	return &api.DeleteFormulaResponse{}, nil
}
