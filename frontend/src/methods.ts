import {
    CreateAPITokenRequest,
    ListFormulasRequest,
    Formula, Job,
    ListJobsRequest,
    CreateFormulaRequest,
    CreateJobRequest,
    UpdateFormulaRequest,
    UpdateJobRequest,
    DeleteFormulaRequest,
    DeleteJobRequest,
    Base,
    Colorant,
    Contact
} from "./basecoat_pb"

import store from './store'
import Cookies from 'js-cookie'
import { BasecoatClient } from "./BasecoatServiceClientPb"

interface LoginInfo {
    username: string;
    password: string;
}

declare var __API__: string;

class BasecoatClientWrapper {

    client: BasecoatClient
    metadata: { 'Authorization': string }

    constructor() {
        this.client = new BasecoatClient(__API__, null, null);
        this.metadata = { 'Authorization': 'Bearer ' + Cookies.get('token') }
    }

    checkLogin() {
        if (!Cookies.get('username') || !Cookies.get('token')) {
            store.commit('updateLoginState', false)
            return
        }

        store.commit('updateUsername', Cookies.get('username'))
        store.commit('updateLoginState', true)
        this.loadFormulaData();
        this.loadJobData();
    }

    handleLogout() {
        Cookies.remove('username')
        Cookies.remove('token')
        this.checkLogin()
    }

    validateLogin(loginInfo: LoginInfo) {
        let self = this
        store.commit('updateLoginIsLoading', true)

        let tokenRequest = new CreateAPITokenRequest();
        tokenRequest.setUser(loginInfo.username);
        tokenRequest.setPassword(loginInfo.password);
        tokenRequest.setDuration(10368000);
        this.client.createAPIToken(tokenRequest, {}, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Invalid Login Credentials")
                store.commit('updateLoginState', false)
                store.commit('updateLoginIsLoading', false)
                return
            }
            Cookies.set('username', loginInfo.username, { expires: 120, secure: true })
            Cookies.set('token', response.getKey(), { expires: 120, secure: true })
            store.commit('updateUsername', loginInfo.username)
            store.commit('updateLoginState', true)
            store.commit('updateLoginIsLoading', false)
            self.loadFormulaData();
            self.loadJobData();
        })
    }

    loadFormulaData() {
        let listFormulasRequest = new ListFormulasRequest();
        this.client.listFormulas(listFormulasRequest, this.metadata, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Could not load current formulas.")
                return
            }

            let formulaMap: { [key: string]: Formula } = {}
            response.getFormulasMap().forEach(function (value, key) {
                formulaMap[key] = value
            })
            store.commit('updateFormulaData', formulaMap)
            store.commit('updateTotalFormulas')
        })
    }

    loadJobData() {
        let client = new BasecoatClient(__API__, null, null);
        let listJobsRequest = new ListJobsRequest();
        client.listJobs(listJobsRequest, this.metadata, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Could not load current jobs.")
                return
            }

            let jobMap: { [key: string]: Job } = {}
            response.getJobsMap().forEach(function (value, key) {
                jobMap[key] = value
            })
            store.commit('updateJobData', jobMap)
            store.commit('updateTotalJobs')
        })
    }

    submitCreateForm(formulaData: CreateFormulaRequest.AsObject) {
        let self = this
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

        this.client.createFormula(createFormulaRequest, this.metadata, function (err, response) {
            if (err) {
                console.log(err)
                // If formula already exists return helpful error
                if (err.code == 6) {
                    store.commit('displaySnackBar', "Could not create formula. Please make sure formula name is unique.")
                    return
                }
                store.commit('displaySnackBar', "Could not create formula.")
                return
            }
            store.commit("hideCreateFormulaModal")
            self.loadFormulaData();
            self.loadJobData();
            //(self.$refs.createFormulaForm as HTMLFormElement).clearForm();
        })
    }

    submitAddJobForm(jobData: CreateJobRequest.AsObject) {
        let self = this
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

        this.client.createJob(createJobRequest, this.metadata, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Could not create job.")
                return
            }
            store.commit("hideAddJobModal")
            self.loadFormulaData();
            self.loadJobData();
            //(self.$refs.addJobForm as HTMLFormElement).clearForm();
        })
    }

    submitManageFormulaForm(formulaData: UpdateFormulaRequest.AsObject) {
        let self = this
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


        this.client.updateFormula(updateFormulaRequest, this.metadata, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Could not update formula")
                return
            }
            store.commit("hideManageFormulaModal")
            self.loadFormulaData();
            self.loadJobData();
            //(self.$refs.manageFormulaForm as HTMLFormElement).setFormModeView();
        })
    }

    submitManageJobForm(jobData: UpdateJobRequest.AsObject) {
        let self = this
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

        this.client.updateJob(updateJobRequest, this.metadata, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Could not update job")
                return
            }
            self.loadFormulaData();
            self.loadJobData();
            store.commit("hideManageJobModal");
            //(self.$refs.manageJobForm as HTMLFormElement).setFormModeView();
        })
    }

    deleteJob(jobID: string) {
        let self = this
        let deleteJobRequest = new DeleteJobRequest();
        deleteJobRequest.setId(jobID)

        this.client.deleteJob(deleteJobRequest, this.metadata, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Could not delete job")
                return
            }
            store.commit("hideManageJobModal")
            self.loadFormulaData();
            self.loadJobData();
            //(self.$refs.manageJobForm as HTMLFormElement).setFormModeView()
        })
    }

    deleteFormula(formulaID: string) {
        let self = this
        let deleteFormulaRequest = new DeleteFormulaRequest();
        deleteFormulaRequest.setId(formulaID)

        this.client.deleteFormula(deleteFormulaRequest, this.metadata, function (err, response) {
            if (err) {
                console.log(err)
                store.commit('displaySnackBar', "Could not delete formula")
                return
            }
            store.commit("hideManageFormulaModal")
            self.loadFormulaData();
            self.loadJobData();
            //(self.$refs.manageFormulaForm as HTMLFormElement).setFormModeView()
        })
    }
}

export default BasecoatClientWrapper
