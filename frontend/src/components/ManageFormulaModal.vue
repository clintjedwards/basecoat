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
                    hint="custom unique number pertaining to formula"
                    :rules="numberRules"
                    v-model="formulaData.number"
                    required
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <v-divider></v-divider>
              <br>

              <!-- Base -->
              <v-spacer>Base</v-spacer>
              <v-layout
                row
                wrap
                v-for="(base, index) in formulaData.base"
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
              <br>
              <br>

              <!-- Colorants -->
              <v-spacer>Colorants</v-spacer>
              <v-layout
                row
                wrap
                v-for="(colorant, index) in formulaData.colorants"
                v-bind:key="`colorant-${index}`"
              >
                <v-list v-show="formMode === 'view'" style="width:100%;">
                  <v-list-tile>
                    <v-list-tile-content>
                      <v-list-tile-title v-text="colorant.name"></v-list-tile-title>
                    </v-list-tile-content>
                    <v-list-tile-avatar>
                      <template>{{ colorant.amount }}</template>
                    </v-list-tile-avatar>
                  </v-list-tile>
                  <v-divider></v-divider>
                </v-list>

                <v-flex xs9 sm9 v-show="formMode === 'edit'">
                  <v-text-field label="Colorant Name" v-model="colorant.name"></v-text-field>
                </v-flex>
                <v-flex xs2 sm3 v-show="formMode === 'edit'">
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
              <br>

              <!-- Jobs -->
              <v-spacer>Jobs</v-spacer>
              <v-layout>
                <v-list two-line v-show="formMode === 'view'" style="width:100%;">
                  <template v-for="(job, index) in formulaData.jobs">
                    <v-list-tile v-bind:key="`job-tile-${index}`">
                      <v-list-tile-content>
                        <v-list-tile-title v-text="jobDictionary[job].name"></v-list-tile-title>
                        <v-list-tile-sub-title>{{ jobDictionary[job].street }} {{ jobDictionary[job].city }}, {{ jobDictionary[job].state }} {{ jobDictionary[job].zipcode }}</v-list-tile-sub-title>
                      </v-list-tile-content>
                    </v-list-tile>
                    <v-divider v-bind:key="`job-divider-${index}`"></v-divider>
                  </template>
                </v-list>

                <v-flex xs12 v-show="formMode === 'edit'">
                  <v-autocomplete
                    v-model="formulaData.jobs"
                    :items="$store.state.jobData"
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
            <v-spacer></v-spacer>
            <v-btn
              color="blue darken-1"
              flat
              @click="$store.commit('hideManageFormulaModal'); setFormModeView(); populateFormData();"
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
              @click="setFormModeView()"
            >View</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'edit'"
              @click="populateFormData()"
            >Reset</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'edit'"
              @click="handleFormSave()"
            >Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script>
export default {
  props: ["formulaInView"],
  data: function() {
    return {
      formMode: "view",
      formulaData: {},
      nameRules: [v => !!v || "Formula Name is required"],
      numberRules: [v => !!v || "Formula Number is required"]
    };
  },
  watch: {
    formulaInView: function() {
      this.populateFormData();
    }
  },
  computed: {
    jobDictionary: function() {
      let jobDictionary = {};

      this.$store.state.jobData.forEach(function(job) {
        jobDictionary[job.id] = job;
      });

      return jobDictionary;
    }
  },
  methods: {
    setFormModeEdit: function() {
      this.formMode = "edit";
    },
    setFormModeView: function() {
      this.formMode = "view";
    },
    populateFormData: function() {
      this.formulaData = JSON.parse(JSON.stringify(this.formulaInView));
      this.formulaData.base = [];
      this.formulaData.colorants = [];

      for (var key in this.formulaInView.base) {
        this.formulaData.base.push({
          name: key,
          amount: this.formulaInView.base[key]
        });
      }

      for (var key in this.formulaInView.colorants) {
        this.formulaData.colorants.push({
          name: key,
          amount: this.formulaInView.colorants[key]
        });
      }
    },
    addBaseField: function() {
      this.formulaData.base.push({
        name: "",
        amount: ""
      });
    },
    addColorantField: function() {
      this.formulaData.colorants.push({
        name: "",
        amount: ""
      });
    },
    removeColorantField: function(index) {
      this.formulaData.colorants.splice(index, 1);
    },
    removeBaseField: function(index) {
      this.formulaData.base.splice(index, 1);
    },
    clearForm: function() {
      this.$refs.manageFormulaForm.reset();
      this.formulaData.base = [{ name: "", amount: "" }];
      this.formulaData.colorants = [{ name: "", amount: "" }];
    },
    handleFormSave: function() {
      let newFormulaData = this.flattenFormulaData(this.formulaData);
      if (this.$refs.manageFormulaForm.validate()) {
        this.$emit("submit-manage-formula-form", newFormulaData);
      }
    },
    flattenFormulaData: function(formulaData) {
      var flattenedData = JSON.parse(JSON.stringify(this.formulaData));
      flattenedData.base = {};
      flattenedData.colorants = {};

      formulaData.base.forEach(function(base) {
        if (base.name == "") {
          return;
        }
        flattenedData.base[base.name] = base.amount;
      });

      formulaData.colorants.forEach(function(colorant) {
        if (colorant.name == "") {
          return;
        }
        flattenedData.colorants[colorant.name] = colorant.amount;
      });

      return flattenedData;
    }
  }
};
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
