<template>
  <div>
    <span>
      <v-badge color="#ff5252">
        <template v-slot:badge>{{ $store.state.totalJobs }}</template>
        <span class="display-3 font-weight-light search-panel-text">Jobs</span>
      </v-badge>
    </span>
    <span style="margin-left: 4em; width: 100em;">
      <v-text-field
        label="Search"
        prepend-icon="search"
        hint="Search for job data of any kind"
        v-on:input="$store.commit('updateJobTableSearchTerm', $event)"
      ></v-text-field>
    </span>
    <span style="margin-left: 4em;">
      <v-btn color="secondary" v-on:click="$store.commit('showAddJobModal')">
        <v-icon>create</v-icon>Add Job
      </v-btn>
    </span>
  </div>
</template>

<script>
export default {};
</script>

<style scoped>
.search-panel-text {
  color: #9e9e9e;
}
</style>
