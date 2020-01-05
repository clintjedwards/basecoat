
<template>
  <div class="sizer data-table">
    <v-container>
      <v-switch class="pull-right" v-model="showEmptyContractors" label="Show Empty Contractors"></v-switch>
      <v-layout row>
        <v-flex>
          <v-card elevation="5" v-for="contractor in contractorList" :key="contractor.id">
            <div v-show="filterJobList(contractor.jobsList).length > 0 || showEmptyContractors">
              <v-list v-show="contractor.id != ''">
                <v-list-tile avatar @click="navigateToContractor(contractor.id)">
                  <v-list-tile-avatar>
                    <v-icon>business</v-icon>
                  </v-list-tile-avatar>

                  <v-list-tile-content>
                    <v-list-tile-title class="headline font-weight-light">{{ contractor.company }}</v-list-tile-title>
                  </v-list-tile-content>
                </v-list-tile>
              </v-list>
              <v-divider></v-divider>
              <div v-if="filterJobList(contractor.jobsList).length > 0">
                <v-list two-line>
                  <v-list-tile
                    v-for="jobID in contractor.jobsList"
                    :key="jobID"
                    v-show="filterJobList(contractor.jobsList).includes(jobID)"
                    avatar
                    @click="navigateToJob(jobID)"
                  >
                    <v-list-tile-avatar class="pl-5">
                      <v-icon>work</v-icon>
                    </v-list-tile-avatar>

                    <v-list-tile-content>
                      <v-list-tile-title
                        class="title font-weight-light"
                      >{{ getJobObject(jobID).name }}</v-list-tile-title>
                      <v-list-tile-sub-title
                        class="subheading font-weight-light"
                      >{{ getAddress(jobID) }}</v-list-tile-sub-title>
                    </v-list-tile-content>
                  </v-list-tile>
                </v-list>
              </div>
            </div>
          </v-card>
        </v-flex>
      </v-layout>
    </v-container>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { Job, Contact, Contractor } from "../basecoat_message_pb";

export default Vue.extend({
  data: function() {
    return {
      showEmptyContractors: false,
      contractorList: [] as Contractor.AsObject[],
      jobMap: {} as { [key: string]: Job },
      jobFilter: [] as string[]
    };
  },
  mounted() {
    this.contractorMapToList();
    this.jobMap = this.$store.state.jobData;
  },
  created() {
    this.$store.subscribe((mutation, state) => {
      if (mutation.type === "updateJobDataFilter") {
        this.jobFilter = this.$store.state.jobDataFilter;
      }
    });

    this.$store.subscribe((mutation, state) => {
      if (
        mutation.type === "updateJobData" ||
        mutation.type === "updateContractorData"
      ) {
        this.contractorMapToList();
        this.jobMap = this.$store.state.jobData;
      }
    });
  },
  methods: {
    navigateToJob: function(jobID: string) {
      this.$router.push("/jobs/" + jobID);
    },
    navigateToContractor: function(contractorID: string) {
      this.$router.push("/jobs/contractors/" + contractorID);
    },
    contractorMapToList: function() {
      let contractorDataMap: { [key: string]: Contractor } = this.$store.state
        .contractorData;
      let newContractorList: Contractor.AsObject[] = [];

      for (const [key, value] of Object.entries(contractorDataMap)) {
        newContractorList.push(value.toObject());
      }
      newContractorList.push(this.insertOrphanedJobs());
      this.contractorList = newContractorList;
    },
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
    },
    filterJobList: function(jobsList: string[]) {
      let filteredJobList: string[] = [];
      let self = this;
      jobsList.forEach(function(jobID, index) {
        if (self.jobFilter.length === 0 || self.jobFilter.includes(jobID)) {
          filteredJobList.push(jobID);
        }
      });

      return filteredJobList;
    },
    insertOrphanedJobs: function() {
      let jobDataMap: { [key: string]: Job } = this.$store.state.jobData;
      let emptyContractor: Contractor.AsObject = {
        id: "",
        company: "",
        contact: {
          name: "",
          phone: "",
          email: ""
        } as Contact.AsObject,
        jobsList: []
      };

      for (const [key, value] of Object.entries(jobDataMap)) {
        let currentJob = value.toObject();
        if (currentJob.contractorId === "") {
          emptyContractor.jobsList.push(key);
        }
      }
      return emptyContractor;
    }
  }
});
</script>

<style scoped>
div.v-card {
  margin-bottom: 0.3%;
}

.sizer {
  width: 65%;
}

div.v-list__tile__title {
  height: auto;
}

div.v-list__tile__sub-title {
  height: auto;
}

.title {
  line-height: normal !important;
}

.pull-right {
  justify-content: flex-end;
}

.v-input--selection-controls {
  margin-top: 0%;
  padding-top: 0%;
}
</style>
