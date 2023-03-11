<template>
  <div>
    <v-spacer v-show="formula.jobsList.length === 0">No Jobs Listed</v-spacer>
    <h2
      class="font-weight-light text-center"
      v-show="formula.jobsList.length > 0"
      justify-space-around
    >
      Jobs
      <v-tooltip right>
        <template v-slot:activator="{ on }">
          <v-icon small color="text--secondary" style="vertical-align: middle" v-on="on">info</v-icon>
        </template>
        <span>Job sites that have requested this formula</span>
      </v-tooltip>
    </h2>

    <v-layout>
      <v-list two-line v-show="formMode === 'view'" style="width:100%;">
        <template v-for="(jobID, index) in formula.jobsList">
          <v-list-tile v-bind:key="`job-tile-${index}`">
            <v-list-tile-avatar>
              <v-icon>work</v-icon>
            </v-list-tile-avatar>
            <v-list-tile-content>
              <v-list-tile-title v-text="getJobObject(jobID).name"></v-list-tile-title>
              <v-list-tile-sub-title>{{ getAddress(jobID) }}</v-list-tile-sub-title>
            </v-list-tile-content>
          </v-list-tile>
          <v-divider v-bind:key="`job-divider-${index}`"></v-divider>
        </template>
      </v-list>

      <v-flex xs12 v-show="formMode === 'edit'">
        <v-autocomplete
          v-model="formula.jobsList"
          :items="jobDataToList"
          item-text="name"
          item-value="id"
          hide-selected
          label="Link jobs to formula"
          placeholder="Start typing to Search"
          multiple
          clearable
          counter
        >
          <template slot="item" slot-scope="data">
            <v-list-tile-content>
              <v-list-tile-title v-html="data.item.name"></v-list-tile-title>
              <v-list-tile-sub-title v-html="data.item.address.street"></v-list-tile-sub-title>
            </v-list-tile-content>
          </template>
        </v-autocomplete>
      </v-flex>
    </v-layout>
    <br />
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import { Job, Formula } from "../../basecoat_message_pb";

export default Vue.extend({
  props: ["formMode", "formula"],
  data: function() {
    return {};
  },
  computed: {
    jobDataToList: function() {
      let jobDataMap: { [key: string]: Job } = this.$store.state.jobData;
      let jobDataList: Job.AsObject[] = [];

      for (const [key, value] of Object.entries(jobDataMap)) {
        jobDataList.push(value.toObject());
      }

      return jobDataList;
    }
  },
  methods: {
    getJobObject: function(jobID: string) {
      let jobMap: { [key: string]: Job } = this.$store.state.jobData;
      return jobMap[jobID].toObject();
    },
    getAddress: function(jobID: string) {
      let job: Job.AsObject = this.getJobObject(jobID);
      let strBuilder: string[] = [];

      if (job.address === undefined) {
        return "";
      }

      strBuilder.push(job.address.street);
      strBuilder.push(" ");
      strBuilder.push(job.address.city);
      if (job.address.city !== "") {
        strBuilder.push(", ");
      } else {
        strBuilder.push(" ");
      }
      strBuilder.push(job.address.state);
      strBuilder.push(" ");
      strBuilder.push(job.address.zipcode);

      return strBuilder.join("");
    }
  }
});
</script>
<style scoped>
</style>
