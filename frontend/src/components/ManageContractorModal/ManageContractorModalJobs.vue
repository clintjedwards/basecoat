<template>
  <div>
    <v-spacer v-show="contractor.jobsList.length === 0">No Jobs Listed</v-spacer>
    <h2
      class="font-weight-light text-center"
      v-show="contractor.jobsList.length > 0"
      justify-space-around
    >
      Jobs
      <v-tooltip right>
        <template v-slot:activator="{ on }">
          <v-icon small color="text--secondary" style="vertical-align: middle" v-on="on">info</v-icon>
        </template>
        <span>Jobs linked to this contractor</span>
      </v-tooltip>
    </h2>

    <v-layout>
      <v-list two-line v-show="formMode === 'view'" style="width:100%;">
        <template v-for="(jobID, index) in contractor.jobsList">
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
          v-model="contractor.jobsList"
          :items="jobDataToList"
          item-text="name"
          item-value="id"
          hide-selected
          label="Link jobs to contractor"
          placeholder="Start typing to Search"
          multiple
          clearable
          counter
        >
          <template slot="item" slot-scope="data">
            <v-list-tile-content>
              <v-list-tile-title v-html="data.item.name"></v-list-tile-title>
              <v-list-tile-sub-title v-html="data.item.street"></v-list-tile-sub-title>
            </v-list-tile-content>
          </template>
        </v-autocomplete>
      </v-flex>
    </v-layout>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import { Job, Contractor } from "../../basecoat_message_pb";

interface modifiedJob {
  id: string;
  name: string;
  street: string;
}

export default Vue.extend({
  props: ["formMode", "contractor"],
  data: function() {
    return {};
  },
  computed: {
    jobDataToList(): modifiedJob[] {
      let jobDataMap: { [key: string]: Job } = this.$store.state.jobData;
      let contractorDataMap: { [key: string]: Contractor } = this.$store.state
        .contractorData;
      let jobDataList: modifiedJob[] = [];
      for (const [key, value] of Object.entries(jobDataMap)) {
        let job: modifiedJob;
        let street: string = "";
        let jobAddress = value.getAddress();
        if (jobAddress != undefined) {
          street = jobAddress.getStreet();
        }

        // if the job is either linked to a contractor or in our jobs list include it
        // this is because we need our own jobs to be included in our jobs list in order
        // to support unlinking them
        if (
          contractorDataMap[value.getContractorId()] !== undefined &&
          !this.contractor.jobsList.includes(key)
        ) {
          continue;
        }

        job = {
          id: value.getId(),
          name: value.getName(),
          street: street
        };
        jobDataList.push(job);
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
