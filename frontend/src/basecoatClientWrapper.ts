import {
    CreateAPITokenRequest,
    ListFormulasRequest,
    Formula, Job,
    ListJobsRequest,
    CreateFormulaRequest,
    CreateJobRequest,
    GetFormulaRequest,
    GetJobRequest,
    UpdateFormulaRequest,
    UpdateJobRequest,
    DeleteFormulaRequest,
    DeleteJobRequest,
    Base,
    Colorant,
    Contact
} from "./basecoat_pb"

import Cookies from 'js-cookie'
import { BasecoatClient } from "./BasecoatServiceClientPb"

interface LoginInfo {
    username: string;
    password: string;
}

interface formulaMap { [key: string]: Formula }
interface jobMap { [key: string]: Job }

declare var __API__: string; // The api endpoint that the client will talk to

// BasecoatClientWrapper is a wrapper for all frontend to backend communication
class BasecoatClientWrapper {

    client: BasecoatClient

    constructor() {
        this.client = new BasecoatClient(__API__, null, null);
    }

    // isUserLoggedIn determines if the user should be kicked back to the login route
    isUserLoggedIn(): boolean {
        if (!Cookies.get('username') || !Cookies.get('token')) {
            return false
        }
        return true
    }

    // handleLogin is a Promise that returns whether a login was successful
    // or not and sets relevant cookies
    handleLogin(loginInfo: LoginInfo): Promise<string> {
        let tokenRequest = new CreateAPITokenRequest();
        tokenRequest.setUser(loginInfo.username);
        tokenRequest.setPassword(loginInfo.password);
        tokenRequest.setDuration(10368000);

        return new Promise((resolve, reject) => {
            this.client.createAPIToken(tokenRequest, {}, function (err, response) {
                if (err) {
                    reject(err)
                    return
                }
                Cookies.set('username', loginInfo.username, { expires: 120, secure: true })
                Cookies.set('token', response.getKey(), { expires: 120, secure: true })
                resolve()
            })
        })
    }

    // handleLogout signs a user out by removing relevant cookies
    handleLogout() {
        Cookies.remove('username')
        Cookies.remove('token')
    }

    //getFormula retrieves a single formula by ID
    getFormula(formulaID: string): Promise<Formula | undefined> {
        let metadata = { 'Authorization': 'Bearer ' + Cookies.get('token') }
        let getFormulaRequest = new GetFormulaRequest();
        getFormulaRequest.setId(formulaID);

        return new Promise((resolve, reject) => {
            this.client.getFormula(getFormulaRequest, metadata, function (err, response) {
                if (err) {
                    console.log(err)
                    reject(undefined)
                    return
                }
                resolve(response.getFormula())
            });
        })
    }

    //getFormulaData retrieves all formulas from the backend
    getFormulaData(): Promise<formulaMap | undefined> {
        let metadata = { 'Authorization': 'Bearer ' + Cookies.get('token') }
        let listFormulasRequest = new ListFormulasRequest();

        return new Promise((resolve, reject) => {
            if (!this.isUserLoggedIn()) {
                reject(undefined)
            }
            this.client.listFormulas(listFormulasRequest, metadata, function (err, response) {
                if (err) {
                    console.log(err)
                    reject(undefined)
                }

                let formulas: formulaMap = {}
                response.getFormulasMap().forEach(function (value, key) {
                    formulas[key] = value
                })

                resolve(formulas)
            })
        })
    }

    //getJob retrieves a single job by ID
    getJob(jobID: string): Promise<Job | undefined> {
        let metadata = { 'Authorization': 'Bearer ' + Cookies.get('token') }
        let getJobRequest = new GetJobRequest();
        getJobRequest.setId(jobID);

        return new Promise((resolve, reject) => {
            this.client.getJob(getJobRequest, metadata, function (err, response) {
                if (err) {
                    console.log(err);
                    reject(undefined)
                    return
                }
                resolve(response.getJob())
            });
        })
    }

    //getJobData retrieves all jobs from the backend
    getJobData(): Promise<jobMap | undefined> {
        let metadata = { 'Authorization': 'Bearer ' + Cookies.get('token') }
        let listJobsRequest = new ListJobsRequest();

        return new Promise((resolve, reject) => {
            if (!this.isUserLoggedIn()) {
                reject(undefined)
            }
            this.client.listJobs(listJobsRequest, metadata, function (err, response) {
                if (err) {
                    console.log(err)
                    reject(undefined)
                }

                let jobs: jobMap = {}
                response.getJobsMap().forEach(function (value, key) {
                    jobs[key] = value
                })

                resolve(jobs)
            })
        })
    }

    //submitCreateFormulaForm submits the formula create form
    submitCreateFormulaForm(formulaData: CreateFormulaRequest.AsObject): Promise<string> {
        let metadata = { 'Authorization': 'Bearer ' + Cookies.get('token') }
        return new Promise((resolve, reject) => {

            if (!this.isUserLoggedIn()) {
                reject()
            }

            let createFormulaRequest = new CreateFormulaRequest();
            createFormulaRequest.setName(formulaData.name);
            createFormulaRequest.setNumber(formulaData.number);
            createFormulaRequest.setNotes(formulaData.notes);
            createFormulaRequest.setJobsList(formulaData.jobsList);

            let basesList: Base[] = []
            formulaData.basesList.forEach(function (item, index) {
                let newBase = new Base();
                newBase.setType(item.type)
                newBase.setName(item.name)
                newBase.setAmount(item.amount)

                basesList.push(newBase)
            });

            createFormulaRequest.setBasesList(basesList);

            let colorantsList: Colorant[] = []
            formulaData.colorantsList.forEach(function (item, index) {
                let newColorant = new Colorant();
                newColorant.setType(item.type)
                newColorant.setName(item.name)
                newColorant.setAmount(item.amount)

                colorantsList.push(newColorant)
            })
            createFormulaRequest.setColorantsList(colorantsList);

            this.client.createFormula(createFormulaRequest, metadata, function (err, response) {
                if (err) {
                    reject(err)
                }
                resolve()
            })
        })
    }

