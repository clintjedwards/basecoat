import Vue from 'vue'
import Vuex, { MutationTree } from 'vuex'
import Vuetify from 'vuetify'
import VueCookies from 'vue-cookies'
import VueRouter from 'vue-router'

import * as bcInterface from './basecoatInterfaces'

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

Vue.use(Vuex)
Vue.use(Vuetify)
Vue.use(VueCookies)
Vue.use(VueRouter)

// This is automatically set by the build process
declare var __API__: string;

const routes = [
    { path: '/', redirect: '/formulas' },
    { path: '/formulas', component: FormulasPage },
    { path: '/formulas/:id', component: ManageFormulaModal, props: true },
    { path: '/jobs', component: JobsPage },
    { path: '/jobs/:id', component: ManageJobsModal },
]

const router = new VueRouter({
    routes
})

interface RootState {
    formulaData: bcInterface.FormulaMap,
    jobData: bcInterface.JobMap,
    totalFormulas: number,
    totalJobs: number,
    formulaTableSearchTerm: string,
    jobTableSearchTerm: string,
    username: string,
    displayCreateFormulaModal: boolean,
    displayAddJobModal: boolean,
    displayManageFormulaModal: boolean,
    displayManageJobsModal: boolean,
    displayLoginModal: boolean,
    currentTab: string,
    isLoggedIn: boolean,
    loginIsLoading: boolean,
    displaySnackBar: boolean,
    snackBarText: string,
    colorantTypes: bcInterface.colorantTypeMap
}

const state: RootState = {
    formulaData: {},
    jobData: {},
    totalFormulas: 0,
    totalJobs: 0,
    formulaTableSearchTerm: "",
    jobTableSearchTerm: "",
    username: "Unknown",
    displayCreateFormulaModal: false,
    displayAddJobModal: false,
    displayManageFormulaModal: false,
    displayManageJobsModal: false,
    displayLoginModal: false,
    currentTab: "formulas",
    isLoggedIn: false,
    loginIsLoading: false,
    displaySnackBar: false,
    snackBarText: "",
    colorantTypes: {
        "Benjamin Moore": { imageURL: "/images/benjamin-moore.png", userMessage: "Use Benjamin Moore Colorant Only" },
        "PPG Pittsburgh Paints": { imageURL: "/images/ppg.png", userMessage: "Use PPG Colorant Only" }
    }
}

const mutations: MutationTree<RootState> = {
    showCreateFormulaModal(state) {
        state.displayCreateFormulaModal = true
    },
    hideCreateFormulaModal(state) {
        state.displayCreateFormulaModal = false
    },
    showAddJobModal(state) {
        state.displayAddJobModal = true
    },
    hideAddJobModal(state) {
        state.displayAddJobModal = false
        app.$router.push('/jobs');
    },
    showManageFormulaModal(state) {
        state.displayManageFormulaModal = true;
    },
    hideManageFormulaModal(state) {
        state.displayManageFormulaModal = false;
        app.$router.push('/formulas');
    },
    showManageJobsModal(state) {
        state.displayManageJobsModal = true;
    },
    hideManageJobsModal(state) {
        state.displayManageJobsModal = false
    },
    updateTotalFormulas(state) {
        state.totalFormulas = Object.keys(state.formulaData).length
    },
    updateTotalJobs(state) {
        state.totalJobs = Object.keys(state.jobData).length
    },
    updateUsername(state, username: string) {
        state.username = username
    },
    updateFormulaTableSearchTerm(state, searchTerm: string) {
        state.formulaTableSearchTerm = searchTerm
    },
    updateJobTableSearchTerm(state, searchTerm: string) {
        state.jobTableSearchTerm = searchTerm
    },
    updateFormulaData(state, formulaData: bcInterface.FormulaMap) {
        state.formulaData = formulaData
    },
    updateJobData(state, jobData: bcInterface.JobMap) {
        state.jobData = jobData
    },
    updateCurrentTab(state, tabName: string) {
        state.currentTab = tabName
    },
    displaySnackBar(state, text: string) {
        state.snackBarText = text
        state.displaySnackBar = true
    },
    updateLoginIsLoading(state, isLoading: boolean) {
        state.loginIsLoading = isLoading
    },
    updateLoginState(state, isLoggedIn: boolean) {
        state.isLoggedIn = isLoggedIn
        if (!isLoggedIn) {
            state.displayLoginModal = true
            return
        }

        state.displayLoginModal = false
    }
}

const store = new Vuex.Store<RootState>({
    state,
    mutations
})

let client: BasecoatClient

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
        client = new BasecoatClient(__API__, null, null);

        setInterval(() => {
            if (this.$store.state.isLoggedIn) {
                this.loadFormulaData();
                this.loadJobData();
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
        checkLogin: function () {
            if (!this.$cookies.isKey('username') || !this.$cookies.isKey('token')) {
                store.commit('updateLoginState', false)
                return
            }

            store.commit('updateUsername', this.$cookies.get('username'))
            store.commit('updateLoginState', true)
            this.loadFormulaData();
            this.loadJobData();
        },
        validateLogin: function (loginInfo: bcInterface.LoginInfo) {
            let self = this
            store.commit('updateLoginIsLoading', true)

            let tokenRequest = new CreateAPITokenRequest();
            tokenRequest.setUser(loginInfo.username);
            tokenRequest.setPassword(loginInfo.password);
            tokenRequest.setDuration(10368000);
            client.createAPIToken(tokenRequest, {}, function (err, response) {
                if (err) {
                    console.log(err)
                    store.commit('displaySnackBar', "Invalid Login Credentials")
                    store.commit('updateLoginState', false)
                    store.commit('updateLoginIsLoading', false)
                    return
                }
                self.$cookies.set('username', loginInfo.username, "4m", undefined, undefined, true)
                self.$cookies.set('token', response.getKey(), "4m", undefined, undefined, true)
                store.commit('updateUsername', loginInfo.username)
                store.commit('updateLoginState', true)
                store.commit('updateLoginIsLoading', false)
                self.loadFormulaData();
                self.loadJobData();
            })
        },
        handleLogout: function () {
            let self = this
            self.$cookies.remove('username')
            self.$cookies.remove('token')
            this.checkLogin()
        },
        loadFormulaData: function () {
            let listFormulasRequest = new ListFormulasRequest();
            let metadata = { 'Authorization': 'Bearer ' + this.$cookies.get('token') }
            client.listFormulas(listFormulasRequest, metadata, function (err, response) {
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
        },
        loadJobData: function () {
            let listJobsRequest = new ListJobsRequest();
            let metadata = { 'Authorization': 'Bearer ' + this.$cookies.get('token') }
            client.listJobs(listJobsRequest, metadata, function (err, response) {
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
        },
        submitCreateForm: function (formulaData: CreateFormulaRequest.AsObject) {
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

            let metadata = { 'Authorization': 'Bearer ' + this.$cookies.get('token') }
            client.createFormula(createFormulaRequest, metadata, function (err, response) {
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
                (self.$refs.createFormulaForm as HTMLFormElement).clearForm();
            })
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
