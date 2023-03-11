<template>
  <div>
    <h2
      class="font-weight-light text-center"
      v-show="contractor.id || formMode === 'edit'"
      justify-space-around
    >
      Contractor
      <v-tooltip right>
        <template v-slot:activator="{ on }">
          <v-icon small color="text--secondary" style="vertical-align: middle" v-on="on">info</v-icon>
        </template>
        <span>
          A contractor is the company who requests the formula.
          <br />A contractor might make the request for a formula at many job sites.
        </span>
      </v-tooltip>
    </h2>
    <div v-show="contractor.id && formMode === 'view'">
      <!-- contractor info -->
      <div>
        <v-layout row justify-center class="text-xs-center">
          <!-- view mode -->
          <v-flex>
            <h4 class="display-1 font-weight-light text-capitalize">{{ contractor.company }}</h4>
          </v-flex>
        </v-layout>
        <v-layout row justify-center class="text-xs-center">
          <v-flex>
            <h6 class="subheading font-weight-light text-capitalize">{{ contractor.contact.name }}</h6>
          </v-flex>
          <v-flex>
            <h6 class="subheading font-weight-light">{{ contractor.contact.email }}</h6>
          </v-flex>
          <v-flex>
            <h6 class="subheading font-weight-light">{{ contractor.contact.phone }}</h6>
          </v-flex>
        </v-layout>
      </div>
    </div>
    <div v-show="formMode === 'edit'">
      <!-- edit mode -->
      <v-layout justify-center>
        <v-flex xs12 sm12 md12>
          <v-autocomplete
            browser-autocomplete="off"
            :items="contractorDataToList"
            item-text="company"
            item-value="id"
            v-model="contractor.id"
            label="Add Existing Contractor"
            @change="loadContractor"
          ></v-autocomplete>
        </v-flex>
      </v-layout>
      <v-layout justify-center>
        <v-flex xs12 sm6 md6>
          <v-btn block flat color="error darken-1" v-on:click="deleteContractor">Remove Contractor</v-btn>
        </v-flex>
      </v-layout>
    </div>
    <br v-show="contractor.id || formMode === 'edit'" />
    <v-divider v-show="contractor.id || formMode === 'edit'"></v-divider>
    <br v-show="contractor.id || formMode === 'edit'" />
  </div>
</template>

<script lang="ts">
import Vue from "vue";

import { CreateJobRequest } from "../../basecoat_transport_pb";
import { Contact, Contractor } from "../../basecoat_message_pb";

export default Vue.extend({
  props: ["formMode", "contractor"],
  data: function() {
    return {
      nameRules: [
        function(v: string) {
          if (!!v) {
            return true;
          }
          return "Contractor Name is required";
        }
      ]
    };
  },
  computed: {
    contractorDataToList: function() {
      let contractorDataMap: { [key: string]: Contractor } = this.$store.state
        .contractorData;

      interface modifiedContractor {
        id: string;
        company: string;
      }

      let contractorDataList = [];

      for (const [key, value] of Object.entries(contractorDataMap)) {
        let modifiedContractor: modifiedContractor = {
          id: key,
          company: value.getCompany()
        };
        contractorDataList.push(modifiedContractor);
      }

      return contractorDataList;
    }
  },
  methods: {
    deleteContractor: function() {
      let emptyContractor = {
        id: "",
        company: "",
        contact: {
          name: "",
          email: "",
          phone: ""
        } as Contact.AsObject
      } as Contractor.AsObject;
      this.$emit("update:contractor", emptyContractor);
    },
    loadContractor: function(contractorID: string) {
      if (contractorID) {
        let contractorDataMap: { [key: string]: Contractor } = this.$store.state
          .contractorData;

        let currentContractor = contractorDataMap[contractorID];
        this.$emit("update:contractor", currentContractor.toObject());
      }
    }
  }
});
</script>

<style scoped>
h6 {
  color: #2e3131;
}
</style>


