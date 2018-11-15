package api

import (
	"fmt"
	"net/http"

	"github.com/clintjedwards/basecoat/models"
	"github.com/clintjedwards/basecoat/utils"
	"github.com/gorilla/mux"
)

func (restAPI *API) getFormulaHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	formula, err := restAPI.getFormula(vars["id"])
	if err == errFormulaNotFound {
		utils.SendResponse(w, http.StatusNotFound, fmt.Sprintf("%s: %s", errFormulaNotFound.Error(), vars["id"]), true)
		return
	} else if err != nil {
		utils.StructuredLog("error", "could not retrieve formula", err)
		utils.SendResponse(w, http.StatusInternalServerError, "could not retrieve formula", true)
		return
	}

	utils.SendResponse(w, http.StatusOK, formula, false)
	return
}

func (restAPI *API) getFormulasHandler(w http.ResponseWriter, req *http.Request) {
	formulas, err := restAPI.getFormulas()
	if err != nil {
		utils.StructuredLog("error", "could not retrieve formulas", err)
		utils.SendResponse(w, http.StatusInternalServerError, "could not retrieve formulas", true)
		return
	}

	utils.SendResponse(w, http.StatusOK, formulas, false)
	return
}

func (restAPI *API) createFormulaHandler(w http.ResponseWriter, req *http.Request) {
	newFormula := models.Formula{}
	err := utils.ParseJSON(req.Body, &newFormula)
	if err != nil {
		utils.StructuredLog("error", "could not decode json body", err)
		utils.SendResponse(w, http.StatusBadRequest, "could not decode json body", true)
		return
	}
	defer req.Body.Close()

	err = restAPI.createFormula(&newFormula)
	if err == errFormulaExists {
		utils.SendResponse(w, http.StatusConflict, fmt.Sprintf("formula %s:%s already exists", newFormula.Name, newFormula.Number), true)
		return
	} else if err != nil {
		utils.StructuredLog("error", "could not create formula", err)
		utils.SendResponse(w, http.StatusInternalServerError, "could not create formula", true)
		return
	}

	utils.StructuredLog("info", "formula registered", newFormula)
	utils.SendResponse(w, http.StatusCreated, "", false)
	return
}

func (restAPI *API) updateFormulaHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	var updatedFormula models.Formula

	err := utils.ParseJSON(req.Body, &updatedFormula)
	if err != nil {
		utils.StructuredLog("error", "could not parse json", err)
		utils.SendResponse(w, http.StatusBadRequest, errJSONParseFailure, true)
		return
	}
	defer req.Body.Close()

	err = restAPI.updateFormula(vars["id"], &updatedFormula)
	if err == errFormulaNotFound {
		utils.SendResponse(w, http.StatusNotFound, fmt.Sprintf("%s: %s", errFormulaNotFound.Error(), vars["id"]), true)
		return
	} else if err != nil {
		utils.StructuredLog("error", "could not update formula", err)
		utils.SendResponse(w, http.StatusInternalServerError, "could not update formula", true)
		return
	}

	utils.StructuredLog("info", "formula updated", updatedFormula)
	utils.SendResponse(w, http.StatusNoContent, "", false)
	return
}

func (restAPI *API) deleteFormulaHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	err := restAPI.deleteFormula(vars["id"])
	if err == errFormulaNotFound {
		utils.SendResponse(w, http.StatusNotFound, fmt.Sprintf("%s: %s", errFormulaNotFound.Error(), vars["id"]), true)
		return
	} else if err != nil {
		utils.StructuredLog("error", "could not delete formula", err)
		utils.SendResponse(w, http.StatusInternalServerError, "could not delete formula", true)
		return
	}

	utils.StructuredLog("info", "formula deleted", vars["id"])
	utils.SendResponse(w, http.StatusNoContent, "", false)

	return
}
