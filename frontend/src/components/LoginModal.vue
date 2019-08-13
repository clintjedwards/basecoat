<template>
  <v-layout row justify-center>
    <v-dialog v-model="showModal" max-width="600px" persistent>
      <v-form ref="loginForm" lazy-validation>
        <v-card>
          <v-card-text>
            <!-- Banner -->
            <v-layout justify-center wrap>
              <img class="icon center" src="images/paintbrush.svg" />
              <v-flex xs12 sm12 md12>
                <h2 class="display-3 font-weight-light center">Basecoat</h2>
              </v-flex>
            </v-layout>
            <v-spacer></v-spacer>
            <br />
            <br />

            <v-layout wrap justify-center>
              <v-flex xs12 sm8 md8>
                <v-text-field
                  label="Username"
                  :rules="nameRules"
                  v-model="loginInfo.username"
                  outline
                  required
                ></v-text-field>
              </v-flex>
              <v-flex xs12 sm8 md8>
                <v-text-field
                  label="Password"
                  :append-icon="showPassword ? 'visibility' : 'visibility_off'"
                  :type="showPassword ? 'text' : 'password'"
                  :rules="passwordRules"
                  v-model="loginInfo.password"
                  outline
                  required
                  @click:append="showPassword = !showPassword"
                ></v-text-field>
              </v-flex>
            </v-layout>
          </v-card-text>
          <v-layout wrap justify-center>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
                color="blue darken-1"
                :loading="this.loading"
                :disabled="this.loading"
                @click="handleLogin()"
              >Login</v-btn>
            </v-card-actions>
          </v-layout>
          <br />
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script lang="ts">
import Vue from "vue";
import BasecoatClientWrapper from "../basecoatClientWrapper";

let client: BasecoatClientWrapper;
client = new BasecoatClientWrapper();

export default Vue.extend({
  data: function() {
    return {
      showModal: true,
      showPassword: false,
      loading: false,
      loginInfo: {
        username: "",
        password: ""
      },
      nameRules: [
        function(v: string) {
          if (!!v) {
            return true;
          }
          return "Username is required";
        }
      ],
      passwordRules: [
        function(v: string) {
          if (!!v) {
            return true;
          }
          return "Password is required";
        }
      ]
    };
  },
  methods: {
    handleLogin: function() {
      if ((this.$refs.loginForm as HTMLFormElement).validate()) {
        this.loading = true;

        client
          .handleLogin(this.loginInfo)
          .then(() => {
            this.$store.commit("updateUsername", this.loginInfo.username);
            //Redirect user to original page if coming from another
            if (this.$route.query.redirect) {
              this.$router.push(this.$route.query.redirect.toString());
            } else {
              this.$router.push("/");
            }
            this.loading = false;
          })
          .catch(error => {
            console.log(error);
            this.$store.commit("showSnackBar", "Invalid Login Credentials");
            this.loading = false;
          });
      }
    }
  }
});
</script>

<style scoped>
.center {
  text-align: center;
}
</style>
