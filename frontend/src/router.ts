import Vue from 'vue'
import VueRouter from 'vue-router'

import FormulasPage from "./components/FormulasPage.vue"
import JobsPage from "./components/JobsPage.vue"
import ManageFormulaModal from "./components/ManageFormulaModal.vue"
import ManageJobsModal from "./components/ManageJobsModal.vue"

Vue.use(VueRouter)

const routes = [
    { path: '/', redirect: '/formulas' },
    { path: '/formulas', component: FormulasPage },
    { path: '/formulas/:id', component: ManageFormulaModal, props: true },
    { path: '/jobs', component: JobsPage },
    { path: '/jobs/:id', component: ManageJobsModal, props: true },
]

const router = new VueRouter({
    routes
})

export default router

