<template>
  <div>
    <h2 class="font-weight-light text-center" justify-space-around>
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
    <v-layout wrap v-if="contractor.company !== ''">
      <v-flex xs12 sm12 md12>
        <v-text-field
          label="Contractor Name"
          :rules="nameRules"
          :disabled="true"
          v-model="contractor.company"
          required
        ></v-text-field>
      </v-flex>
      <v-flex xs12 sm4 md4>
        <v-text-field label="Contact Name" :disabled="true" v-model="contractor.contact.name"></v-text-field>
      </v-flex>
      <v-flex xs12 sm4 md4>
        <v-text-field label="Contractor Email" :disabled="true" v-model="contractor.contact.email"></v-text-field>
      </v-flex>
      <v-flex xs12 sm4 md4>
        <v-text-field label="Contractor Phone" :disabled="true" v-model="contractor.contact.phone"></v-text-field>
      </v-flex>
    </v-layout>
    <v-layout justify-center v-if="contractor.company === ''">
      <v-flex xs12 sm12 md12>
        <v-autocomplete
          browser-autocomplete="off"
          :items="contractorDataToList"
          item-text="company"
          item-value="id"
          label="Add Existing Contractor"
          @change="loadContractor"
          clearable
        ></v-autocomplete>
      </v-flex>
    </v-layout>
    <v-layout justify-center v-if="contractor.company !== ''">
      <v-flex xs12 sm6 md6>
        <v-btn block flat color="error darken-1" v-on:click="deleteContractor">Remove Contractor</v-btn>
      </v-flex>
    </v-layout>
  </div>
</template>

<script lang="ts">
import Vue from "vue";

import { CreateJobRequest } from "../../basecoat_transport_pb";
import { Contact, Contractor } from "../../basecoat_message_pb";

export default Vue.extend({
  props: ["contractor"],
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

<style scoped></style>


