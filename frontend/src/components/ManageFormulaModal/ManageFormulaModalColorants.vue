<template>
  <div>
    <v-spacer v-show="formula.colorantsList.length === 0">No Colorants Listed</v-spacer>
    <h2
      v-show="formula.colorantsList.length > 0"
      class="font-weight-light text-center"
      justify-space-around
    >
      Colorants List
      <v-tooltip right>
        <template v-slot:activator="{ on }">
          <v-icon small color="text--secondary" style="vertical-align: middle" v-on="on">info</v-icon>
        </template>
        <span>A colorant is a concentrated dispersion of colour pigment that is used to tint a base paint</span>
      </v-tooltip>
    </h2>
    <v-list v-show="formMode === 'view' && colorantOverallTypeSet" style="width:100%;">
      <v-list-tile>
        <v-list-tile-avatar tile>
          <v-img max-height="25" max-width="30" v-bind:src="colorantTypeInfo.imageURL"></v-img>
        </v-list-tile-avatar>
        <v-list-tile-content>
          <v-list-tile-title v-text="colorantTypeInfo.userMessage"></v-list-tile-title>
        </v-list-tile-content>
      </v-list-tile>
      <v-divider></v-divider>
    </v-list>
    <v-flex xs6 sm6 v-show="formula.colorantsList.length != 0 && formMode === 'edit'">
      <v-select
        :items="colorantTypesToList"
        label="Colorant Type"
        v-model="currentColorantType"
        v-on:change="fillColorantTypes(currentColorantType)"
        clearable
      ></v-select>
    </v-flex>
    <v-layout
      row
      wrap
      v-for="(colorant, index) in formula.colorantsList"
      v-bind:key="`colorant-${index}`"
    >
      <v-list v-show="formMode === 'view'" style="width:100%;">
        <v-list-tile>
          <v-list-tile-avatar>
            <v-icon>color_lens</v-icon>
          </v-list-tile-avatar>
          <v-list-tile-content>
            <v-list-tile-sub-title v-show="currentColorantType === ''" v-text="colorant.type"></v-list-tile-sub-title>
            <v-list-tile-title v-text="colorant.name"></v-list-tile-title>
          </v-list-tile-content>
          <v-list-tile-avatar>
            <template>{{ colorant.amount }}</template>
          </v-list-tile-avatar>
        </v-list-tile>
        <v-divider></v-divider>
      </v-list>

      <v-flex xs2 sm2 v-show="formMode === 'edit'">
        <v-text-field label="Type" v-model="colorant.type"></v-text-field>
      </v-flex>
      <v-flex xs7 sm7 v-show="formMode === 'edit'">
        <v-text-field label="Colorant Name" v-model="colorant.name"></v-text-field>
      </v-flex>
      <v-flex xs3 sm3 v-show="formMode === 'edit'">
        <v-text-field
          label="Amount"
          v-model="colorant.amount"
          append-outer-icon="delete"
          v-on:click:append-outer="removeColorantField(index)"
        ></v-text-field>
      </v-flex>
    </v-layout>
    <div v-show="formMode === 'edit'">
      <v-btn flat color="primary" v-on:click="addColorantField">Add Colorant</v-btn>
    </div>
    <br />
    <v-divider></v-divider>
    <br />
  </div>
</template>
<script lang="ts">
import Vue from "vue";
import { Colorant } from "../../basecoat_message_pb";

export default Vue.extend({
  props: ["formMode", "formula"],
  data: function() {
    return {
      colorantOverallTypeSet: false,
      colorantTypeInfo: {},
      currentColorantType: ""
    };
  },
  watch: {
    formMode: function() {
      this.parseColorantListForSameType;
    }
  },
  mounted() {
    this.parseColorantListForSameType;
  },
  computed: {
    colorantTypesToList(): string[] {
      let colorantTypeMap: { [key: string]: string } = this.$store.state
        .colorantTypes;
      let colorantTypeList: string[] = [];

      for (const [key, value] of Object.entries(colorantTypeMap)) {
        let colorantType: string;
        colorantType = key;
        colorantTypeList.push(colorantType);
      }

      return colorantTypeList;
    },
    parseColorantListForSameType: function() {
      let self = this;
      if (self.formula.colorantsList.length < 1) {
        this.colorantOverallTypeSet = false;
        this.currentColorantType = "";
        return;
      }
      for (let i = 0; i < self.formula.colorantsList.length; ++i) {
        if (
          self.formula.colorantsList[i].type !=
          self.formula.colorantsList[0].type
        ) {
          this.colorantOverallTypeSet = false;
          this.currentColorantType = "";
          return;
        }
      }
      if (
        self.formula.colorantsList[0].type in this.$store.state.colorantTypes
      ) {
        this.colorantOverallTypeSet = true;
        this.colorantTypeInfo = this.$store.state.colorantTypes[
          self.formula.colorantsList[0].type
        ];
        this.currentColorantType = self.formula.colorantsList[0].type;
        return;
      }
      this.colorantOverallTypeSet = false;
      this.currentColorantType = "";
    }
  },
  methods: {
    addColorantField: function() {
      let type: string = "";
      if (this.currentColorantType != "") {
        type = this.currentColorantType;
      }
      this.formula.colorantsList.push({
        type: type,
        name: "",
        amount: ""
      });
    },
    fillColorantTypes: function(type: string) {
      this.formula.colorantsList.forEach(function(colorant: Colorant.AsObject) {
        colorant.type = type;
      });
    },
    removeColorantField: function(index: number) {
      this.formula.colorantsList.splice(index, 1);
    }
  }
});
</script>
<style scoped>
</style>
