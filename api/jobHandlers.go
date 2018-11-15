package api

import (
	"fmt"
	"net/http"

	"github.com/clintjedwards/basecoat/models"
	"github.com/clintjedwards/basecoat/utils"
	"github.com/gorilla/mux"
)

func (restAPI *API) getJobHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	job, err := restAPI.getJob(vars["id"])
	if err == errJobNotFound {
		utils.SendResponse(w, http.StatusNotFound, fmt.Sprintf("%s: %s", errJobNotFound.Error(), vars["id"]), true)
		return
	} else if err != nil {
		utils.StructuredLog("error", "could not retrieve job", err)
		utils.SendResponse(w, http.StatusInternalServerError, "could not retrieve job", true)
		return
	}

	utils.SendResponse(w, http.StatusOK, job, false)
	return
}

func (restAPI *API) getJobsHandler(w http.ResponseWriter, req *http.Request) {
	jobs, err := restAPI.getJobs()
	if err != nil {
		utils.StructuredLog("error", "could not retrieve jobs", err)
		utils.SendResponse(w, http.StatusInternalServerError, "could not retrieve jobs", true)
		return
	}

	utils.SendResponse(w, http.StatusOK, jobs, false)
	return
}

func (restAPI *API) createJobHandler(w http.ResponseWriter, req *http.Request) {
	newJob := models.Job{}
	err := utils.ParseJSON(req.Body, &newJob)
	if err != nil {
		utils.StructuredLog("error", "could not decode json body", err)
		utils.SendResponse(w, http.StatusBadRequest, "could not decode json body", true)
		return
	}
	defer req.Body.Close()

	err = restAPI.createJob(&newJob)
	if err == errJobExists {
		utils.SendResponse(w, http.StatusConflict, fmt.Sprintf("job %s already exists", newJob.Name), true)
		return
	} else if err != nil {
		utils.StructuredLog("error", "could not create job", err)
		utils.SendResponse(w, http.StatusInternalServerError, "could not create job", true)
		return
	}

	utils.StructuredLog("info", "job registered", newJob)
	utils.SendResponse(w, http.StatusCreated, "", false)
	return
}

func (restAPI *API) updateJobHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	var updatedJob models.Job

	err := utils.ParseJSON(req.Body, &updatedJob)
	if err != nil {
		utils.StructuredLog("error", "could not parse json", err)
		utils.SendResponse(w, http.StatusBadRequest, errJSONParseFailure, true)
		return
	}
	defer req.Body.Close()

	err = restAPI.updateJob(vars["id"], &updatedJob)
	if err == errJobNotFound {
		utils.SendResponse(w, http.StatusNotFound, fmt.Sprintf("%s: %s", errJobNotFound.Error(), vars["id"]), true)
		return
	} else if err != nil {
		utils.StructuredLog("error", "could not update jobs", err)
		utils.SendResponse(w, http.StatusInternalServerError, "could not update job", true)
		return
	}

	utils.StructuredLog("info", "job updated", updatedJob)
	utils.SendResponse(w, http.StatusNoContent, "", false)
	return
}

func (restAPI *API) deleteJobHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	err := restAPI.deleteJob(vars["id"])
	if err == errJobNotFound {
		utils.SendResponse(w, http.StatusNotFound, fmt.Sprintf("%s: %s", errJobNotFound.Error(), vars["id"]), true)
		return
	} else if err != nil {
		utils.StructuredLog("error", "could not delete job", err)
		utils.SendResponse(w, http.StatusInternalServerError, "could not delete job", true)
		return
	}

	utils.StructuredLog("info", "job deleted", vars["id"])
	utils.SendResponse(w, http.StatusNoContent, "", false)

	return
}
