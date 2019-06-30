import Vue from 'vue'
import VueRouter from 'vue-router'

import FormulasPage from "./components/FormulasPage.vue"
import JobsPage from "./components/JobsPage.vue"
import ManageFormulaModal from "./components/ManageFormulaModal.vue"
import ManageJobModal from "./components/ManageJobModal.vue"


Vue.use(VueRouter)

const routes = [
    { path: '/', redirect: '/formulas' },
    {
        path: '/formulas',
        name: 'formulas',
        component: FormulasPage,
        children: [{ path: ':id', name: 'formulaModal', component: ManageFormulaModal }]
    },
    {
        path: '/jobs',
        name: 'jobs',
        component: JobsPage,
        children: [{ path: ':id', name: 'jobModal', component: ManageJobModal }]
    },
]

const router = new VueRouter({
    routes
})

export default router
