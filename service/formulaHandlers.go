package service

import (
	"context"
	"time"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/utils"
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
		if err == utils.ErrFormulaNotFound {
			return &api.GetFormulaResponse{}, status.Error(codes.NotFound, "formula requested not found")
		}
		return &api.GetFormulaResponse{}, status.Error(codes.Internal, "failed to retrieve formula from database")
	}

	return &api.GetFormulaResponse{Formula: formula}, nil
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
		Id:        string(utils.GenerateRandString(basecoat.config.Backend.IDLength)),
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

	// If the user has not entered a formula number just make it the ID
	if newFormula.Number == "" {
		newFormula.Number = newFormula.Id
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

	err := basecoat.storage.AddFormula(account, newFormula.Id, &newFormula)
	if err != nil {
		if err == utils.ErrFormulaExists {
			return &api.CreateFormulaResponse{}, status.Error(codes.AlreadyExists, "could not save formula; formula already exists")
		}
		utils.StructuredLog(utils.LogLevelError, "could not save formula", err)
		return &api.CreateFormulaResponse{}, status.Error(codes.Internal, "could not save formula")
	}

	if newFormula.Jobs != nil {
		// Append formula id to formula list in job
		for _, jobID := range newFormula.Jobs {
			job, err := basecoat.storage.GetJob(account, jobID)
			if err != nil {
				utils.StructuredLog(utils.LogLevelWarn, "could not retrieve job when attempting to update formula list", jobID)
				continue
			}

			job.Formulas = append(job.Formulas, newFormula.Id)

			err = basecoat.storage.UpdateJob(account, jobID, job)
			if err != nil {
				utils.StructuredLog(utils.LogLevelError, "could not update job", err)
				continue
			}
		}
	}

	utils.StructuredLog(utils.LogLevelInfo, "formula created", newFormula)

	return &api.CreateFormulaResponse{Id: newFormula.Id}, nil
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
		if err == utils.ErrFormulaNotFound {
			return &api.UpdateFormulaResponse{}, status.Error(codes.NotFound, "could not update formula; formula key not found")
		}
		utils.StructuredLog(utils.LogLevelError, "could not update formula", err)
		return &api.UpdateFormulaResponse{}, status.Error(codes.Internal, "could not update formula")
	}

	additions, removals := utils.FindListUpdates(currentFormula.Jobs, updatedFormula.Jobs)
	// Append formula id to formula list in job
	for _, jobID := range additions {
		job, err := basecoat.storage.GetJob(account, jobID)
		if err != nil {
			utils.StructuredLog(utils.LogLevelWarn, "could not retrieve job when attempting to update formula list", jobID)
			continue
		}

		job.Formulas = append(job.Formulas, currentFormula.Id)

		err = basecoat.storage.UpdateJob(account, jobID, job)
		if err != nil {
			utils.StructuredLog(utils.LogLevelError, "could not update job", err)
			continue
		}
	}

	// Remove formula id from formula list in job
	for _, jobID := range removals {
		job, err := basecoat.storage.GetJob(account, jobID)
		if err != nil {
			utils.StructuredLog(utils.LogLevelWarn, "could not retrieve job when attempting to update formula list", jobID)
			continue
		}

		job.Formulas = utils.RemoveStringFromList(job.Formulas, currentFormula.Id)

		err = basecoat.storage.UpdateJob(account, jobID, job)
		if err != nil {
			utils.StructuredLog(utils.LogLevelError, "could not update job", err)
			continue
		}
	}

	utils.StructuredLog(utils.LogLevelInfo, "formula updated", updatedFormula)

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

	// Remove this formula id from all jobs
	currentFormula, _ := basecoat.storage.GetFormula(account, request.Id)
	for _, jobID := range currentFormula.Jobs {
		job, err := basecoat.storage.GetJob(account, jobID)
		if err != nil {
			utils.StructuredLog(utils.LogLevelWarn, "could not retrieve job when attempting to update formula list", jobID)
			continue
		}

		updatedFormulaList := utils.RemoveStringFromList(job.Formulas, currentFormula.Id)
		job.Formulas = updatedFormulaList

		err = basecoat.storage.UpdateJob(account, jobID, job)
		if err != nil {
			utils.StructuredLog(utils.LogLevelWarn, "could not update job", err)
			continue
		}
	}

	err := basecoat.storage.DeleteFormula(account, request.Id)
	if err != nil {
		if err == utils.ErrFormulaNotFound {
			return &api.DeleteFormulaResponse{}, status.Error(codes.NotFound, "could not delete formula; formula key not found")
		}
		utils.StructuredLog(utils.LogLevelError, "could not delete formula", err)
		return &api.DeleteFormulaResponse{}, status.Error(codes.Internal, "could not delete formula")
	}

	utils.StructuredLog(utils.LogLevelInfo, "formula deleted", request.Id)

	return &api.DeleteFormulaResponse{}, nil
}
