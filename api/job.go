package api

import (
	"net/http"
	"time"

	"github.com/clintjedwards/basecoat/models"
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

	return nil
}

func (restAPI *API) updateJob(id string, updatedJob *models.Job) error {

	updatedJob.Modified = time.Now().Unix()

	_, err := restAPI.db.Model(updatedJob).Where("id = ?", id).UpdateNotNull()
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

	err = restAPI.db.Delete(currentJob)
	if err != nil {
		return err
	}

	return nil
}
