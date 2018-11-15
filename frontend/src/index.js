import Vue from 'vue'
import Vuex from 'vuex'
import Vuetify from 'vuetify'
import VueCookies from 'vue-cookies'
import axios from 'axios'

import PageHeader from "./components/PageHeader.vue"
import FormulaSearchPanel from "./components/FormulaSearchPanel.vue"
import JobsSearchPanel from "./components/JobsSearchPanel.vue"
import FormulaTable from "./components/FormulaTable.vue"
import JobTable from "./components/JobTable.vue"
import CreateFormulaModal from "./components/CreateFormulaModal.vue"
import AddJobModal from "./components/AddJobModal.vue"
import ManageFormulaModal from "./components/ManageFormulaModal.vue"
import ManageJobsModal from "./components/ManageJobsModal.vue"
import LoginModal from "./components/LoginModal.vue"

Vue.use(Vuex)
Vue.use(Vuetify)
Vue.use(VueCookies)


const store = new Vuex.Store({
    state: {
        formulaData: [],
        jobData: [],
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
        formulaInView: {},
        jobInView: {},
        isLoggedIn: false,
        loginIsLoading: false,
        displaySnackBar: false,
        snackBarText: ""
    },
    mutations: {
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
        },
        showManageFormulaModal(state, formulaID) {
            state.displayManageFormulaModal = true
            state.formulaInView = state.formulaData[formulaID]
        },
        hideManageFormulaModal(state) {
            state.displayManageFormulaModal = false
        },
        showManageJobsModal(state, jobID) {
            state.displayManageJobsModal = true
            state.jobInView = state.jobData[jobID]
        },
        hideManageJobsModal(state) {
            state.displayManageJobsModal = false
        },
        updateTotalFormulas(state) {
            state.totalFormulas = state.formulaData.length
        },
        updateTotalJobs(state) {
            state.totalJobs = state.jobData.length
        },
        updateUsername(state, username) {
            state.username = username
        },
        updateFormulaTableSearchTerm(state, searchTerm) {
            state.formulaTableSearchTerm = searchTerm
        },
        updateJobTableSearchTerm(state, searchTerm) {
            state.jobTableSearchTerm = searchTerm
        },
        updateFormulaData(state, formulaData) {
            state.formulaData = formulaData
        },
        updateJobData(state, jobData) {
            state.jobData = jobData
        },
        updateCurrentTab(state, tabName) {
            state.currentTab = tabName
        },
        displaySnackBar(state, text) {
            state.snackBarText = text
            state.displaySnackBar = true
        },
        updateLoginIsLoading(state, isLoading) {
            state.loginIsLoading = isLoading
        },
        updateLoginState(state, isLoggedIn) {
            state.isLoggedIn = isLoggedIn
            if (!isLoggedIn) {
                state.displayLoginModal = true
                return
            }

            state.displayLoginModal = false
        }
    },
})

const app = new Vue({
    el: '#app',
    store,
    components: {
        PageHeader,
        FormulaSearchPanel,
        JobsSearchPanel,
        FormulaTable,
        JobTable,
        CreateFormulaModal,
        AddJobModal,
        ManageFormulaModal,
        ManageJobsModal,
        LoginModal
    },
    methods: {
        checkLogin: function () {
            if (!this.$cookies.isKey('username') || !this.$cookies.isKey('token')) {
                store.commit('updateLoginState', false)
                return
            }

            store.commit('updateUsername', this.$cookies.get('username'))
            store.commit('updateLoginState', true)
        },
        validateLogin: function (loginInfo) {
            store.commit('updateLoginIsLoading', true)
            const headers = {
                'Content-Type': 'application/json',
                'Authorization': loginInfo.token,
            };
            axios
                .get('/auth/api', {
                    headers
                })
                .then(function (response) {
                    if (response.status == 200) {
                        $cookies.set('username', loginInfo.username, "4m", null, null, null, true)
                        $cookies.set('token', loginInfo.token, "4m", null, null, null, true)
                        store.commit('updateUsername', loginInfo.username)
                        store.commit('updateLoginState', true)
                    } else {
                        console.log("loginFailed")
                        store.commit('updateLoginState', false)
                    }
                })
                .catch((error) => {
                    console.log(error)
                    store.commit('displaySnackBar', "Invalid Login Credentials")
                    store.commit('updateLoginState', false)
                })
                .finally(function () {
                    store.commit('updateLoginIsLoading', false)
                })
        },
        loadFormulaData: function () {
            const headers = {
                'Content-Type': 'application/json',
                'Authorization': this.$cookies.get('token'),
            };
            axios
                .get('/formulas', {
                    headers
                })
                .then(function (response) {
                    store.commit('updateFormulaData', response.data)
                })
                .catch(error => console.log(error))
                .finally(function () {
                    store.commit('updateTotalFormulas')
                })
        },
        loadJobData: function () {
            axios
                .get('/jobs')
                .then(response => (store.commit('updateJobData', response.data)))
                .catch(error => console.log(error))
                .finally(function () {
                    store.commit('updateTotalJobs')
                })
        },
        submitCreateForm: function (formulaData) {
            var self = this

            const headers = {
                'Content-Type': 'application/json',
                'Authorization': this.$cookies.get('token'),
            };

            axios
                .post('/formulas', formulaData, {
                    headers
                })
                .then(function (response) {
                    store.commit("hideCreateFormulaModal")
                })
                .catch(function (error) {
                    console.log(error)
                })
                .then(function () {
                    self.loadFormulaData()
                })
                .then(function () {
                    self.$refs.createFormulaForm.clearForm()
                })
        },
        submitAddJobForm: function (jobData) {
            var self = this

            const headers = {
                'Content-Type': 'application/json',
                'Authorization': this.$cookies.get('token'),
            };

            axios
                .post('/jobs', jobData, {
                    headers
                })
                .then(function (response) {
                    store.commit("hideAddJobModal")
                })
                .catch(function (error) {
                    console.log(error)
                })
                .then(function () {
                    self.loadJobData()
                })
                .then(function () {
                    self.$refs.addJobForm.clearForm()
                })
        },
        submitManageFormulaForm: function (formulaData) {
            var self = this

            const headers = {
                'Content-Type': 'application/json',
                'Authorization': this.$cookies.get('token'),
            };

            axios
                .put('/formulas/' + formulaData.id, formulaData, {
                    headers
                })
                .then(function (response) {
                    store.commit("hideManageFormulaModal")
                })
                .catch(function (error) {
                    console.log(error)
                })
                .then(function () {
                    self.loadFormulaData()
                })
                .then(function () {
                    self.$refs.manageFormulaForm.setFormModeView()
                })
        },
        submitManageJobsForm: function (jobData) {
            var self = this

            const headers = {
                'Content-Type': 'application/json',
                'Authorization': this.$cookies.get('token'),
            };

            axios
                .put('/jobs/' + jobData.id, jobData, {
                    headers
                })
                .then(function (response) {
                    store.commit("hideManageJobsModal")
                })
                .catch(function (error) {
                    console.log(error)
                })
                .then(function () {
                    self.loadJobData()
                })
                .then(function () {
                    self.$refs.manageJobsForm.setFormModeView()
                })
        }
    },
    mounted() {
        this.checkLogin();

        this.loadFormulaData();
        this.loadJobData();
        setInterval(() => {
            this.loadFormulaData();
            this.loadJobData();
        }, 180000); //3mins
    }
})
