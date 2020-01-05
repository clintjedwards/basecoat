import Cookies from "js-cookie";
import Vue from "vue";
import Vuetify from "vuetify";
import BasecoatClientWrapper from "./basecoatClientWrapper";
import PageFooter from "./components/PageFooter.vue";
import PageHeader from "./components/PageHeader.vue";
import router from "./router";
import store from "./store";

Vue.use(Vuetify);

let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

// On the first page load of the app we load all data.
// To make sure that the data is ready before we try to load the page we
// use a 'setIsInitialized' function that tells any page of the application that
// data is ready to be used. We structure all of this in a promise to make sure
// concurrency doesn't fuck us.
router.beforeEach((to, from, next) => {
  if (to.path === "/login") {
    next();
    return;
  }
  // if user is not logged in redirect them to login
  if (!client.isUserLoggedIn()) {
    next({ name: "login", query: { redirect: to.path } });
    return;
  }
  // if store is initialized go straight to next route, if not load the app
  if (store.state.isInitialized) {
    next();
    return;
  }

  var formulaPromise = client.getFormulaData();
  var jobsPromise = client.getJobData();
  var contractorsPromise = client.getContractorData();

  Promise.all([formulaPromise, jobsPromise, contractorsPromise]).then(
    values => {
      store.commit("updateFormulaData", values[0]);
      store.commit("updateJobData", values[1]);
      store.commit("updateContractorData", values[2]);
      store.commit("setIsInitialized");
      next();
      return;
    }
  );
});

const app = new Vue({
  el: "#app",
  store,
  router,
  components: {
    PageFooter,
    PageHeader
  },
  mounted() {
    client.getSystemInfo().then(systemInfo => {
      if (systemInfo) {
        store.commit("updateAppInfo", systemInfo);
      }
    });

    setInterval(() => {
      client.getFormulaData().then(formulas => {
        store.commit("updateFormulaData", formulas);
      });
      client.getJobData().then(jobs => {
        store.commit("updateJobData", jobs);
      });
      client.getContractorData().then(contractors => {
        store.commit("updateContractorData", contractors);
      });
    }, 180000); //3mins

    store.commit("updateUsername", Cookies.get("username"));
  }
});
