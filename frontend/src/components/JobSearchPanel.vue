<template>
  <div>
    <span>
      <v-badge color="#ff5252">
        <template v-slot:badge>{{ Object.keys($store.state.jobData).length }}</template>
        <span class="display-3 font-weight-light search-panel-text">Jobs</span>
      </v-badge>
    </span>
    <span style="margin-left: 4em; width: 100em;">
      <v-text-field
        label="Search"
        prepend-icon="search"
        hint="Search for job data of any kind"
        v-on:input="debounceInput($event)"
      ></v-text-field>
    </span>
    <span style="margin-left: 4em;">
      <v-btn color="secondary" v-on:click="showCreateJobModal()">
        <v-icon>create</v-icon>Create Job
      </v-btn>
    </span>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import _ from "lodash";

import BasecoatClientWrapper from "../basecoatClientWrapper";

let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

export default Vue.extend({
  data: function() {
    return {};
  },
  methods: {
    showCreateJobModal: function() {
      this.$router.push({ name: "jobCreateModal" });
    },
    debounceInput: _.debounce(function(this: any, searchTerm) {
      if (!searchTerm) {
        this.$store.commit("updateJobDataFilter", []);
        return;
      }
      client.searchJobs(searchTerm).then(hits => {
        if (hits != undefined) {
          this.$store.commit("updateJobDataFilter", hits);
          return;
        }
        this.$store.commit("updateJobDataFilter", []);
        return;
      });
    }, 625)
  }
});
</script>

<style scoped>
.search-panel-text {
  color: #9e9e9e;
}
</style>
