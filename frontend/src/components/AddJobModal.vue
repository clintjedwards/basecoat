<template>
  <v-layout row justify-center>
    <v-dialog v-model="$store.state.displayAddJobModal" max-width="600px" persistent>
      <v-form ref="addJobForm" lazy-validation>
        <v-card>
          <v-card-title>
            <span class="headline">
              <slot name="heading"></slot>
            </span>
          </v-card-title>
          <v-card-text>
            <!-- Job Info -->
            <v-container grid-list-md>
              <v-spacer>General Company Information</v-spacer>
              <v-layout wrap>
                <v-flex xs12 sm12 md12>
                  <v-text-field
                    label="Company Name"
                    :rules="nameRules"
                    v-model="jobData.name"
                    required
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <v-layout wrap>
                <v-flex xs12 sm6 md6>
                  <v-text-field label="Contact Name" v-model="jobData.contact_name"></v-text-field>
                </v-flex>
                <v-flex xs12 sm6 md6>
                  <v-text-field
                    label="Contact Info"
                    hint="This can be an email address, phone number, etc"
                    v-model="jobData.contact_info"
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <v-divider></v-divider>
              <br>

              <!-- Address -->
              <v-spacer>Address Information</v-spacer>
              <v-layout wrap>
                <v-flex xs12 sm12 md12>
                  <v-text-field label="Street" hint="Street and number" v-model="jobData.street"></v-text-field>
                </v-flex>
                <v-flex xs12 sm12 md12>
                  <v-text-field
                    label="Street 2"
                    hint="Apartment number, suite, unit, building, floor, etc."
                    v-model="jobData.street2"
                  ></v-text-field>
                </v-flex>
                <v-flex xs12 sm4 md4>
                  <v-text-field label="City" v-model="jobData.city"></v-text-field>
                </v-flex>
                <v-flex xs12 sm4 md4>
                  <v-select :items="states" label="State" v-model="jobData.state"></v-select>
                </v-flex>
                <v-flex xs12 sm4 md4>
                  <v-text-field label="Zipcode" v-model="jobData.zipcode"></v-text-field>
                </v-flex>
              </v-layout>
              <v-divider></v-divider>
              <br>

              <!-- Notes -->
              <v-spacer>Miscellaneous Information</v-spacer>
              <v-layout>
                <v-flex xs12>
                  <v-textarea name="input-7-1" label="Notes" v-model="jobData.notes" auto-grow></v-textarea>
                </v-flex>
              </v-layout>
            </v-container>
          </v-card-text>

          <!-- Buttons -->
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="blue darken-1" flat v-on:click="$store.commit('hideAddJobModal')">Close</v-btn>
            <v-btn color="blue darken-1" flat v-on:click="handleFormSave()">Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script>
export default {
  data: function() {
    return {
      jobData: {
        name: "",
        contact_name: "",
        contact_info: "",
        street: "",
        street2: "",
        city: "",
        state: "",
        zipcode: "",
        notes: ""
      },
      states: [
        "Alabama",
        "Alaska",
        "American Samoa",
        "Arizona",
        "Arkansas",
        "California",
        "Colorado",
        "Connecticut",
        "Delaware",
        "District of Columbia",
        "Federated States of Micronesia",
        "Florida",
        "Georgia",
        "Guam",
        "Hawaii",
        "Idaho",
        "Illinois",
        "Indiana",
        "Iowa",
        "Kansas",
        "Kentucky",
        "Louisiana",
        "Maine",
        "Marshall Islands",
        "Maryland",
        "Massachusetts",
        "Michigan",
        "Minnesota",
        "Mississippi",
        "Missouri",
        "Montana",
        "Nebraska",
        "Nevada",
        "New Hampshire",
        "New Jersey",
        "New Mexico",
        "New York",
        "North Carolina",
        "North Dakota",
        "Northern Mariana Islands",
        "Ohio",
        "Oklahoma",
        "Oregon",
        "Palau",
        "Pennsylvania",
        "Puerto Rico",
        "Rhode Island",
        "South Carolina",
        "South Dakota",
        "Tennessee",
        "Texas",
        "Utah",
        "Vermont",
        "Virgin Island",
        "Virginia",
        "Washington",
        "West Virginia",
        "Wisconsin",
        "Wyoming"
      ],
      nameRules: [v => !!v || "Company Name is required"]
    };
  },
  methods: {
    clearForm: function() {
      this.$refs.addJobForm.reset();
    },
    handleFormSave: function() {
      if (this.$refs.addJobForm.validate()) {
        this.$emit("submit-add-job-form", this.jobData);
      }
    }
  }
};
</script>
