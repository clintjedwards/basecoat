<template>
  <div>
    <h2 class="font-weight-light" justify-space-around>
      Miscellaneous
      <v-tooltip right>
        <template v-slot:activator="{ on }">
          <v-icon small color="text--secondary" style="vertical-align: middle" v-on="on">info</v-icon>
        </template>
        <span>Other information that might be included about this job</span>
      </v-tooltip>
    </h2>
    <v-layout>
      <v-flex xs12>
        <v-textarea name="input-7-1" label="Notes" v-model="job.notes" auto-grow></v-textarea>
      </v-flex>
    </v-layout>
  </div>
</template>
<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  props: ["job"],
  data: function() {
    return {
      jobData: {
        notes: ""
      }
    };
  }
});
</script>
<style scoped>
</style>
