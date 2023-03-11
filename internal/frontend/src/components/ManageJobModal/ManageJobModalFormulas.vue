<template>
  <div>
    <v-spacer v-show="job.formulasList.length === 0">No Formulas Listed</v-spacer>
    <h2
      class="font-weight-light text-center"
      v-show="job.formulasList.length > 0"
      justify-space-around
    >
      Formulas
      <v-tooltip right>
        <template v-slot:activator="{ on }">
          <v-icon small color="text--secondary" style="vertical-align: middle" v-on="on">info</v-icon>
        </template>
        <span>Formulas linked to this job</span>
      </v-tooltip>
    </h2>
    <div>
      <v-layout>
        <v-list two-line v-show="formMode === 'view'" style="width:100%;">
          <template v-for="(formulaID, index) in job.formulasList">
            <v-list-tile v-bind:key="`formula-tile-${index}`">
              <v-list-tile-avatar>
                <v-icon>invert_colors</v-icon>
              </v-list-tile-avatar>
              <v-list-tile-content>
                <v-list-tile-title v-text="$store.state.formulaData[formulaID].getName()"></v-list-tile-title>
                <v-list-tile-sub-title v-text="$store.state.formulaData[formulaID].getNumber()"></v-list-tile-sub-title>
              </v-list-tile-content>
            </v-list-tile>
            <v-divider v-bind:key="`formula-divider-${index}`"></v-divider>
          </template>
        </v-list>

        <v-flex xs12 v-show="formMode === 'edit'">
          <v-autocomplete
            v-model="job.formulasList"
            :items="formulaDataToList"
            item-text="name"
            item-value="id"
            hide-selected
            label="Link formulas to this job"
            placeholder="Start typing to Search"
            multiple
            clearable
            counter
          >
            <template slot="item" slot-scope="data">
              <v-list-tile-content>
                <v-list-tile-title v-html="data.item.name"></v-list-tile-title>
                <v-list-tile-sub-title v-html="data.item.number"></v-list-tile-sub-title>
              </v-list-tile-content>
            </template>
          </v-autocomplete>
        </v-flex>
      </v-layout>
    </div>
    <br />
    <v-divider></v-divider>
    <br />
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import { Formula } from "../../basecoat_message_pb";

interface modifiedFormula {
  id: string;
  name: string;
  number: string;
}

export default Vue.extend({
  props: ["formMode", "job"],
  data: function() {
    return {};
  },
  computed: {
    formulaDataToList: function() {
      let formulaDataMap: { [key: string]: Formula } = this.$store.state
        .formulaData;
      let formulaDataList: Formula[] = [];

      for (const [key, value] of Object.entries(formulaDataMap)) {
        formulaDataList.push(value);
      }
      let modifiedFormulaList: modifiedFormula[] = [];
      let formula: Formula;

      for (formula of formulaDataList) {
        let modifiedFormula: modifiedFormula = {
          id: "",
          name: "",
          number: ""
        };

        modifiedFormula.id = formula.getId();
        modifiedFormula.name = formula.getName();
        modifiedFormula.number = formula.getNumber();

        modifiedFormulaList.push(modifiedFormula);
      }

      return modifiedFormulaList;
    }
  }
});
</script>
<style scoped>
</style>
