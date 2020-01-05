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
    <!-- general contractor -->
    <v-layout row justify-center class="text-xs-center">
      <!-- view mode -->
      <v-flex v-show="formMode === 'view'">
        <h2 class="display-3 font-weight-light text-capitalize">{{ contractor.company }}</h2>
      </v-flex>
      <!-- edit mode -->
      <v-flex xs12 v-show="formMode === 'edit'">
        <v-text-field
          label="Contractor Name"
          :rules="nameRules"
          v-model="contractor.company"
          required
        ></v-text-field>
      </v-flex>
    </v-layout>
    <br />
    <!-- contact info -->
    <!-- view mode -->
    <v-layout
      row
      justify-center
      class="text-xs-center"
      v-show="contractor.contact && formMode === 'view'"
    >
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
    <!-- edit mode -->
    <v-layout
      row
      justify-center
      class="text-xs-center"
      v-show="contractor.contact && formMode === 'edit'"
    >
      <v-flex>
        <v-text-field label="Contact Name" v-model="contractor.contact.name"></v-text-field>
      </v-flex>
      <v-flex>
        <v-text-field label="Contact Email" v-model="contractor.contact.email"></v-text-field>
      </v-flex>
      <v-flex>
        <v-text-field label="Contact Phone" v-model="contractor.contact.phone"></v-text-field>
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
  props: ["formMode", "contractor"],
  data: function() {
    return {
      nameRules: [
        function(v: string) {
          if (!!v) {
            return true;
          }
          return "Company name is required";
        }
      ]
    };
  }
});
</script>
<style scoped>
</style>
