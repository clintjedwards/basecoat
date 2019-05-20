<template>
  <v-layout row justify-center>
    <v-dialog v-model="$store.state.displayManageJobsModal" max-width="600px" persistent>
      <v-form ref="manageJobsForm" lazy-validation>
        <v-card>
          <v-card-title>
            <span>ID: {{ jobData.id }}</span>
          </v-card-title>
          <v-card-text>
            <!-- Job Info -->
            <v-container grid-list-md>
              <v-spacer>General Company Information</v-spacer>
              <v-layout wrap>
                <v-flex xs12 sm12 md12>
                  <h2
                    class="display-3 font-weight-light text-capitalize"
                    v-show="formMode === 'view'"
                  >{{ jobData.name }}</h2>
                  <v-text-field
                    label="Company Name"
                    :rules="nameRules"
                    v-model="jobData.name"
                    v-show="formMode === 'edit'"
                    required
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <br>
              <v-layout v-show="formMode === 'view'">
                <v-flex xs12 sm6 md6>
                  <h6
                    class="subheading font-weight-light text-capitalize"
                  >{{ jobData.contact.name }}</h6>
                </v-flex>
                <v-flex xs12 sm6 md6>
                  <h6 class="subheading font-weight-light pull-right">{{ jobData.contact.info }}</h6>
                </v-flex>
              </v-layout>
              <v-layout wrap v-show="formMode === 'edit'">
                <v-flex xs12 sm12 md12>
                  <v-text-field label="Contact Name" v-model="jobData.contact.name"></v-text-field>
                </v-flex>
                <v-flex xs12 sm12 md12>
                  <v-text-field
                    label="Contact Info"
                    hint="This can be an email address, phone number, etc"
                    v-model="jobData.contact.info"
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <v-divider></v-divider>
              <br>

              <!-- Address -->
              <v-spacer>Address Information</v-spacer>
              <v-layout wrap>
                <v-flex xs12 sm12 md12>
                  <h5
                    class="headline font-weight-light text-capitalize"
                    v-show="formMode === 'view'"
                  >{{ jobData.street }}</h5>
                  <v-text-field
                    label="Street"
                    v-model="jobData.street"
                    v-show="formMode === 'edit'"
                  ></v-text-field>
                </v-flex>
                <v-flex xs12 sm12 md12>
                  <h5
                    class="subheading font-weight-light text-capitalize"
                    v-show="formMode === 'view'"
                  >{{ jobData.street2 }}</h5>
                  <v-text-field
                    label="Street 2"
                    hint="Apt Num, Extra Information, etc"
                    v-model="jobData.street2"
                    v-show="formMode === 'edit'"
                  ></v-text-field>
                </v-flex>
                <v-flex xs12 sm12 md12 v-show="formMode === 'view'">
                  <h5 class="subheading font-weight-light text-capitalize">
                    {{ jobData.city }}
                    <template v-if="jobData.city !== ''">,</template>
                    {{ jobData.state }} {{ jobData.zipcode }}
                  </h5>
                </v-flex>
                <v-flex xs12 sm4 md4>
                  <v-text-field label="City" v-model="jobData.city" v-show="formMode === 'edit'"></v-text-field>
                </v-flex>
                <v-flex xs12 sm4 md4>
                  <v-select
                    :items="states"
                    label="State"
                    v-model="jobData.state"
                    v-show="formMode === 'edit'"
                  ></v-select>
                </v-flex>
                <v-flex xs12 sm4 md4>
                  <v-text-field
                    label="Zipcode"
                    v-model="jobData.zipcode"
                    v-show="formMode === 'edit'"
                  ></v-text-field>
                </v-flex>
              </v-layout>
              <v-divider></v-divider>
              <br>

              <!-- Notes -->
              <v-spacer v-show="formMode === 'edit'">Miscellaneous Information</v-spacer>
              <v-layout>
                <v-flex xs12 v-show="formMode === 'view'">
                  <v-spacer>Notes</v-spacer>
                  <br>
                  <pre>{{ jobData.notes }}</pre>
                </v-flex>
                <v-flex xs12 v-show="formMode === 'edit'">
                  <v-textarea name="input-7-1" label="Notes" v-model="jobData.notes" auto-grow></v-textarea>
                </v-flex>
              </v-layout>
            </v-container>
          </v-card-text>

          <!-- Buttons -->
          <v-card-actions>
            <v-btn
              color="error darken-1"
              flat
              v-show="formMode === 'view' && !showConfirmDelete"
              @click="showConfirmDelete = true"
            >Delete</v-btn>
            <v-btn
              color="error darken-1"
              flat
              v-show="formMode === 'view' && showConfirmDelete"
              @click="handleFormDelete()"
            >Confirm Delete</v-btn>
            <v-spacer></v-spacer>
            <v-btn
              color="blue darken-1"
              flat
              @click="$store.commit('hideManageJobsModal'); setFormModeView(); populateFormData();"
            >Close</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'view'"
              @click="setFormModeEdit()"
            >Edit</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'edit'"
              @click="setFormModeView()"
            >View</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'edit'"
              @click="populateFormData()"
            >Reset</v-btn>
            <v-btn
              color="blue darken-1"
              flat
              v-show="formMode === 'edit'"
              @click="handleFormSave()"
            >Save</v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script lang="ts">
