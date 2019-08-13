import Vue from 'vue'
import Vuetify from 'vuetify'
import VueCookies from 'vue-cookies'

import store from './store'
import router from './router'

import PageHeader from "./components/PageHeader.vue"
import FormulaSearchPanel from "./components/FormulaSearchPanel.vue"
import JobSearchPanel from "./components/JobSearchPanel.vue"
import FormulaTable from "./components/FormulaTable.vue"
import JobTable from "./components/JobTable.vue"
import CreateFormulaModal from "./components/CreateFormulaModal.vue"
import AddJobModal from "./components/AddJobModal.vue"
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

Vue.use(Vuetify)
Vue.use(VueCookies)

declare var __API__: string;

let client: BasecoatClient

const app = new Vue({
    el: '#app',
    store,
    router,
    components: {
        PageHeader,
        FormulaSearchPanel,
        JobSearchPanel,
        FormulaTable,
        JobTable,
        CreateFormulaModal,
        AddJobModal,
        LoginModal
    },
    created: function () {
        client = new BasecoatClient(__API__, null, null);
    },
    mounted() {
        this.checkLogin();

        setInterval(() => {
            if (this.$store.state.isLoggedIn) {
                this.loadFormulaData();
                this.loadJobData();
            }
        }, 180000); //3mins
    }
})
