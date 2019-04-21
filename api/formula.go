package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/clintjedwards/basecoat/models"
	"github.com/clintjedwards/basecoat/utils"
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

	if newFormula.Jobs != nil {

		// Append formula id to formula list in job
		for _, jobID := range newFormula.Jobs {
			job, err := restAPI.getJob(strconv.Itoa(jobID))
			if err != nil {
				utils.StructuredLog("error", "could not retrieve job when attempting to update formula list", jobID)
				continue
			}

			job.Formulas = append(job.Formulas, newFormula.ID)

			_, err = restAPI.db.Model(job).Where("id = ?", jobID).Update()
			if err != nil {
				utils.StructuredLog("error",
					"could not update formula list when attempting to job",
					map[string]string{"error": err.Error(), "jobID": strconv.Itoa(jobID)})
				continue
			}
		}
	}

	return nil
}

func (restAPI *API) updateFormula(id string, updatedFormula *models.Formula) error {

	updatedFormula.Modified = time.Now().Unix()

	currentFormula, _ := restAPI.getFormula(id)
	additions, removals := utils.FindListUpdates(currentFormula.Jobs, updatedFormula.Jobs)

	// Append formula id to formula list in job
	for _, jobID := range additions {
		job, err := restAPI.getJob(strconv.Itoa(jobID))
		if err != nil {
			utils.StructuredLog("error", "could not retrieve job when attempting to update formula list", jobID)
			continue
		}

		job.Formulas = append(job.Formulas, currentFormula.ID)

		_, err = restAPI.db.Model(job).Where("id = ?", jobID).Update()
		if err != nil {
			continue
		}
	}

	// Remove formula id from formula list in job
	for _, jobID := range removals {
		job, err := restAPI.getJob(strconv.Itoa(jobID))
		if err != nil {
			utils.StructuredLog("error", "could not retrieve job when attempting to update formula list", jobID)
			continue
		}

		job.Formulas = utils.RemoveIntFromList(job.Formulas, currentFormula.ID)

		_, err = restAPI.db.Model(job).Where("id = ?", jobID).Update()
		if err != nil {
			continue
		}
	}

	_, err := restAPI.db.Model(updatedFormula).Where("id = ?", id).Update()
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

	// Remove this formula id from all jobs
	for _, jobID := range currentFormula.Jobs {
		updatedJob, err := restAPI.getJob(strconv.Itoa(jobID))
		if err != nil {
			continue
		}

		updatedFormulaList := utils.RemoveIntFromList(updatedJob.Formulas, currentFormula.ID)
		updatedJob.Formulas = updatedFormulaList

		_, err = restAPI.db.Model(updatedJob).Where("id = ?", currentFormula.ID).Update()
		if err != nil {
			continue
		}
	}

	err = restAPI.db.Delete(currentFormula)
	if err != nil {
		return err
	}

	return nil
}
