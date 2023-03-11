<template>
  <v-layout row>
    <v-dialog v-model="showModal" max-width="600px">
      <v-form v-if="!loading" ref="manageJobForm" lazy-validation>
        <v-card>
          <v-card-title>
            <span>ID: {{ job.id }}</span>
          </v-card-title>
          <v-card-text>
            <!-- Contractor Info -->
            <v-container>
              <manage-job-modal-contractor :formMode="formMode" :contractor.sync="contractor"></manage-job-modal-contractor>

              <!-- Job Info -->
              <h2 class="font-weight-light text-center">Job Site</h2>
              <manage-job-modal-site-info :formMode="formMode" :job.sync="job"></manage-job-modal-site-info>
              <br />
              <v-divider></v-divider>
              <br />

              <!-- Address -->
              <manage-job-modal-address :formMode="formMode" :address.sync="job.address"></manage-job-modal-address>

              <!-- Formulas -->
              <manage-job-modal-formulas :formMode="formMode" :job.sync="job"></manage-job-modal-formulas>

              <!-- Notes-->
              <manage-job-modal-notes :formMode="formMode" :job.sync="job"></manage-job-modal-notes>
            </v-container>
          </v-card-text>

          <!-- Buttons -->
          <v-card-actions>
            <v-btn
              color="error darken-1"
              flat
              v-show="formMode === 'view' && !showConfirmDelete"
              @click="showConfirmDelete = true"
            >Delete</v-btn>
            <v-btn
              color="error darken-1"
              flat
              v-show="formMode === 'view' && showConfirmDelete"
              @click="handleFormDelete()"
            >Confirm Delete</v-btn>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" flat @click="handleCloseForm(); setFormModeView();">Close</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'view'"
              @click="setFormModeEdit()"
            >Edit</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'edit'"
              @click="setFormModeView()"
            >View</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'edit'"
              @click="handleResetForm()"
            >Reset</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'edit'"
              @click="handleUpdateJob()"
            >Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script lang="ts">
import Vue from "vue";
import ManageJobModalSiteInfo from "./ManageJobModalSiteInfo.vue";
import ManageJobModalContractor from "./ManageJobModalContractor.vue";
import ManageJobModalAddress from "./ManageJobModalAddress.vue";
import ManageJobModalFormulas from "./ManageJobModalFormulas.vue";
import ManageJobModalNotes from "./ManageJobModalNotes.vue";

import { GetJobRequest, UpdateJobRequest } from "../../basecoat_transport_pb";
import {
  Contact,
  Job,
  Formula,
  Address,
  Contractor
} from "../../basecoat_message_pb";
import { SnackBarColor, SnackBar } from "../../snackbar";
import BasecoatClientWrapper from "../../basecoatClientWrapper";

let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

export default Vue.extend({
  components: {
    ManageJobModalSiteInfo,
    ManageJobModalContractor,
    ManageJobModalAddress,
    ManageJobModalFormulas,
    ManageJobModalNotes
  },
  data: function() {
    return {
      showModal: true,
      loading: true,
      formMode: "view",
      showConfirmDelete: false,
      job: {
        id: "",
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
      } as UpdateJobRequest.AsObject,
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
  watch: {
    showModal(val) {
      if (!val) {
        this.$router.push({ name: "jobs" });
      }
    }
  },
  async mounted() {
    let loadJob = await this.getJob(this.$route.params.id);
    if (loadJob === undefined) {
      return;
    }
    this.job = loadJob.toObject();
    if (this.job.address === undefined) {
      this.job.address = {
        street: "",
        street2: "",
        state: "",
        city: "",
        zipcode: ""
      } as Address.AsObject;
    }

    if (this.job.contact === undefined) {
      this.job.contact = {
        name: "",
        email: "",
        phone: ""
      } as Contact.AsObject;
    }

    if (this.job.contractorId !== "") {
      let loadContractor = await this.getContractor(this.job.contractorId);
      if (!loadContractor) {
        return;
      }
      this.contractor = loadContractor.toObject();
      if (this.contractor.contact === undefined) {
        this.contractor.contact = {
          name: "",
          email: "",
          phone: ""
        } as Contact.AsObject;
      }
    }

    this.loading = false;
  },
  methods: {
    setFormModeEdit: function() {
      this.formMode = "edit";
    },
    setFormModeView: function() {
      this.formMode = "view";
    },
    getJob: async function(jobID: string) {
      try {
        let job = await client.getJob(jobID);
        return job;
      } catch (err) {
        this.$store.commit("updateSnackBar", {
          text: "Could not load job",
          color: SnackBarColor.Error,
          display: true
        } as SnackBar);
      }
    },
    getContractor: async function(contractorID: string) {
      try {
        let contractor = await client.getContractor(contractorID);
        return contractor;
      } catch (err) {
        this.$store.commit("updateSnackBar", {
          text: "Could not load contractor",
          color: SnackBarColor.Error,
          display: true
        } as SnackBar);
      }
    },
    handleResetForm: async function() {
      let loadJob = await this.getJob(this.$route.params.id);
      if (loadJob === undefined) {
        return;
      }
      this.job = loadJob.toObject();
      if (this.job.address === undefined) {
        this.job.address = {
          street: "",
          street2: "",
          state: "",
          city: "",
          zipcode: ""
        } as Address.AsObject;
      }

      if (this.job.contact === undefined) {
        this.job.contact = {
          name: "",
          email: "",
          phone: ""
        } as Contact.AsObject;
      }

      if (this.job.contractorId !== "") {
        let loadContractor = await this.getContractor(this.job.contractorId);
        if (!loadContractor) {
          return;
        }
        this.contractor = loadContractor.toObject();

        if (this.contractor.contact === undefined) {
          this.contractor.contact = {
            name: "",
            email: "",
            phone: ""
          } as Contact.AsObject;
        }
      }
    },
    handleCloseForm: function() {
      this.setFormModeView();
      this.$router.push({ name: "jobs" });
    },
    handleUpdateJob: async function() {
      if ((this.$refs.manageJobForm as HTMLFormElement).validate()) {
        this.job.contractorId = this.contractor.id;

        try {
          await client.submitManageJobForm(this.job);
        } catch (err) {
          console.log(err);
          this.$store.commit("updateSnackBar", {
            text: "Could not update job",
            color: SnackBarColor.Error,
            display: true
          } as SnackBar);
        }

        this.$store.commit("updateSnackBar", {
          text: "Updated Job: " + this.job.name,
          color: SnackBarColor.Success,
          display: true
        } as SnackBar);
        this.handleCloseForm();

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
    handleFormDelete: function() {
      client
        .deleteJob(this.job.id)
        .then(() => {
          client
            .getJobData()
            .then(jobs => {
              this.$store.commit("updateJobData", jobs);
              this.handleCloseForm();
            })
            .catch(() => {
              this.handleCloseForm();
            });
        })
        .catch(() => {
          this.$store.commit("updateSnackBar", {
            text: "Could not delete job",
            color: SnackBarColor.Error,
            display: true
          } as SnackBar);
        });
      this.showConfirmDelete = false;
    }
  }
});
</script>

<style scoped>
pre {
  white-space: pre-wrap;
}

div.v-card__title {
  padding-bottom: 0%;
}
</style>
