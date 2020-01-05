<template>
  <v-layout row>
    <v-dialog v-model="showModal" max-width="600px">
      <v-form v-if="!loading" ref="manageContractorForm" lazy-validation>
        <v-card>
          <v-card-title>
            <span>ID: {{ contractor.id }}</span>
          </v-card-title>
          <v-card-text>
            <v-container>
              <manage-contractor-modal-info :formMode="formMode" :contractor.sync="contractor"></manage-contractor-modal-info>
              <manage-contractor-modal-jobs :formMode="formMode" :contractor.sync="contractor"></manage-contractor-modal-jobs>
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
              @click="handleUpdateContractor()"
            >Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script lang="ts">
import Vue from "vue";
import ManageContractorModalInfo from "../ManageContractorModal/ManageContractorModalInfo.vue";
import ManageContractorModalJobs from "../ManageContractorModal/ManageContractorModalJobs.vue";

import {
  GetContractorRequest,
  UpdateContractorRequest
} from "../../basecoat_transport_pb";
import { Contact, Contractor } from "../../basecoat_message_pb";
import { SnackBarColor, SnackBar } from "../../snackbar";
import BasecoatClientWrapper from "../../basecoatClientWrapper";

let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

export default Vue.extend({
  components: {
    ManageContractorModalInfo,
    ManageContractorModalJobs
  },
  data: function() {
    return {
      showModal: true,
      loading: true,
      formMode: "view",
      showConfirmDelete: false,
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
    let loadContractor = await this.getContractor(this.$route.params.id);
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

    this.loading = false;
  },
  methods: {
    setFormModeEdit: function() {
      this.formMode = "edit";
    },
    setFormModeView: function() {
      this.formMode = "view";
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
      let loadContractor = await this.getContractor(this.$route.params.id);
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
    },
    handleCloseForm: function() {
      this.setFormModeView();
      this.$router.push({ name: "jobs" });
    },
    handleUpdateContractor: async function() {
      if ((this.$refs.manageContractorForm as HTMLFormElement).validate()) {
        try {
          await client.submitManageContractorForm(this.contractor);
        } catch (err) {
          console.log(err);
          this.$store.commit("updateSnackBar", {
            text: "Could not update contractor",
            color: SnackBarColor.Error,
            display: true
          } as SnackBar);
        }

        this.$store.commit("updateSnackBar", {
          text: "Updated Contractor: " + this.contractor.company,
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
        .deleteContractor(this.contractor.id)
        .then(() => {
          client
            .getContractorData()
            .then(contractors => {
              this.$store.commit("updateContractorData", contractors);
              this.handleCloseForm();
            })
            .catch(() => {
              this.handleCloseForm();
            });
        })
        .catch(() => {
          this.$store.commit("updateSnackBar", {
            text: "Could not delete contractor",
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
