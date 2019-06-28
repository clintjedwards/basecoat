import Vue from 'vue'
import Vuetify from 'vuetify'

import { VerifyLogin, LoadFormulaData, LoadJobData, SubmitCreateForm } from './methods'
import * as bcInterface from './basecoatInterfaces'
import store from './store'
import router from './router'

import PageHeader from "./components/PageHeader.vue"
import FormulasPage from "./components/FormulasPage.vue"
import JobsPage from "./components/JobsPage.vue"
import CreateFormulaModal from "./components/CreateFormulaModal.vue"
import AddJobModal from "./components/AddJobModal.vue"
import ManageFormulaModal from "./components/ManageFormulaModal.vue"
import ManageJobsModal from "./components/ManageJobsModal.vue"
import LoginModal from "./components/LoginModal.vue"

import { BasecoatClient } from "./BasecoatServiceClientPb"
import {
    CreateAPITokenRequest,
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

Vue.use(Vuetify)

// This is automatically set by the build process
declare var __API__: string;

// Create a basecoat client to communicate with grpc-web backend
let client = new BasecoatClient(__API__, null, null);

// We handle login out here because we need to handle it before
// vuerouter kicks in
VerifyLogin(client)

const app = new Vue({
    el: '#app',
    store,
    router,
    components: {
        PageHeader,
        FormulasPage,
        JobsPage,
        CreateFormulaModal,
        AddJobModal,
        ManageFormulaModal,
        ManageJobsModal,
        LoginModal
    },
    created: function () {
        setInterval(() => {
            if (this.$store.state.isLoggedIn) {
                LoadFormulaData
                LoadJobData
            }
        }, 180000); //3mins
    },
    methods: {
        navigateToFormulas() {
            this.$router.push('/formulas');
        },
        navigateToJobs() {
            this.$router.push('/jobs');
        },
        submitCreateForm: function (formulaData: CreateFormulaRequest.AsObject) {
            let self = this;
            SubmitCreateForm(client, formulaData);
            (self.$refs.createFormulaForm as HTMLFormElement).clearForm();
        },
        submitAddJobForm: function (jobData: CreateJobRequest.AsObject) {
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

            let metadata = { 'Authorization': 'Bearer ' + this.$cookies.get('token') }
            client.createJob(createJobRequest, metadata, function (err, response) {
                if (err) {
                    console.log(err)
                    store.commit('displaySnackBar', "Could not create job.")
                    return
                }
                store.commit("hideAddJobModal")
                self.loadFormulaData();
                self.loadJobData();
                (self.$refs.addJobForm as HTMLFormElement).clearForm();
            })
        },
        submitManageFormulaForm: function (formulaData: UpdateFormulaRequest.AsObject) {
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


            let metadata = { 'Authorization': 'Bearer ' + this.$cookies.get('token') }
            client.updateFormula(updateFormulaRequest, metadata, function (err, response) {
                if (err) {
                    console.log(err)
                    store.commit('displaySnackBar', "Could not update formula")
                    return
                }
                store.commit("hideManageFormulaModal")
                self.loadFormulaData();
                self.loadJobData();
                (self.$refs.manageFormulaForm as HTMLFormElement).setFormModeView();
            })
        },
        submitManageJobsForm: function (jobData: UpdateJobRequest.AsObject) {
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

            let metadata = { 'Authorization': 'Bearer ' + this.$cookies.get('token') }
            client.updateJob(updateJobRequest, metadata, function (err, response) {
                if (err) {
                    console.log(err)
                    store.commit('displaySnackBar', "Could not update job")
                    return
                }
                self.loadFormulaData();
                self.loadJobData();
                store.commit("hideManageJobsModal");
                (self.$refs.manageJobsForm as HTMLFormElement).setFormModeView();
            })
        },
        deleteJob: function (jobID: string) {
            let self = this
            let deleteJobRequest = new DeleteJobRequest();
            deleteJobRequest.setId(jobID)

            let metadata = { 'Authorization': 'Bearer ' + this.$cookies.get('token') }
            client.deleteJob(deleteJobRequest, metadata, function (err, response) {
                if (err) {
                    console.log(err)
                    store.commit('displaySnackBar', "Could not delete job")
                    return
                }
                store.commit("hideManageJobsModal")
                self.loadFormulaData();
                self.loadJobData();
                (self.$refs.manageJobsForm as HTMLFormElement).setFormModeView()
            })
        },
        deleteFormula: function (formulaID: string) {
            let self = this
            let deleteFormulaRequest = new DeleteFormulaRequest();
            deleteFormulaRequest.setId(formulaID)

            let metadata = { 'Authorization': 'Bearer ' + this.$cookies.get('token') }
            client.deleteFormula(deleteFormulaRequest, metadata, function (err, response) {
                if (err) {
                    console.log(err)
                    store.commit('displaySnackBar', "Could not delete formula")
                    return
                }
                store.commit("hideManageFormulaModal")
                self.loadFormulaData();
                self.loadJobData();
                (self.$refs.manageFormulaForm as HTMLFormElement).setFormModeView()
            })
        },
    }
})
