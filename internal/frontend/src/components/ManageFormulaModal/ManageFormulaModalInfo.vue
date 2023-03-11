<template>
  <div>
    <h2 class="font-weight-light text-center" justify-space-around>
      Formula
      <v-tooltip right>
        <template v-slot:activator="{ on }">
          <v-icon small color="text--secondary" style="vertical-align: middle" v-on="on">info</v-icon>
        </template>
        <span>A formula is a set of colorants and bases that are used to make a certain color</span>
      </v-tooltip>
    </h2>
    <v-layout row justify-center class="text-xs-center">
      <v-flex v-show="formMode === 'view'">
        <h2 class="display-3 font-weight-light text-capitalize">{{ formula.name }}</h2>
      </v-flex>
      <v-flex v-show="formMode === 'edit'">
        <v-text-field label="Formula Name" :rules="nameRules" v-model="formula.name" required></v-text-field>
      </v-flex>
    </v-layout>
    <v-layout row justify-center class="text-xs-center">
      <v-flex v-show="formMode === 'view'">
        <h4 class="display-1 font-weight-light text-capitalize">{{ formula.number }}</h4>
      </v-flex>
      <v-flex v-show="formMode === 'edit'">
        <v-text-field
          label="Formula Number"
          hint="Custom number used to reference formula"
          v-model="formula.number"
          required
        ></v-text-field>
      </v-flex>
    </v-layout>
    <br />
    <v-divider></v-divider>
    <br />
  </div>
</template>
<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  props: ["formMode", "formula"],
  data: function() {
    return {
      nameRules: [
        function(v: string) {
          if (!!v) {
            return true;
          }
          return "Formula Name is required";
        }
      ]
    };
  }
});
</script>
<style scoped>
</style>
