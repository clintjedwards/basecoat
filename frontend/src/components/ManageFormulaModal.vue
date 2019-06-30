<template>
  <v-layout row justify-center>
    <v-dialog v-model="$store.state.displayManageFormulaModal" max-width="600px" persistent>
      <v-form ref="manageFormulaForm" lazy-validation>
        <v-card>
          <v-card-title>
            <span>ID: {{ formulaData.id }}</span>
          </v-card-title>
          <v-card-text>
            <!-- Formula Info -->
            <v-container grid-list-md>
              <v-spacer>Formula Info</v-spacer>
              <v-layout align-center justify-center wrap>
                <v-flex xs12 sm12 md12 v-show="formMode === 'view'">
                  <h2 class="display-3 font-weight-light text-capitalize">{{ formulaData.name }}</h2>
                </v-flex>
                <v-flex xs12 sm6 md6 v-show="formMode === 'edit'">
                  <v-text-field
                    label="Formula Name"
                    :rules="nameRules"
                    v-model="formulaData.name"
                    required
                  ></v-text-field>
                </v-flex>
                <v-flex xs12 sm8 md8 v-show="formMode === 'view'">
                  <h4 class="display-1 font-weight-light text-capitalize">{{ formulaData.number }}</h4>
                </v-flex>
                <v-flex xs12 sm6 md6 v-show="formMode === 'edit'">
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
              <v-spacer v-show="formulaData.basesList.length === 0">No Bases Listed</v-spacer>
              <v-spacer v-show="formulaData.basesList.length > 0">Base</v-spacer>
              <v-layout
                row
                wrap
                v-for="(base, index) in formulaData.basesList"
                v-bind:key="`base-${index}`"
              >
                <v-list v-show="formMode === 'view'" style="width:100%;">
                  <v-list-tile>
                    <v-list-tile-content>
                      <v-list-tile-title v-text="base.name"></v-list-tile-title>
                    </v-list-tile-content>
                    <v-list-tile-avatar>
                      <template>{{ base.amount }}</template>
                    </v-list-tile-avatar>
                  </v-list-tile>
                  <v-divider></v-divider>
                </v-list>

                <v-flex xs10 sm9 v-show="formMode === 'edit'">
                  <v-text-field label="Base Name" v-model="base.name"></v-text-field>
                </v-flex>
                <v-flex xs2 sm3 v-show="formMode === 'edit'">
                  <v-text-field
                    label="Amount"
                    v-model="base.amount"
                    append-outer-icon="delete"
                    v-on:click:append-outer="removeBaseField(index)"
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <div v-show="formMode === 'edit'">
                <v-btn flat color="primary" v-on:click="addBaseField">Add Base</v-btn>
              </div>
              <br />
              <br />

              <!-- Colorants -->
              <v-spacer v-show="formulaData.colorantsList.length === 0">No Colorants Listed</v-spacer>
              <v-spacer v-show="formulaData.colorantsList.length > 0">Colorants</v-spacer>

              <v-list v-show="formMode === 'view' && colorantOverallTypeSet" style="width:100%;">
                <v-list-tile>
                  <v-list-tile-avatar tile>
                    <v-img max-height="25" max-width="30" v-bind:src="colorantTypeInfo.imageURL"></v-img>
                  </v-list-tile-avatar>
                  <v-list-tile-content>
                    <v-list-tile-title v-text="colorantTypeInfo.userMessage"></v-list-tile-title>
                  </v-list-tile-content>
                </v-list-tile>
                <v-divider></v-divider>
              </v-list>
              <v-flex xs6 sm6 v-show="formulaData.colorantsList.length != 0 && formMode === 'edit'">
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
                <v-list v-show="formMode === 'view'" style="width:100%;">
                  <v-list-tile>
                    <v-list-tile-avatar>
                      <v-icon>invert_colors</v-icon>
                    </v-list-tile-avatar>
                    <v-list-tile-content>
                      <v-list-tile-sub-title
                        v-show="currentColorantType === ''"
                        v-text="colorant.type"
                      ></v-list-tile-sub-title>
                      <v-list-tile-title v-text="colorant.name"></v-list-tile-title>
                    </v-list-tile-content>
                    <v-list-tile-avatar>
                      <template>{{ colorant.amount }}</template>
                    </v-list-tile-avatar>
                  </v-list-tile>
                  <v-divider></v-divider>
                </v-list>

                <v-flex xs2 sm2 v-show="formMode === 'edit'">
                  <v-text-field label="Type" v-model="colorant.type"></v-text-field>
                </v-flex>
                <v-flex xs7 sm7 v-show="formMode === 'edit'">
                  <v-text-field label="Colorant Name" v-model="colorant.name"></v-text-field>
                </v-flex>
                <v-flex xs3 sm3 v-show="formMode === 'edit'">
                  <v-text-field
                    label="Amount"
                    v-model="colorant.amount"
                    append-outer-icon="delete"
                    v-on:click:append-outer="removeColorantField(index)"
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <div v-show="formMode === 'edit'">
                <v-btn flat color="primary" v-on:click="addColorantField">Add Colorant</v-btn>
              </div>
              <br />

              <!-- Jobs -->
              <v-spacer v-show="formulaData.jobsList.length === 0">No Jobs Listed</v-spacer>
              <v-spacer v-show="formulaData.jobsList.length > 0">Jobs</v-spacer>

              <v-layout>
                <v-list two-line v-show="formMode === 'view'" style="width:100%;">
                  <template v-for="(jobID, index) in formulaData.jobsList">
                    <v-list-tile v-bind:key="`job-tile-${index}`">
                      <v-list-tile-avatar>
                        <v-icon>work</v-icon>
                      </v-list-tile-avatar>
                      <v-list-tile-content>
                        <v-list-tile-title v-text="$store.state.jobData[jobID].getName()"></v-list-tile-title>
                        <v-list-tile-sub-title>
                          {{ $store.state.jobData[jobID].getStreet() }} {{ $store.state.jobData[jobID].getCity() }}
                          <template
                            v-if="$store.state.jobData[jobID].getCity() !== ''"
                          >,</template>
                          {{ $store.state.jobData[jobID].getState() }} {{ $store.state.jobData[jobID].getZipcode() }}
                        </v-list-tile-sub-title>
                      </v-list-tile-content>
                    </v-list-tile>
                    <v-divider v-bind:key="`job-divider-${index}`"></v-divider>
                  </template>
                </v-list>

                <v-flex xs12 v-show="formMode === 'edit'">
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
                <v-flex xs12 v-show="formMode === 'view'">
                  <v-spacer>Notes</v-spacer>
                  <pre>{{ formulaData.notes }}</pre>
                </v-flex>
                <v-flex xs12 v-show="formMode === 'edit'">
                  <v-textarea name="input-7-1" label="Notes" v-model="formulaData.notes" auto-grow></v-textarea>
                </v-flex>
              </v-layout>
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
            <v-btn
              color="blue darken-1"
              flat
              @click="$store.commit('hideManageFormulaModal'); setFormModeView();"
            >Close</v-btn>
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
              @click="parseColorantListForSameType(); setFormModeView();"
            >View</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'edit'"
              @click="getFormula(this.formulaData.id); parseColorantListForSameType();"
            >Reset</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'edit'"
              @click="handleFormSave(); parseColorantListForSameType();"
            >Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script lang="ts">
