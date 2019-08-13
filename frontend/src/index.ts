import Vue from 'vue'
import Vuetify from 'vuetify'

import Cookies from 'js-cookie'

import PageHeader from "./components/PageHeader.vue"

import BasecoatClientWrapper from './basecoatClientWrapper';

import store from './store'
import router from './router'

Vue.use(Vuetify)

let client: BasecoatClientWrapper
client = new BasecoatClientWrapper()

router.beforeEach((to, from, next) => {

    if (to.path === '/login') {
        next()
        return
    }
    if (!client.isUserLoggedIn()) {
        next({ name: "login", query: { redirect: to.path } })
        return
    }
    if (store.state.isInitialized) {
        next()
        return
    }

    var formulaPromise = client.getFormulaData()
    var jobsPromise = client.getJobData()

    Promise.all([formulaPromise, jobsPromise]).then((values) => {
        store.commit("updateFormulaData", values[0])
        store.commit("updateJobData", values[1])
        store.commit("setIsInitialized")
        next()
        return
    })
})

const app = new Vue({
    el: '#app',
    store,
    router,
    components: {
        PageHeader,
    },
    mounted() {
        setInterval(() => {
            client.getFormulaData().then((formulas) => {
                store.commit("updateFormulaData", formulas)
            })
            client.getJobData().then((jobs) => {
                store.commit("updateJobData", jobs)
            })
        }, 180000); //3mins

        store.commit("updateUsername", Cookies.get('username'));
    },
})
