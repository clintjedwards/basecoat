<template>
  <v-container style="width: 80%;">
    <v-layout row justify-center>
      <v-flex md2 class="text-xs-center">
        <v-badge color="#ff5252">
          <template v-slot:badge>{{ Object.keys($store.state.jobData).length }}</template>
          <span class="display-3 font-weight-light search-panel-text">Jobs</span>
        </v-badge>
      </v-flex>
      <v-flex md8>
        <v-text-field
          label="Search"
          prepend-icon="search"
          hint="Search for job data of any kind"
          v-on:input="debounceInput($event)"
        ></v-text-field>
      </v-flex>
      <v-flex md2 style="margin-left: 1em;">
        <v-layout column>
          <v-flex>
            <v-btn color="secondary" v-on:click="showCreateJobModal()">
              <v-icon>create</v-icon>Create Job
            </v-btn>
          </v-flex>
          <v-flex>
            <v-btn color="secondary" v-on:click="showCreateContractorModal()">
              <v-icon>create</v-icon>Create Contractor
            </v-btn>
          </v-flex>
        </v-layout>
      </v-flex>
    </v-layout>
  </v-container>
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
    showCreateContractorModal: function() {
      this.$router.push({ name: "contractorCreateModal" });
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
