<template>
  <v-layout row justify-center>
    <v-dialog v-model="showModal" max-width="600px">
      <v-form ref="createContractorForm" lazy-validation>
        <v-card>
          <v-card-text>
            <v-container>
              <h2 class="font-weight-light text-center" justify-space-around>
                Contractor
                <v-tooltip right>
                  <template v-slot:activator="{ on }">
                    <v-icon
                      small
                      color="text--secondary"
                      style="vertical-align: middle"
                      v-on="on"
                    >info</v-icon>
                  </template>
                  <span>
                    A contractor is the company who requests the formula.
                    <br />A contractor might make the request for a formula at many job sites.
                  </span>
                </v-tooltip>
              </h2>
              <v-layout wrap>
                <v-flex xs12 sm12 md12>
                  <v-text-field
                    label="Contractor Name"
                    :rules="nameRules"
                    v-model="contractor.company"
                    required
                  ></v-text-field>
                </v-flex>
                <v-flex xs12 sm4 md4>
                  <v-text-field label="Contractor Name" v-model="contractor.contact.name"></v-text-field>
                </v-flex>
                <v-flex xs12 sm4 md4>
                  <v-text-field label="Contractor Email" v-model="contractor.contact.email"></v-text-field>
                </v-flex>
                <v-flex xs12 sm4 md4>
                  <v-text-field label="Contractor Phone" v-model="contractor.contact.phone"></v-text-field>
                </v-flex>
              </v-layout>
              <!-- Jobs -->
              <v-spacer>
                Jobs
                <v-tooltip right>
                  <template v-slot:activator="{ on }">
                    <v-icon
                      small
                      color="text--secondary"
                      style="vertical-align: middle"
                      v-on="on"
                    >info</v-icon>
                  </template>
                  <span>
                    From here you can link jobs to this contractor.
                    <br />This list does not include jobs that already have a contractor.
                    <br />
                    <br />To include a job that has a contractor already linked, navigate to the job
                    <br />or contractor and unlink it first.
                  </span>
                </v-tooltip>
              </v-spacer>
              <v-layout>
                <v-flex xs12>
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
            </v-container>
          </v-card-text>

          <!-- Buttons -->
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" flat v-on:click="handleCloseForm()">Close</v-btn>
            <v-btn color="blue darken-1" flat v-on:click="handleCreateContractor()">Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>
<script lang="ts">
import Vue from "vue";
import { Contractor, Contact, Job } from "../../basecoat_message_pb";
import { CreateContractorRequest } from "../../basecoat_transport_pb";
import { SnackBarColor, SnackBar } from "../../snackbar";

import BasecoatClientWrapper from "../../basecoatClientWrapper";
let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

interface modifiedJob {
  id: string;
  name: string;
  street: string;
}

export default Vue.extend({
  data: function() {
    return {
      showModal: true,
      contractor: {
        company: "",
        contact: {
          name: "",
          email: "",
          phone: ""
        } as Contact.AsObject,
        jobsList: []
      } as CreateContractorRequest.AsObject,
      nameRules: [
        function(v: string) {
          if (!!v) {
            return true;
          }
          return "Contractor Name is required";
        }
      ]
    };
  },
  watch: {
    showModal(val) {
      if (!val) {
        this.$router.push({ name: "jobs" });
      }
    }
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

        if (contractorDataMap[value.getContractorId()] !== undefined) {
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
    handleCreateContractor: async function() {
      if ((this.$refs.createContractorForm as HTMLFormElement).validate()) {
        try {
          await client.submitCreateContractorForm(this.contractor);
        } catch (err) {
          console.log(err);
          this.$store.commit("updateSnackBar", {
            text: "Could not create contractor",
            color: SnackBarColor.Error,
            display: true
          } as SnackBar);
          return;
        }

        this.$store.commit("updateSnackBar", {
          text: "Created Contractor: " + this.contractor.company,
          color: SnackBarColor.Success,
          display: true
        } as SnackBar);
        this.clearForm();
        this.$router.push({ name: "jobs" });

        // reload job data
        client
          .getJobData()
          .then(jobs => {
            this.$store.commit("updateJobData", jobs);
          })
          .catch(() => {
            console.log("could not load jobs");
          });

        // reload contractor data
        client
          .getContractorData()
          .then(contractors => {
            this.$store.commit("updateContractorData", contractors);
          })
          .catch(() => {
            console.log("could not load contractors");
          });
      }
    },
    handleCloseForm: function() {
      this.$router.push({ name: "jobs" });
    },
    clearForm: function() {
      (this.$refs.createContractorForm as HTMLFormElement).reset();
    }
  }
});
</script>
<style scoped>
</style>
