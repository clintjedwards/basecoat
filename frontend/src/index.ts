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

import BasecoatClientWrapper from './methods';

Vue.use(Vuetify)
Vue.use(VueCookies)

let client: BasecoatClientWrapper

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
        client = new BasecoatClientWrapper()
    },
    mounted() {
        client.checkLogin()

        setInterval(() => {
            if (this.$store.state.isLoggedIn) {
                client.loadFormulaData();
                client.loadJobData();
            }
        }, 180000); //3mins
    }
})
