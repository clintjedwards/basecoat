<template>
  <v-layout row>
    <v-dialog v-model="showModal" max-width="600px">
      <v-form v-if="!loading" ref="manageFormulaForm" lazy-validation>
        <v-card>
          <v-card-title>
            <span>ID: {{ formula.id }}</span>
          </v-card-title>
          <v-card-text>
            <!-- Formula Info -->
            <v-container>
              <manage-formula-modal-info :formMode="formMode" :formula.sync="formula"></manage-formula-modal-info>

              <!-- Base -->
              <manage-formula-modal-bases :formMode="formMode" :formula.sync="formula"></manage-formula-modal-bases>

              <!-- Colorants -->
              <manage-formula-modal-colorants :formMode="formMode" :formula.sync="formula"></manage-formula-modal-colorants>

              <!-- Jobs -->
              <manage-formula-modal-jobs :formMode="formMode" :formula.sync="formula"></manage-formula-modal-jobs>

              <!-- Notes -->
              <manage-formula-modal-notes :formMode="formMode" :formula.sync="formula"></manage-formula-modal-notes>
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
              @click="setFormModeView();"
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
              @click="handleUpdateFormula();"
            >Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script lang="ts">
import Vue from "vue";
import ManageFormulaModalInfo from "./ManageFormulaModalInfo.vue";
import ManageFormulaModalBases from "./ManageFormulaModalBases.vue";
import ManageFormulaModalColorants from "./ManageFormulaModalColorants.vue";
import ManageFormulaModalJobs from "./ManageFormulaModalJobs.vue";
import ManageFormulaModalNotes from "./ManageFormulaModalNotes.vue";

import { UpdateFormulaRequest } from "../../basecoat_transport_pb";
import { Base, Colorant, Job, Formula } from "../../basecoat_message_pb";
import BasecoatClientWrapper from "../../basecoatClientWrapper";
import { SnackBar, SnackBarColor } from "../../snackbar";

let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

export default Vue.extend({
  components: {
    ManageFormulaModalInfo,
    ManageFormulaModalBases,
    ManageFormulaModalColorants,
    ManageFormulaModalJobs,
    ManageFormulaModalNotes
  },
  data: function() {
    return {
      showModal: true,
      loading: true,
      formMode: "view",
      showConfirmDelete: false,
      formula: {
        id: "",
        name: "",
        number: "",
        jobsList: [],
        basesList: [] as Base.AsObject[],
        colorantsList: [] as Colorant.AsObject[],
        notes: ""
      } as UpdateFormulaRequest.AsObject
    };
  },
  async mounted() {
    let loadFormula = await this.getFormula(this.$route.params.id);
    if (loadFormula === undefined) {
      return;
    }
    this.formula = loadFormula.toObject();
    this.loading = false;
  },
  watch: {
    showModal(val) {
      if (!val) {
        this.$router.push({ name: "formulas" });
      }
    }
  },
  methods: {
    setFormModeEdit: function() {
      this.formMode = "edit";
    },
    setFormModeView: function() {
      this.formMode = "view";
    },
    getFormula: async function(formulaID: string) {
      try {
        let formula = await client.getFormula(formulaID);
        return formula;
      } catch (err) {
        this.$store.commit("updateSnackBar", {
          text: "Could not load formula",
          color: SnackBarColor.Error,
          display: true
        } as SnackBar);
      }
    },
    handleResetForm: async function() {
      let loadFormula = await this.getFormula(this.$route.params.id);
      if (loadFormula === undefined) {
        return;
      }
      this.formula = loadFormula.toObject();
    },
    handleCloseForm: function() {
      this.setFormModeView();
      this.$router.push({ name: "formulas" });
    },
    handleUpdateFormula: async function() {
      if ((this.$refs.manageFormulaForm as HTMLFormElement).validate()) {
        try {
          await client.submitManageFormulaForm(this.formula);
        } catch (err) {
          console.log(err);
          this.$store.commit("updateSnackBar", {
            text: "Could not update formula",
            color: SnackBarColor.Error,
            display: true
          } as SnackBar);
        }
      }

      this.$store.commit("updateSnackBar", {
        text: "Updated Formula: " + this.formula.name,
        color: SnackBarColor.Success,
        display: true
      } as SnackBar);
      this.handleCloseForm();

      // reload formula data
      client
        .getFormulaData()
        .then(formulas => {
          this.$store.commit("updateFormulaData", formulas);
        })
        .catch(() => {
          console.log("could not load formulas");
        });
    },
    handleFormDelete: function() {
      client
        .deleteFormula(this.formula.id)
        .then(() => {
          client
            .getFormulaData()
            .then(formulas => {
              this.$store.commit("updateFormulaData", formulas);
              this.handleCloseForm();
            })
            .catch(() => {
              this.handleCloseForm();
            });
        })
        .catch(() => {
          this.$store.commit("updateSnackBar", {
            text: "Could not delete formula",
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
h2 {
  text-align: center;
}

h4 {
  text-align: center;
  color: #9e9e9e;
}

pre {
  white-space: pre-wrap;
}
</style>
