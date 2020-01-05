<template>
  <v-container style="width: 80%;">
    <v-layout row justify-center>
      <v-flex md2 class="text-xs-center">
        <v-badge color="#ff5252">
          <template v-slot:badge>{{ Object.keys($store.state.formulaData).length }}</template>
          <span class="display-3 font-weight-light search-panel-text">Formulas</span>
        </v-badge>
      </v-flex>
      <v-flex md8>
        <v-text-field
          label="Search"
          prepend-icon="search"
          hint="Search for formula data of any kind"
          v-on:input="debounceInput($event)"
        ></v-text-field>
      </v-flex>
      <v-flex md2 style="margin-left: 1em;">
        <v-layout column>
          <v-flex>
            <v-btn color="primary" v-on:click="showCreateFormulaModal()">
              <v-icon>create</v-icon>Create Formula
            </v-btn>
          </v-flex>
        </v-layout>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import _ from "lodash";

import BasecoatClientWrapper from "../basecoatClientWrapper";

let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

export default Vue.extend({
  data: function() {
    return {};
  },
  methods: {
    showCreateFormulaModal: function() {
      this.$router.push({ name: "createFormulaModal" });
    },
    debounceInput: _.debounce(function(this: any, searchTerm) {
      if (!searchTerm) {
        this.$store.commit("updateFormulaDataFilter", []);
        return;
      }
      client.searchFormulas(searchTerm).then(hits => {
        if (hits != undefined) {
          this.$store.commit("updateFormulaDataFilter", hits);
          return;
        }
        this.$store.commit("updateFormulaDataFilter", []);
        return;
      });
    }, 625)
  }
});
</script>

<style scoped>
.search-panel-text {
  color: #9e9e9e;
}
</style>
