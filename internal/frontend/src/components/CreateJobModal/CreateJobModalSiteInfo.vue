<template>
  <div>
    <h2 class="font-weight-light" justify-space-around>
      Job Site
      <v-tooltip right>
        <template v-slot:activator="{ on }">
          <v-icon small color="text--secondary" style="vertical-align: middle" v-on="on">info</v-icon>
        </template>
        <span>
          A job site is the place the produced formula has been sent to.
          <br />This might be the location that the contractor who requested a formula specifies
        </span>
      </v-tooltip>
    </h2>
    <v-layout wrap>
      <v-flex xs12 sm12 md12>
        <v-text-field label="Job Site Name" :rules="nameRules" v-model="job.name" required></v-text-field>
      </v-flex>
    </v-layout>
    <v-layout wrap>
      <v-flex xs12 sm4 md4>
        <v-text-field label="Contact Name" v-model="job.contact.name"></v-text-field>
      </v-flex>
      <v-flex xs12 sm4 md4>
        <v-text-field label="Contact Email" v-model="job.contact.email"></v-text-field>
      </v-flex>
      <v-flex xs12 sm4 md4>
        <v-text-field label="Contact Phone" v-model="job.contact.phone"></v-text-field>
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
      nameRules: [
        function(v: string) {
          if (!!v) {
            return true;
          }
          return "Job Site Name is required";
        }
      ]
    };
  }
});
</script>
<style scoped>
</style>
