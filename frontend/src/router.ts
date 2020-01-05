import Vue from "vue";
import VueRouter from "vue-router";
import CreateContractorModal from "./components/CreateContractorModal/CreateContractorModal.vue";
import CreateFormulaModal from "./components/CreateFormulaModal.vue";
import CreateJobModal from "./components/CreateJobModal/CreateJobModal.vue";
import FormulasPage from "./components/FormulasPage.vue";
import JobsPage from "./components/JobsPage.vue";
import LoginModal from "./components/LoginModal.vue";
import ManageContractorModal from "./components/ManageContractorModal/ManageContractorModal.vue";
import ManageFormulaModal from "./components/ManageFormulaModal/ManageFormulaModal.vue";
import ManageJobModal from "./components/ManageJobModal/ManageJobModal.vue";
import NotFound from "./components/NotFound.vue";

Vue.use(VueRouter);

const routes = [
  { path: "/", redirect: "/formulas" },
  {
    path: "/formulas",
    name: "formulas",
    component: FormulasPage,
    children: [
      {
        path: "create",
        name: "createFormulaModal",
        component: CreateFormulaModal
      },
      {
        path: ":id",
        name: "manageFormulaModal",
        component: ManageFormulaModal
      }
    ]
  },
  {
    path: "/jobs",
    name: "jobs",
    component: JobsPage,
    children: [
      {
        path: "create",
        name: "jobCreateModal",
        component: CreateJobModal
      },
      {
        path: ":id",
        name: "manageJobModal",
        component: ManageJobModal
      },
      {
        path: "contractors/create",
        name: "contractorCreateModal",
        component: CreateContractorModal
      },
      {
        path: "contractors/:id",
        name: "manageContractorModal",
        component: ManageContractorModal
      }
    ]
  },
  {
    path: "/login",
    name: "login",
    component: LoginModal
  },
  { path: "*", component: NotFound }
];

const router = new VueRouter({
  routes,
  mode: "history"
});

export default router;
