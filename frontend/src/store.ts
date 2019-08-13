import Vue from 'vue'
import Vuex, { MutationTree } from 'vuex'

import { Formula, Job } from "./basecoat_pb"

Vue.use(Vuex)

interface FormulaMap {
    [key: string]: Formula;
}

interface JobMap {
    [key: string]: Job;
}

interface colorantType {
    imageURL: string
    userMessage: string
}

interface colorantTypeMap {
    [key: string]: colorantType;
}

interface RootState {
    isInitialized: boolean,
    username: string,
    snackBarText: string,
    displaySnackBar: boolean,

    // Formula Data
    formulaData: FormulaMap,
    formulaTableSearchTerm: string,
    colorantTypes: colorantTypeMap,

    // Job Data
    jobData: JobMap,
    jobTableSearchTerm: string,
}

const state: RootState = {
    // set so we can wait to load the store before accessing some components that depend on it
    isInitialized: false,
    username: "Unknown",
    snackBarText: "",
    displaySnackBar: false,

    // Formula Data
    formulaData: {},
    formulaTableSearchTerm: "",
    colorantTypes: {
        "Benjamin Moore": { imageURL: "/images/benjamin-moore.png", userMessage: "Use Benjamin Moore Colorant Only" },
        "PPG Pittsburgh Paints": { imageURL: "/images/ppg.png", userMessage: "Use PPG Colorant Only" }
    },

    // Job Data
    jobData: {},
    jobTableSearchTerm: "",
}

const mutations: MutationTree<RootState> = {
    setIsInitialized(state) {
        state.isInitialized = true
    },
    updateUsername(state, username: string) {
        state.username = username
    },
    //Formula Data Mutators
    updateFormulaTableSearchTerm(state, searchTerm: string) {
        state.formulaTableSearchTerm = searchTerm
    },
    updateFormulaData(state, formulaData: FormulaMap) {
        state.formulaData = formulaData
    },
    //Job Data Mutators
    updateJobTableSearchTerm(state, searchTerm: string) {
        state.jobTableSearchTerm = searchTerm
    },
    updateJobData(state, jobData: JobMap) {
        state.jobData = jobData
    },
    showSnackBar(state, text: string) {
        state.snackBarText = text
        state.displaySnackBar = true
    },
}

const store = new Vuex.Store<RootState>({
    state,
    mutations
})

export default store
