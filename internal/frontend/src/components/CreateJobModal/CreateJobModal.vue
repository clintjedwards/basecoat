<template>
  <v-layout row justify-center>
    <v-dialog v-model="showModal" max-width="600px">
      <v-form ref="createJobForm" lazy-validation>
        <v-card>
          <v-card-title>
            <span class="headline">
              <slot name="heading"></slot>
            </span>
          </v-card-title>
          <v-card-text>
            <v-container>
              <!-- Contractor Info -->
              <create-job-modal-contractor :contractor.sync="contractor"></create-job-modal-contractor>
              <v-divider></v-divider>
              <br />

              <!-- Site Info-->
              <create-job-modal-site-info :job.sync="job"></create-job-modal-site-info>
              <v-divider></v-divider>
              <br />

              <!-- Address -->
              <create-job-modal-address :job.sync="job"></create-job-modal-address>
              <v-divider></v-divider>
              <br />

              <!-- Notes -->
              <create-job-modal-notes :job.sync="job"></create-job-modal-notes>
            </v-container>
          </v-card-text>

          <!-- Buttons -->
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" flat v-on:click="handleCloseForm()">Close</v-btn>
            <v-btn color="blue darken-1" flat v-on:click="handleCreateJob()">Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script lang="ts">
import Vue from "vue";
import CreateJobModalContractor from "./CreateJobModalContractor.vue";
import CreateJobModalSiteInfo from "./CreateJobModalSiteInfo.vue";
import CreateJobModalAddress from "./CreateJobModalAddress.vue";
import CreateJobModalNotes from "./CreateJobModalNotes.vue";
import {
  CreateJobRequest,
  CreateContractorRequest
} from "../../basecoat_transport_pb";
import { Contact, Address, Contractor } from "../../basecoat_message_pb";
import { SnackBarColor, SnackBar } from "../../snackbar";
import BasecoatClientWrapper from "../../basecoatClientWrapper";

let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

export default Vue.extend({
  components: {
    CreateJobModalContractor,
    CreateJobModalSiteInfo,
    CreateJobModalAddress,
    CreateJobModalNotes
  },
  data: function() {
    return {
      showModal: true,
      job: {
        name: "",
        address: {
          street: "",
          street2: "",
          state: "",
          city: "",
          zipcode: ""
        } as Address.AsObject,
        notes: "",
        formulasList: [],
        contractorId: "",
        contact: {
          name: "",
          email: "",
          phone: ""
        } as Contact.AsObject
      } as CreateJobRequest.AsObject,
      contractor: {
        id: "",
        company: "",
        contact: {
          name: "",
          email: "",
          phone: ""
        } as Contact.AsObject
      } as Contractor.AsObject
    };
  },
  methods: {
    clearForm: function() {
      (this.$refs.createJobForm as HTMLFormElement).reset();
    },
    handleCloseForm: function() {
      this.$router.push({ name: "jobs" });
    },
    // if the contractor has an id we can add it in as the contractorID and create
    // job as usual, if not we'll need to wait for the contractor to be created first
    handleCreateJob: async function() {
      if ((this.$refs.createJobForm as HTMLFormElement).validate()) {
        this.job.contractorId = this.contractor.id;

        try {
          await client.submitCreateJobForm(this.job);
        } catch (err) {
          console.log(err);
          this.$store.commit("updateSnackBar", {
            text: "Could not create job",
            color: SnackBarColor.Error,
            display: true
          } as SnackBar);
          return;
        }

        this.$store.commit("updateSnackBar", {
          text: "Created Job: " + this.job.name,
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
    }
  }
});
</script>
