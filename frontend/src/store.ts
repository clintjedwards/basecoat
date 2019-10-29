import Vue from "vue";
import Vuex, { MutationTree } from "vuex";
import { Formula, Job } from "./basecoat_pb";

Vue.use(Vuex);

interface FormulaMap {
  [key: string]: Formula;
}

interface JobMap {
  [key: string]: Job;
}

interface colorantType {
  imageURL: string;
  userMessage: string;
}

interface colorantTypeMap {
  [key: string]: colorantType;
}

interface RootState {
  isInitialized: boolean;
  username: string;
  snackBarText: string;
  displaySnackBar: boolean;

  // Formula Data
  formulaData: FormulaMap;
  formulaDataFilter: string[];
  colorantTypes: colorantTypeMap;

  // Job Data
  jobData: JobMap;
  jobDataFilter: string[];
}

const state: RootState = {
  // set so we can wait to load the store before accessing some components that depend on it
  isInitialized: false,
  username: "Unknown",
  snackBarText: "",
  displaySnackBar: false,

  // Formula Data
  formulaData: {},
  formulaDataFilter: [],
  colorantTypes: {
    "Benjamin Moore": {
      imageURL: "/images/benjamin-moore.png",
      userMessage: "Use Benjamin Moore Colorant Only"
    },
    "PPG Pittsburgh Paints": {
      imageURL: "/images/ppg.png",
      userMessage: "Use PPG Colorant Only"
    }
  },

  // Job Data
  jobData: {},
  jobDataFilter: []
};

const mutations: MutationTree<RootState> = {
  setIsInitialized(state) {
    state.isInitialized = true;
  },
  updateUsername(state, username: string) {
    state.username = username;
  },
  //Formula Data Mutators
  updateFormulaDataFilter(state, newFilterList: string[]) {
    state.formulaDataFilter = newFilterList;
  },
  updateFormulaData(state, formulaData: FormulaMap) {
    state.formulaData = formulaData;
  },
  //Job Data Mutators
  updateJobDataFilter(state, newFilterList: string[]) {
    state.jobDataFilter = newFilterList;
  },
  updateJobData(state, jobData: JobMap) {
    state.jobData = jobData;
  },
  showSnackBar(state, text: string) {
    state.snackBarText = text;
    state.displaySnackBar = true;
  }
};

const store = new Vuex.Store<RootState>({
  state,
  mutations
});

export default store;