    //submitCreateJobForm submits a new job to the backend
    submitCreateJobForm(jobData: CreateJobRequest.AsObject): Promise<string> {
        let metadata = { 'Authorization': 'Bearer ' + Cookies.get('token') }
        return new Promise((resolve, reject) => {

            if (!this.isUserLoggedIn()) {
                reject()
            }

            let createJobRequest = new CreateJobRequest();
            createJobRequest.setName(jobData.name);
            createJobRequest.setStreet(jobData.street)
            createJobRequest.setStreet2(jobData.street2)
            createJobRequest.setCity(jobData.city)
            createJobRequest.setState(jobData.state)
            createJobRequest.setZipcode(jobData.zipcode)
            createJobRequest.setNotes(jobData.notes)
            createJobRequest.setFormulasList(jobData.formulasList)

            if (jobData.contact != undefined) {
                let contact = new Contact();
                contact.setName(jobData.contact.name)
                contact.setInfo(jobData.contact.info)
                createJobRequest.setContact(contact)
            }

            this.client.createJob(createJobRequest, metadata, function (err, response) {
                if (err) {
                    reject(err)
                }
                resolve()
            })
        })
    }

    submitManageFormulaForm(formulaData: UpdateFormulaRequest.AsObject): Promise<string> {
        let metadata = { 'Authorization': 'Bearer ' + Cookies.get('token') }
        return new Promise((resolve, reject) => {

            if (!this.isUserLoggedIn()) {
                reject()
            }

            let updateFormulaRequest = new UpdateFormulaRequest();
            updateFormulaRequest.setId(formulaData.id)
            updateFormulaRequest.setName(formulaData.name)
            updateFormulaRequest.setNumber(formulaData.number)
            updateFormulaRequest.setNotes(formulaData.notes)
            updateFormulaRequest.setJobsList(formulaData.jobsList)

            let basesList: Base[] = []
            formulaData.basesList.forEach(function (item, index) {
                let newBase = new Base();
                newBase.setType(item.type)
                newBase.setName(item.name)
                newBase.setAmount(item.amount)

                basesList.push(newBase)
            });

            updateFormulaRequest.setBasesList(basesList);

            let colorantsList: Colorant[] = []
            formulaData.colorantsList.forEach(function (item, index) {
                let newColorant = new Colorant();
                newColorant.setType(item.type)
                newColorant.setName(item.name)
                newColorant.setAmount(item.amount)

                colorantsList.push(newColorant)
            })
            updateFormulaRequest.setColorantsList(colorantsList);

            this.client.updateFormula(updateFormulaRequest, metadata, function (err, response) {
                if (err) {
                    reject(err)
                }
                resolve()
            })
        })
    }

    submitManageJobForm(jobData: UpdateJobRequest.AsObject): Promise<string> {
        let metadata = { 'Authorization': 'Bearer ' + Cookies.get('token') }
        return new Promise((resolve, reject) => {

            if (!this.isUserLoggedIn()) {
                reject()
            }

            let updateJobRequest = new UpdateJobRequest();
            updateJobRequest.setId(jobData.id)
            updateJobRequest.setName(jobData.name)
            updateJobRequest.setStreet(jobData.street)
            updateJobRequest.setStreet2(jobData.street2)
            updateJobRequest.setCity(jobData.city)
            updateJobRequest.setState(jobData.state)
            updateJobRequest.setZipcode(jobData.zipcode)
            updateJobRequest.setNotes(jobData.notes)
            updateJobRequest.setFormulasList(jobData.formulasList)

            if (jobData.contact != undefined) {
                let contact = new Contact();
                contact.setName(jobData.contact.name)
                contact.setInfo(jobData.contact.info)
                updateJobRequest.setContact(contact)
            }

            this.client.updateJob(updateJobRequest, metadata, function (err, response) {
                if (err) {
                    reject(err)
                }
                resolve()
            })
        })
    }

    deleteFormula(formulaID: string): Promise<string> {
        let metadata = { 'Authorization': 'Bearer ' + Cookies.get('token') }
        let deleteFormulaRequest = new DeleteFormulaRequest();
        deleteFormulaRequest.setId(formulaID)

        return new Promise((resolve, reject) => {

            if (!this.isUserLoggedIn()) {
                reject()
            }

            this.client.deleteFormula(deleteFormulaRequest, metadata, function (err, response) {
                if (err) {
                    reject(err)
                }
                resolve()
            })
        })
    }

    deleteJob(jobID: string): Promise<string> {
        let metadata = { 'Authorization': 'Bearer ' + Cookies.get('token') }
        let deleteJobRequest = new DeleteJobRequest();
        deleteJobRequest.setId(jobID)

        return new Promise((resolve, reject) => {

            if (!this.isUserLoggedIn()) {
                reject()
            }

            this.client.deleteJob(deleteJobRequest, metadata, function (err, response) {
                if (err) {
                    reject(err)
                }
                resolve()
            })
        })
    }
}

export default BasecoatClientWrapper
