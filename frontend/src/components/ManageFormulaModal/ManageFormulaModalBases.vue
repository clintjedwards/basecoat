<template>
  <div>
    <v-spacer v-show="formula.basesList.length === 0">No Bases Listed</v-spacer>
    <h2
      v-show="formula.basesList.length > 0"
      class="font-weight-light text-center"
      justify-space-around
    >
      Base List
      <v-tooltip right>
        <template v-slot:activator="{ on }">
          <v-icon small color="text--secondary" style="vertical-align: middle" v-on="on">info</v-icon>
        </template>
        <span>A base is a paint medium specifically manufactured for mixing colors</span>
      </v-tooltip>
    </h2>
    <v-layout row wrap v-for="(base, index) in formula.basesList" v-bind:key="`base-${index}`">
      <v-list v-show="formMode === 'view'" style="width:100%;">
        <v-list-tile>
          <v-list-tile-avatar>
            <v-icon>invert_colors</v-icon>
          </v-list-tile-avatar>
          <v-list-tile-content>
            <v-list-tile-title v-text="base.name"></v-list-tile-title>
          </v-list-tile-content>
          <v-list-tile-avatar>
            <template>{{ base.amount }}</template>
          </v-list-tile-avatar>
        </v-list-tile>
        <v-divider></v-divider>
      </v-list>

      <v-flex xs10 sm9 v-show="formMode === 'edit'">
        <v-text-field label="Base Name" v-model="base.name"></v-text-field>
      </v-flex>
      <v-flex xs2 sm3 v-show="formMode === 'edit'">
        <v-text-field
          label="Size"
          v-model="base.amount"
          append-outer-icon="delete"
          v-on:click:append-outer="removeBaseField(index)"
        ></v-text-field>
      </v-flex>
    </v-layout>
    <div v-show="formMode === 'edit'">
      <v-btn flat color="primary" v-on:click="addBaseField">Add Base</v-btn>
    </div>
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
    return {};
  },
  methods: {
    addBaseField: function() {
      this.formula.basesList.push({
        type: "",
        name: "",
        amount: ""
      });
    },
    removeBaseField: function(index: number) {
      this.formula.basesList.splice(index, 1);
    }
  }
});
</script>
<style scoped>
</style>