import Vue from "vue";
import {
  GetFormulaRequest,
  UpdateFormulaRequest,
  Base,
  Colorant,
  Job,
  Formula
} from "../basecoat_pb";

import { BasecoatClient } from "../BasecoatServiceClientPb";

let baseList: Base.AsObject[] = [];
let colorantList: Colorant.AsObject[] = [];

let formulaData: UpdateFormulaRequest.AsObject = {
  id: "",
  name: "",
  number: "",
  jobsList: [],
  basesList: baseList,
  colorantsList: colorantList,
  notes: ""
};

let formula: Formula;
declare var __API__: string;
let client = new BasecoatClient(__API__, null, null);

export default Vue.extend({
  data: function() {
    return {
      formMode: "view",
      showConfirmDelete: false,
      formulaData: formulaData,
      colorantOverallTypeSet: false,
      colorantTypeInfo: {},
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
  mounted() {
    if (this.$route.name === "formulaModal" && Vue.cookies.isKey("token")) {
      this.getFormula(this.$route.params.id);
      this.$store.commit("showManageFormulaModal");
    }
  },
  watch: {
    "formulaData.colorantsList": function() {
      this.parseColorantListForSameType();
    }
  },
  computed: {
    jobDataToList: function() {
      interface modifiedJob {
        id: string;
        name: string;
        street: string;
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
          street: ""
        };

        modifiedJob.id = job.getId();
        modifiedJob.name = job.getName();
        modifiedJob.street = job.getStreet();

        modifiedJobList.push(modifiedJob);
      }

      return modifiedJobList;
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
    setFormModeEdit: function() {
      this.formMode = "edit";
    },
    setFormModeView: function() {
      this.formMode = "view";
    },
    getFormula: function(formulaID: string) {
      let self = this;
      let getFormulaRequest = new GetFormulaRequest();
      getFormulaRequest.setId(formulaID);
      let metadata = { Authorization: "Bearer " + Vue.cookies.get("token") };
      client.getFormula(getFormulaRequest, metadata, function(err, response) {
        if (err) {
          console.log(err);
          self.$store.commit("displaySnackBar", "Could not load formula.");
          return;
        }
        self.populateFormData(response.getFormula());
      });
    },
    populateFormData: function(currentFormula: Formula | undefined) {
      if (currentFormula === undefined) {
        console.log(
          "could not load formula while trying to populate form data"
        );
        this.$store.commit("displaySnackBar", "Could not load formula.");
        return;
      }

      this.formulaData.id = currentFormula.getId();
      this.formulaData.name = currentFormula.getName();
      this.formulaData.number = currentFormula.getNumber();
      this.formulaData.notes = currentFormula.getNotes();
      this.formulaData.jobsList = currentFormula.getJobsList();

      //We need to format this as an object because the protomessage type resoves as a weird array
      let basesList: Base.AsObject[] = [];
      currentFormula
        .getBasesList()
        .forEach(function(item: Base, index: number) {
          let newBase: Base.AsObject;
          newBase = {
            type: item.getType(),
            name: item.getName(),
            amount: item.getAmount()
          };

          basesList.push(newBase);
        });

      //We need to format this as an object because the protomessage type resoves as a weird array
      let colorantsList: Colorant.AsObject[] = [];
      currentFormula
        .getColorantsList()
        .forEach(function(item: Colorant, index: number) {
          let newColorant: Colorant.AsObject;
          newColorant = {
            type: item.getType(),
            name: item.getName(),
            amount: item.getAmount()
          };

          colorantsList.push(newColorant);
        });

      this.formulaData.basesList = basesList;
      this.formulaData.colorantsList = colorantsList;
    },
    parseColorantListForSameType: function() {
      let self = this;

      if (self.formulaData.colorantsList.length < 1) {
        this.colorantOverallTypeSet = false;
        this.currentColorantType = "";
        return;
      }

      for (let i = 0; i < self.formulaData.colorantsList.length; ++i) {
        if (
          self.formulaData.colorantsList[i].type !=
          self.formulaData.colorantsList[0].type
        ) {
          this.colorantOverallTypeSet = false;
          this.currentColorantType = "";
          return;
        }
      }

      if (
        self.formulaData.colorantsList[0].type in
        this.$store.state.colorantTypes
      ) {
        this.colorantOverallTypeSet = true;
        this.colorantTypeInfo = this.$store.state.colorantTypes[
          self.formulaData.colorantsList[0].type
        ];
        this.currentColorantType = self.formulaData.colorantsList[0].type;
        return;
      }

      this.colorantOverallTypeSet = false;
      this.currentColorantType = "";
    },
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
      (this.$refs.manageFormulaForm as HTMLFormElement).reset();
      this.formulaData.basesList = [{ type: "", name: "", amount: "" }];
      this.formulaData.colorantsList = [{ type: "", name: "", amount: "" }];
    },
    handleFormSave: function() {
      if ((this.$refs.manageFormulaForm as HTMLFormElement).validate()) {
        this.$emit("submit-manage-formula-form", this.formulaData);
      }
    },
    handleFormDelete: function() {
      this.$emit("delete-formula", this.formulaData.id);
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
