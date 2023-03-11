<template>
  <div>
    <h2 class="font-weight-light text-center" v-show="formMode === 'edit'">Miscellaneous Information</h2>
    <h2 class="font-weight-light text-center" v-show="formMode === 'view'">Notes</h2>
    <div>
      <v-layout>
        <v-flex xs12 v-show="formMode === 'view'">
          <pre>{{ formula.notes }}</pre>
        </v-flex>
        <v-flex xs12 v-show="formMode === 'edit'">
          <v-textarea name="input-7-1" label="Notes" v-model="formula.notes" auto-grow></v-textarea>
        </v-flex>
      </v-layout>
    </div>
  </div>
</template>
<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  props: ["formMode", "formula"],
  data: function() {
    return {};
  }
});
</script>
<style scoped>
pre {
  white-space: pre-wrap;
}
</style>