import Vue from "vue";
import { UpdateJobRequest, Contact, Job } from "../basecoat_pb";

let contact: Contact.AsObject = {
  name: "",
  info: ""
};

let jobData: UpdateJobRequest.AsObject = {
  id: "",
  name: "",
  street: "",
  street2: "",
  city: "",
  state: "",
  zipcode: "",
  notes: "",
  formulasList: [],
  contact: contact
};

let job: Job;

export default Vue.extend({
  props: ["jobIdInView"],
  data: function() {
    return {
      formMode: "view",
      showConfirmDelete: false,
      jobData: jobData,
      jobInView: job,
      nameRules: [
        function(v: string) {
          if (!!v) {
            return true;
          }
          return "Company Name is required";
        }
      ],
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
      ]
    };
  },
  watch: {
    jobIdInView: function() {
      this.jobInView = this.$store.state.jobData[this.jobIdInView];
    },
    jobInView: function() {
      this.populateFormData();
    }
  },
  computed: {
    formulaDataToList: function() {}
  },
  methods: {
    setFormModeEdit: function() {
      this.formMode = "edit";
    },
    setFormModeView: function() {
      this.formMode = "view";
    },
    populateFormData: function() {
      let currentJob = this.jobInView;

      this.jobData.id = currentJob.getId();
      this.jobData.name = currentJob.getName();
      this.jobData.street = currentJob.getStreet();
      this.jobData.street2 = currentJob.getStreet2();
      this.jobData.city = currentJob.getCity();
      this.jobData.state = currentJob.getState();
      this.jobData.zipcode = currentJob.getZipcode();
      this.jobData.notes = currentJob.getNotes();
      this.jobData.formulasList = currentJob.getFormulasList();

      let contact: Contact.AsObject;
      contact = {
        name: "",
        info: ""
      };

      if (currentJob.getContact() != undefined) {
        let currentContact = currentJob.getContact() || new Contact();
        contact = {
          name: currentContact.getName(),
          info: currentContact.getInfo()
        };
      }

      this.jobData.contact = contact;
    },
    handleFormSave: function() {
      if ((this.$refs.manageJobsForm as HTMLFormElement).validate()) {
        this.$emit("submit-manage-jobs-form", this.jobData);
      }
    },
    handleFormDelete: function() {
      this.$emit("delete-job", this.jobData.id);
      this.showConfirmDelete = false;
    }
  }
});
</script>

<style scoped>
h2 {
  text-align: center;
}

h5 {
  text-align: center;
}

h6 {
  color: #2e3131;
}

pre {
  white-space: pre-wrap;
}

.pull-right {
  text-align: right;
}
</style>
