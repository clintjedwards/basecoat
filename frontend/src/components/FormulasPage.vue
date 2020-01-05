<template>
  <div>
    <formula-search-panel class="search-panel"></formula-search-panel>
    <formula-table class="data-table"></formula-table>
    <router-view></router-view>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import FormulaSearchPanel from "./FormulaSearchPanel.vue";
import FormulaTable from "./FormulaTable.vue";

export default Vue.extend({
  components: {
    FormulaTable,
    FormulaSearchPanel
  },
  data: function() {
    return {};
  }
});
</script>

<style scoped>
.data-table {
  margin-top: 5em;
  width: 77%;
  margin-left: auto;
  margin-right: auto;
}
</style>
