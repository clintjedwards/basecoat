import Vue from 'vue'
import Vuex, { MutationTree } from 'vuex'
import * as bcInterface from './basecoatInterfaces'

Vue.use(Vuex)

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
        //app.$router.push('/jobs');
    },
    showManageFormulaModal(state) {
        state.displayManageFormulaModal = true;
    },
    hideManageFormulaModal(state) {
        state.displayManageFormulaModal = false;
        //app.$router.push('/formulas');
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

export default store
