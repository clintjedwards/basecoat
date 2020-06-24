import Vue from "vue";
import Vuex, { MutationTree } from "vuex";
import { Contractor, Formula, Job } from "./basecoat_message_pb";
import { SnackBar, SnackBarColor } from "./snackbar";

Vue.use(Vuex);

interface FormulaMap {
  [key: string]: Formula;
}

interface JobMap {
  [key: string]: Job;
}

interface ContractorMap {
  [key: string]: Contractor;
}

interface colorantType {
  imageURL: string;
  userMessage: string;
}

interface colorantTypeMap {
  [key: string]: colorantType;
}

interface systemInfo {
  build_time: string;
  commit: string;
  debug_enabled: boolean;
  frontend_enabled: boolean;
  semver: string;
}

interface RootState {
  isInitialized: boolean;
  username: string;
  snackBar: SnackBar;
  appInfo: systemInfo;

  // Formula Data
  formulaData: FormulaMap;
  // filter lists contains the results from the search result
  // this allows us to check against this list and render only
  // results that match
  formulaDataFilter: string[];
  colorantTypes: colorantTypeMap;

  // Job Data
  jobData: JobMap;
  // filter lists contains the results from the search result
  // this allows us to check against this list and render only
  // results that match
  jobDataFilter: string[];

  // Contractor Data
  contractorData: ContractorMap;
}

const state: RootState = {
  // set so we can wait to load the store before accessing some components that depend on it
  isInitialized: false,
  username: "Unknown",
  snackBar: {
    text: "",
    display: false,
    color: SnackBarColor.Error,
  },
  appInfo: {
    build_time: "",
    commit: "",
    debug_enabled: true,
    frontend_enabled: false,
    semver: "",
  },

  // Formula Data
  formulaData: {},
  formulaDataFilter: [],
  colorantTypes: {
    "Benjamin Moore": {
      imageURL: "/images/benjamin-moore.png",
      userMessage: "Use Benjamin Moore Colorant Only",
    },
    "PPG Pittsburgh Paints": {
      imageURL: "/images/ppg.png",
      userMessage: "Use PPG Colorant Only",
    },
  },

  // Job Data
  jobData: {},
  jobDataFilter: [],

  // Contractor Data
  contractorData: {},
};

const mutations: MutationTree<RootState> = {
  setIsInitialized(state) {
    state.isInitialized = true;
  },
  updateUsername(state, username: string) {
    state.username = username;
  },
  updateAppInfo(state, systemInfo: systemInfo) {
    state.appInfo = systemInfo;
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
  updateContractorData(state, contractorData: ContractorMap) {
    state.contractorData = contractorData;
  },
  updateSnackBar(state, snackBar: SnackBar) {
    state.snackBar = snackBar;
  },
};

const store = new Vuex.Store<RootState>({
  state,
  mutations,
});

export default store;
