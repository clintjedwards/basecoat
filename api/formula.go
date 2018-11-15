package api

import (
	"net/http"
	"time"

	"github.com/clintjedwards/basecoat/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//RegisterFormulaRoutes adds /formulas to routing list
func (restAPI *API) RegisterFormulaRoutes(router *mux.Router) {

	router.Handle("/formulas", handlers.MethodHandler{
		"GET":  http.HandlerFunc(restAPI.getFormulasHandler),
		"POST": restAPI.requireAuthentication(http.HandlerFunc(restAPI.createFormulaHandler)),
	})

	router.Handle("/formulas/{id}", handlers.MethodHandler{
		"GET":    http.HandlerFunc(restAPI.getFormulaHandler),
		"PUT":    restAPI.requireAuthentication(http.HandlerFunc(restAPI.updateFormulaHandler)),
		"DELETE": restAPI.requireAuthentication(http.HandlerFunc(restAPI.deleteFormulaHandler)),
	})

}

func (restAPI *API) getFormula(id string) (*models.Formula, error) {
	formula := models.Formula{}

	err := restAPI.db.Model(&formula).Where("id = ?", id).First()
	if formula.ID == 0 {
		return nil, errFormulaNotFound
	} else if err != nil {
		return nil, err
	}

	return &formula, nil
}

func (restAPI *API) getFormulas() ([]*models.Formula, error) {
	formulas := []*models.Formula{}

	err := restAPI.db.Model(&formulas).Order("name ASC").Select()
	if err != nil {
		return nil, err
	}

	return formulas, nil
}

func (restAPI *API) createFormula(newFormula *models.Formula) error {
	newFormula.Created = time.Now().Unix()

	response, err := restAPI.db.Model(newFormula).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}

	if response.RowsAffected() == 0 {
		return errFormulaExists
	}

	return nil
}

func (restAPI *API) updateFormula(id string, updatedFormula *models.Formula) error {

	updatedFormula.Modified = time.Now().Unix()

	_, err := restAPI.db.Model(updatedFormula).Where("id = ?", id).UpdateNotNull()
	if err != nil {
		return err
	}

	return nil
}

func (restAPI *API) deleteFormula(id string) error {

	currentFormula, err := restAPI.getFormula(id)
	if err != nil {
		return err
	}

	err = restAPI.db.Delete(currentFormula)
	if err != nil {
		return err
	}

	return nil
}
