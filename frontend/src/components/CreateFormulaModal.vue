<template>
  <v-layout row justify-center>
    <v-dialog v-model="showModal" max-width="600px">
      <v-form ref="createFormulaForm" lazy-validation>
        <v-card>
          <v-card-title>
            <span class="headline">
              <slot name="heading"></slot>
            </span>
          </v-card-title>
          <v-card-text>
            <!-- Formula Info -->
            <v-container grid-list-md>
              <v-spacer>Formula Info</v-spacer>
              <v-layout wrap>
                <v-flex xs12 sm6 md6>
                  <v-text-field
                    label="Formula Name"
                    hint="Must be unique"
                    :rules="nameRules"
                    v-model="formulaData.name"
                    required
                  ></v-text-field>
                </v-flex>
                <v-flex xs12 sm6 md6>
                  <v-text-field
                    label="Formula Number"
                    hint="Custom number used to reference formula"
                    v-model="formulaData.number"
                    required
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <v-divider></v-divider>
              <br />

              <!-- Base -->
              <v-spacer>Base</v-spacer>
              <v-layout
                row
                wrap
                v-for="(base, index) in formulaData.basesList"
                v-bind:key="`base-${index}`"
              >
                <v-flex xs10 sm9>
                  <v-text-field label="Base Name" v-model="base.name"></v-text-field>
                </v-flex>
                <v-flex xs2 sm3>
                  <v-text-field
                    label="Size"
                    v-model="base.amount"
                    append-outer-icon="delete"
                    v-on:click:append-outer="removeBaseField(index)"
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <v-btn flat color="primary" v-on:click="addBaseField">Add Base</v-btn>
              <v-divider></v-divider>
              <br />

              <!-- Colorants -->
              <v-spacer>Colorants</v-spacer>
              <v-flex xs6 sm6 v-show="formulaData.colorantsList.length != 0">
                <v-select
                  :items="colorantTypesToList"
                  label="Colorant Type"
                  v-model="currentColorantType"
                  v-on:change="fillColorantTypes(currentColorantType)"
                  clearable
                ></v-select>
              </v-flex>
              <v-layout
                row
                wrap
                v-for="(colorant, index) in formulaData.colorantsList"
                v-bind:key="`colorant-${index}`"
              >
                <v-flex xs2 sm2>
                  <v-text-field label="Type" v-model="colorant.type"></v-text-field>
                </v-flex>
                <v-flex xs7 sm7>
                  <v-text-field label="Colorant Name" v-model="colorant.name"></v-text-field>
                </v-flex>
                <v-flex xs3 sm3>
                  <v-text-field
                    label="Amount"
                    v-model="colorant.amount"
                    append-outer-icon="delete"
                    v-on:click:append-outer="removeColorantField(index)"
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <v-btn flat color="primary" v-on:click="addColorantField">Add Colorant</v-btn>
              <v-divider></v-divider>
              <br />

              <!-- Jobs -->
              <v-spacer>Jobs</v-spacer>
              <v-layout>
                <v-flex xs12>
                  <v-autocomplete
                    v-model="formulaData.jobsList"
                    :items="jobDataToList"
                    item-text="name"
                    item-value="id"
                    hide-selected
                    label="Link jobs to Formula"
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

              <!-- Notes -->
              <v-layout>
                <v-flex xs12>
                  <v-textarea name="input-7-1" label="Notes" v-model="formulaData.notes" auto-grow></v-textarea>
                </v-flex>
              </v-layout>
            </v-container>
          </v-card-text>

          <!-- Buttons -->
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" flat v-on:click="handleCloseForm()">Close</v-btn>
            <v-btn color="blue darken-1" flat v-on:click="handleCreateFormula()">Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script lang="ts">
import Vue from "vue";
import { Formula, Base, Colorant, Job } from "../basecoat_message_pb";
import { SnackBarColor, SnackBar } from "../snackbar";
import { CreateFormulaRequest } from "../basecoat_transport_pb";
import BasecoatClientWrapper from "../basecoatClientWrapper";

let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

let baseList: Base.AsObject[] = [];
let colorantList: Colorant.AsObject[] = [];

let formulaData: CreateFormulaRequest.AsObject = {
  name: "",
  number: "",
  jobsList: [],
  basesList: baseList,
  colorantsList: colorantList,
  notes: ""
};

interface modifiedJob {
  id: string;
  name: string;
  street: string;
}

export default Vue.extend({
  data: function() {
    return {
      showModal: true,
      formulaData: formulaData,
      currentColorantType: "",
      nameRules: [
        function(v: string) {
          if (!!v) {
            return true;
          }
          return "Formula Name is required";
        }
      ]
    };
  },
  computed: {
    jobDataToList(): modifiedJob[] {
      let jobDataMap: { [key: string]: Job } = this.$store.state.jobData;
      let jobDataList: modifiedJob[] = [];
      for (const [key, value] of Object.entries(jobDataMap)) {
        let job: modifiedJob;
        let street: string = "";
        let jobAddress = value.getAddress();
        if (jobAddress != undefined) {
          street = jobAddress.getStreet();
        }

        job = {
          id: value.getId(),
          name: value.getName(),
          street: street
        };
        jobDataList.push(job);
      }
      return jobDataList;
    },
    colorantTypesToList(): string[] {
      let colorantTypeMap: { [key: string]: string } = this.$store.state
        .colorantTypes;
      let colorantTypeList: string[] = [];

      for (const [key, value] of Object.entries(colorantTypeMap)) {
        let colorantType: string;
        colorantType = key;
        colorantTypeList.push(colorantType);
      }

      return colorantTypeList;
    }
  },
  methods: {
    addBaseField: function() {
      this.formulaData.basesList.push({
        type: "",
        name: "",
        amount: ""
      });
    },
    addColorantField: function() {
      let type: string = "";

      if (this.currentColorantType != "") {
        type = this.currentColorantType;
      }

      this.formulaData.colorantsList.push({
        type: type,
        name: "",
        amount: ""
      });
    },
    fillColorantTypes: function(type: string) {
      this.formulaData.colorantsList.forEach(function(colorant) {
        colorant.type = type;
      });
    },
    removeColorantField: function(index: number) {
      this.formulaData.colorantsList.splice(index, 1);
    },
    removeBaseField: function(index: number) {
      this.formulaData.basesList.splice(index, 1);
    },
    clearForm: function() {
      (this.$refs.createFormulaForm as HTMLFormElement).reset();
      this.formulaData.basesList = [];
      this.formulaData.colorantsList = [];
    },
    handleCloseForm: function() {
      this.$router.push({ name: "formulas" });
    },
    handleCreateFormula: function() {
      if ((this.$refs.createFormulaForm as HTMLFormElement).validate()) {
        client
          .submitCreateFormulaForm(this.formulaData)
          .then(() => {
            client
              .getFormulaData()
              .then(formulas => {
                this.$store.commit("updateFormulaData", formulas);
                this.clearForm();
                this.$router.push({ name: "formulas" });
              })
              .catch(() => {
                this.$store.commit("updateSnackBar", {
                  text: "Could not load formulas",
                  color: SnackBarColor.Error,
                  display: true
                } as SnackBar);
                this.clearForm();
                this.$router.push({ name: "formulas" });
              });
          })
          .catch(() => {
            this.$store.commit("updateSnackBar", {
              text: "Could not create formula",
              color: SnackBarColor.Error,
              display: true
            } as SnackBar);
          });
      }
    }
  }
});
</script>
