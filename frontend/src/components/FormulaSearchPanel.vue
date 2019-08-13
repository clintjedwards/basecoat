<template>
  <div>
    <span>
      <v-badge color="#ff5252">
        <template v-slot:badge>{{ Object.keys($store.state.formulaData).length }}</template>
        <span class="display-3 font-weight-light search-panel-text">Formulas</span>
      </v-badge>
    </span>
    <span style="margin-left: 4em; width: 100em;">
      <v-text-field
        label="Search"
        prepend-icon="search"
        hint="Search for formula data of any kind"
        v-on:input="$store.commit('updateFormulaTableSearchTerm', $event)"
      ></v-text-field>
    </span>
    <span style="margin-left: 4em;">
      <v-btn color="primary" v-on:click="showCreateFormulaModal()">
        <v-icon>create</v-icon>Create Formula
      </v-btn>
    </span>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  data: function() {
    return {};
  },
  methods: {
    showCreateFormulaModal: function() {
      this.$router.push({ name: "createFormulaModal" });
    }
  }
});
</script>

<style scoped>
.search-panel-text {
  color: #9e9e9e;
}
</style>
