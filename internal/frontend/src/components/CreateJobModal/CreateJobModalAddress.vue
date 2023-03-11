<template>
  <div>
    <h2 class="font-weight-light" justify-space-around>Address</h2>
    <v-layout wrap>
      <v-flex xs12 sm12 md12>
        <v-text-field label="Street" hint="Street and number" v-model="job.address.street"></v-text-field>
      </v-flex>
      <v-flex xs12 sm12 md12>
        <v-text-field
          label="Street 2"
          hint="Apartment number, suite, unit, building, floor, etc."
          v-model="job.address.street2"
        ></v-text-field>
      </v-flex>
      <v-flex xs12 sm4 md4>
        <v-text-field label="City" v-model="job.address.city"></v-text-field>
      </v-flex>
      <v-flex xs12 sm4 md4>
        <v-autocomplete
          browser-autocomplete="off"
          :items="states"
          label="State"
          v-model="job.address.state"
        ></v-autocomplete>
      </v-flex>
      <v-flex xs12 sm4 md4>
        <v-text-field label="Zipcode" v-model="job.address.zipcode"></v-text-field>
      </v-flex>
    </v-layout>
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import states from "../states";

let stateList = states;

export default Vue.extend({
  props: ["job"],
  data: function() {
    return {
      states: stateList
    };
  }
});
</script>
<style scoped>
</style>
