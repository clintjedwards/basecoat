<template>
  <div>
    <job-search-panel class="search-panel"></job-search-panel>
    <job-table class="data-table"></job-table>
    <router-view></router-view>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import JobSearchPanel from "./JobSearchPanel.vue";
import JobTable from "./JobTable.vue";

export default Vue.extend({
  components: {
    JobTable,
    JobSearchPanel
  },
  data: function() {
    return {};
  }
});
</script>

<style scoped>
</style>
