import Vue from "vue";
import VueRouter from "vue-router";
import CreateFormulaModal from "./components/CreateFormulaModal.vue";
import CreateJobModal from "./components/CreateJobModal.vue";
import FormulasPage from "./components/FormulasPage.vue";
import JobsPage from "./components/JobsPage.vue";
import LoginModal from "./components/LoginModal.vue";
import ManageFormulaModal from "./components/ManageFormulaModal.vue";
import ManageJobModal from "./components/ManageJobModal.vue";
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
