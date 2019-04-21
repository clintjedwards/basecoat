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

//RegisterJobRoutes adds /jobs to routing list
func (restAPI *API) RegisterJobRoutes(router *mux.Router) {

	router.Handle("/jobs", handlers.MethodHandler{
		"GET":  http.HandlerFunc(restAPI.getJobsHandler),
		"POST": restAPI.requireAuthentication(http.HandlerFunc(restAPI.createJobHandler)),
	})

	router.Handle("/jobs/{id}", handlers.MethodHandler{
		"GET":    http.HandlerFunc(restAPI.getJobHandler),
		"PUT":    restAPI.requireAuthentication(http.HandlerFunc(restAPI.updateJobHandler)),
		"DELETE": restAPI.requireAuthentication(http.HandlerFunc(restAPI.deleteJobHandler)),
	})

}

func (restAPI *API) getJob(id string) (*models.Job, error) {
	job := models.Job{}

	err := restAPI.db.Model(&job).Where("id = ?", id).First()
	if job.ID == 0 {
		return nil, errJobNotFound
	} else if err != nil {
		return nil, err
	}

	return &job, nil
}

func (restAPI *API) getJobs() ([]*models.Job, error) {
	jobs := []*models.Job{}

	err := restAPI.db.Model(&jobs).Order("name ASC").Select()
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (restAPI *API) createJob(newJob *models.Job) error {
	newJob.Created = time.Now().Unix()

	response, err := restAPI.db.Model(newJob).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}

	if response.RowsAffected() == 0 {
		return errJobExists
	}

	if newJob.Formulas != nil {

		// Append job id to job list in formula
		for _, formulaID := range newJob.Formulas {
			formula, err := restAPI.getFormula(strconv.Itoa(formulaID))
			if err != nil {
				utils.StructuredLog("error", "could not retrieve formula when attempting to update job list", formulaID)
				continue
			}

			formula.Jobs = append(formula.Jobs, newJob.ID)

			_, err = restAPI.db.Model(formula).Where("id = ?", formulaID).Update()
			if err != nil {
				utils.StructuredLog("error",
					"could not update job list when attempting to update formula",
					map[string]string{"error": err.Error(), "formulaID": strconv.Itoa(formulaID)})
				continue
			}
		}
	}

	return nil
}

func (restAPI *API) updateJob(id string, updatedJob *models.Job) error {

	updatedJob.Modified = time.Now().Unix()

	// If formulas are updated make sure those updates are reflected in all formula objects

	currentJob, _ := restAPI.getJob(id)

	additions, removals := utils.FindListUpdates(currentJob.Formulas, updatedJob.Formulas)

	// Append job id to formula list in job
	for _, formulaID := range additions {
		formula, err := restAPI.getFormula(strconv.Itoa(formulaID))
		if err != nil {
			continue
		}

		formula.Jobs = append(formula.Jobs, currentJob.ID)

		_, err = restAPI.db.Model(formula).Where("id = ?", formulaID).Update()
		if err != nil {
			continue
		}
	}

	// Remove jobs id from jobs list in formula
	for _, formulaID := range removals {
		formula, err := restAPI.getFormula(strconv.Itoa(formulaID))
		if err != nil {
			continue
		}

		formula.Jobs = utils.RemoveIntFromList(formula.Jobs, currentJob.ID)

		_, err = restAPI.db.Model(formula).Where("id = ?", formulaID).Update()
		if err != nil {
			continue
		}
	}

	_, err := restAPI.db.Model(updatedJob).Where("id = ?", id).Update()
	if err != nil {
		return err
	}

	return nil
}

func (restAPI *API) deleteJob(id string) error {

	currentJob, err := restAPI.getJob(id)
	if err != nil {
		return err
	}

	// Remove this job id from all formulas
	for _, formulaID := range currentJob.Formulas {
		updatedFormula, err := restAPI.getFormula(strconv.Itoa(formulaID))
		if err != nil {
			utils.StructuredLog("error", "could not retrieve formula when attempting to update job list", formulaID)
			continue
		}

		updatedJobList := utils.RemoveIntFromList(updatedFormula.Jobs, currentJob.ID)
		updatedFormula.Jobs = updatedJobList

		_, err = restAPI.db.Model(updatedFormula).Where("id = ?", formulaID).Update()
		if err != nil {
			utils.StructuredLog("error",
				"could not update job list when attempting to update formula",
				map[string]string{"error": err.Error(), "formulaID": strconv.Itoa(formulaID)})
			continue
		}
	}

	err = restAPI.db.Delete(currentJob)
	if err != nil {
		return err
	}

	return nil
}
