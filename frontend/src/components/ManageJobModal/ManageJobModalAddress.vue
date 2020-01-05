<template>
  <div v-if="address.street != '' || formMode === 'edit'">
    <h2 class="font-weight-light text-center">Address</h2>
    <div>
      <v-layout wrap justify-center class="text-xs-center">
        <!-- street -->
        <v-flex xs12>
          <!-- view mode -->
          <h5
            class="headline font-weight-light text-capitalize"
            v-show="formMode === 'view'"
          >{{ address.street }}</h5>

          <!-- edit mode -->
          <v-text-field label="Street" v-model="address.street" v-show="formMode === 'edit'"></v-text-field>
        </v-flex>

        <!-- street 2-->
        <!-- view mode -->
        <v-flex xs12>
          <h5
            class="subheading font-weight-light text-capitalize"
            v-show="formMode === 'view'"
          >{{ address.street2 }}</h5>

          <!-- edit mode -->
          <v-text-field
            label="Street 2"
            hint="Apt Num, Extra Information, etc"
            v-model="address.street2"
            v-show="formMode === 'edit'"
          ></v-text-field>
        </v-flex>

        <!-- city, state, zipcode-->
        <!-- view mode -->
        <v-flex xs12 v-show="formMode === 'view'">
          <h5 class="subheading font-weight-light text-capitalize">
            {{ address.city }}
            <template v-if="address.city !== ''">,</template>
            {{ address.state }} {{ address.zipcode }}
          </h5>
        </v-flex>
        <!-- edit mode -->
        <v-flex xs12 sm4 md4>
          <v-text-field label="City" v-model="address.city" v-show="formMode === 'edit'"></v-text-field>
        </v-flex>
        <v-flex xs12 sm4 md4>
          <v-autocomplete
            browser-autocomplete="off"
            :items="states"
            label="State"
            v-model="address.state"
            v-show="formMode === 'edit'"
          ></v-autocomplete>
        </v-flex>
        <v-flex xs12 sm4 md4>
          <v-text-field label="Zipcode" v-model="address.zipcode" v-show="formMode === 'edit'"></v-text-field>
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
import { Address } from "../../basecoat_message_pb";
import states from "../states";

export default Vue.extend({
  props: ["formMode", "address"],
  data: function() {
    return {
      states: states
    };
  }
});
</script>
<style scoped>
h6 {
  color: #2e3131;
}
</style>
