<template>
  <div>
    <!-- general job site info -->
    <v-layout row justify-center class="text-xs-center">
      <!-- view mode -->
      <v-flex v-show="formMode === 'view'">
        <h2 class="display-3 font-weight-light text-capitalize">{{ job.name }}</h2>
      </v-flex>
      <!-- edit mode -->
      <v-flex v-show="formMode === 'edit'">
        <v-text-field label="Job Site Name" :rules="nameRules" v-model="job.name" required></v-text-field>
      </v-flex>
    </v-layout>
    <br />
    <!-- contact info -->
    <!-- view mode -->
    <v-layout row justify-center class="text-xs-center" v-show="job.contact && formMode === 'view'">
      <v-flex>
        <h6 class="subheading font-weight-light text-capitalize">{{ job.contact.name }}</h6>
      </v-flex>
      <v-flex>
        <h6 class="subheading font-weight-light">{{ job.contact.email }}</h6>
      </v-flex>
      <v-flex>
        <h6 class="subheading font-weight-light">{{ job.contact.phone }}</h6>
      </v-flex>
    </v-layout>
    <!-- edit mode -->
    <v-layout row justify-center class="text-xs-center" v-show="job.contact && formMode === 'edit'">
      <v-flex>
        <v-text-field label="Contact Name" v-model="job.contact.name"></v-text-field>
      </v-flex>
      <v-flex>
        <v-text-field label="Contact Email" v-model="job.contact.email"></v-text-field>
      </v-flex>
      <v-flex>
        <v-text-field label="Contact Phone" v-model="job.contact.phone"></v-text-field>
      </v-flex>
    </v-layout>
  </div>
</template>
<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  props: ["formMode", "job", "contact"],
  data: function() {
    return {
      nameRules: [
        function(v: string) {
          if (!!v) {
            return true;
          }
          return "Job site name is required";
        }
      ]
    };
  }
});
</script>
<style scoped>
</style>
