
<template>
  <v-data-table :headers="headers" :items="jobDataToList" hide-actions disable-filtering>
    <template v-slot:items="props">
      <tr style="cursor: pointer;" :ripple="{ center: true }" @click="navigateToJob(props.item.id)">
        <td class="text-capitalize">{{ props.item.name }}</td>
        <td>
          <span class="text-capitalize">{{ props.item.contact_name }}</span>
          <template v-if="props.item.contact_name !== ''">|</template>
          {{ props.item.contact_info }}
        </td>
        <td>
          {{ props.item.street }} {{ props.item.city }}
          <template v-if="props.item.city !== ''">,</template>
          {{ props.item.state }} {{ props.item.zipcode }}
        </td>
      </tr>
    </template>
  </v-data-table>
</template>

<script lang="ts">
import Vue from "vue";
import { Job, Contact } from "../basecoat_pb";

import BasecoatClientWrapper from "../basecoatClientWrapper";

let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

interface modifiedJob {
  id: string;
  name: string;
  contact_name: string;
  contact_info: string;
  street: string;
  city: string;
  state: string;
  zipcode: string;
}

export default Vue.extend({
  data: function() {
    return {
      headers: [
        {
          text: "Name",
          align: "left",
          value: "name"
        },
        {
          text: "Contact Info",
          value: "contact_info"
        },
        {
          text: "Address",
          value: "address"
        }
      ]
    };
  },
  mounted() {
    this.loadJobData();
  },
  computed: {
    jobDataToList: function() {
      let filteredIDs: string[] = this.$store.state.jobDataFilter;

      let jobDataMap: { [key: string]: Job } = this.$store.state.jobData;
      let jobDataList: Job[] = [];

      for (const [key, value] of Object.entries(jobDataMap)) {
        if (filteredIDs && filteredIDs.length) {
          if (filteredIDs.includes(key)) {
            jobDataList.push(value);
          }
        } else {
          jobDataList.push(value);
        }
      }

      let modifiedJobList: modifiedJob[] = [];
      let job: Job;

      for (job of jobDataList) {
        let modifiedJob: modifiedJob = {
          id: "",
          name: "",
          contact_name: "",
          contact_info: "",
          street: "",
          city: "",
          state: "",
          zipcode: ""
        };

        modifiedJob.id = job.getId();
        modifiedJob.name = job.getName();
        let jobContact = job.getContact();
        if (jobContact != undefined) {
          modifiedJob.contact_name = jobContact.getName();
          modifiedJob.contact_info = jobContact.getInfo();
        }
        modifiedJob.street = job.getStreet();
        modifiedJob.city = job.getCity();
        modifiedJob.state = job.getState();
        modifiedJob.zipcode = job.getZipcode();

        modifiedJobList.push(modifiedJob);
      }

      return modifiedJobList;
    }
  },
  methods: {
    navigateToJob: function(jobID: string) {
      this.$router.push("/jobs/" + jobID);
    },
    loadJobData: function() {
      client.getJobData().then(jobs => {
        this.$store.commit("updateJobData", jobs);
      });
    }
  }
});
</script>

<style scoped>
</style>
