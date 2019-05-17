<template>
  <v-layout row justify-center>
    <v-dialog v-model="$store.state.displayCreateFormulaModal" max-width="600px" persistent>
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
                <v-flex xs10 sm9>
                  <v-text-field label="Base Name" v-model="base.name"></v-text-field>
                </v-flex>
                <v-flex xs2 sm3>
                  <v-text-field
                    label="Amount"
                    v-model="base.amount"
                    append-outer-icon="delete"
                    v-on:click:append-outer="removeBaseField(index)"
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <v-btn flat color="primary" v-on:click="addBaseField">Add Base</v-btn>
              <v-divider></v-divider>
              <br>

              <!-- Colorants -->
              <v-spacer>Colorants</v-spacer>
              <v-layout
                row
                wrap
                v-for="(colorant, index) in formulaData.colorants"
                v-bind:key="`colorant-${index}`"
              >
                <v-flex xs9 sm9>
                  <v-text-field label="Colorant Name" v-model="colorant.name"></v-text-field>
                </v-flex>
                <v-flex xs2 sm3>
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
              <br>

              <!-- Jobs -->
              <v-spacer>Jobs</v-spacer>
              <v-layout>
                <v-flex xs12>
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
                <v-flex xs12>
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
              v-on:click="$store.commit('hideCreateFormulaModal')"
            >Close</v-btn>
            <v-btn color="blue darken-1" flat v-on:click="handleFormSave()">Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script>
export default {
  data: function() {
    return {
      formulaData: {
        name: "",
        number: "",
        jobs: [],
        base: [
          {
            name: "",
            amount: ""
          }
        ],
        colorants: [
          {
            name: "",
            amount: ""
          }
        ],
        notes: ""
      },
      nameRules: [v => !!v || "Formula Name is required"],
      numberRules: [v => !!v || "Formula Number is required"]
    };
  },
  methods: {
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
      this.$refs.createFormulaForm.reset();
      this.formulaData.base = [{ name: "", amount: "" }];
      this.formulaData.colorants = [{ name: "", amount: "" }];
    },
    handleFormSave: function() {
      let newFormulaData = this.flattenFormulaData(this.formulaData);
      if (this.$refs.createFormulaForm.validate()) {
        this.$emit("submit-create-form", newFormulaData);
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
