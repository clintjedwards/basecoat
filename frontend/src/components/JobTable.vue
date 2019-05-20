
<template>
  <v-data-table
    :headers="headers"
    :items="jobDataToList"
    :search="$store.state.jobTableSearchTerm"
    hide-actions
  >
    <template v-slot:items="props">
      <tr
        style="cursor: pointer;"
        :ripple="{ center: true }"
        @click="$store.commit('showManageJobsModal', props.item.id)"
      >
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
  computed: {
    jobDataToList: function() {
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

      let jobDataMap: { [key: string]: Job } = this.$store.state.jobData;
      let jobDataList: Job[] = [];

      for (const [key, value] of Object.entries(jobDataMap)) {
        jobDataList.push(value);
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
  }
});
</script>

<style scoped>
</style>
