<template>
  <v-layout row justify-center>
    <v-dialog v-model="$store.state.displayLoginModal" max-width="600px" persistent>
      <v-form ref="loginForm" lazy-validation>
        <v-card>
          <v-card-text>
            <!-- Banner -->
            <v-layout justify-center wrap>
              <img class="icon center" src="images/paintbrush.svg">
              <v-flex xs12 sm12 md12>
                <h2 class="display-3 font-weight-light center">Basecoat</h2>
              </v-flex>
            </v-layout>
            <v-spacer></v-spacer>
            <br>
            <br>

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
                  :append-icon="showPass ? 'visibility' : 'visibility_off'"
                  :type="showPass ? 'text' : 'password'"
                  :rules="passwordRules"
                  v-model="loginInfo.password"
                  outline
                  required
                  @click:append="showPass = !showPass"
                ></v-text-field>
              </v-flex>
            </v-layout>
          </v-card-text>
          <v-layout wrap justify-center>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
                color="blue darken-1"
                :loading="$store.state.loginIsLoading"
                :disabled="$store.state.loginIsLoading"
                @click="handleLogin()"
                type="submit"
              >Login</v-btn>
            </v-card-actions>
          </v-layout>
          <br>
        </v-card>
      </v-form>
    </v-dialog>
  </v-layout>
</template>

<script>
export default {
  data: function() {
    return {
      showPass: false,
      loginInfo: {
        username: "",
        password: ""
      },
      nameRules: [v => !!v || "Username is required"],
      passwordRules: [v => !!v || "Password is required"]
    };
  },
  methods: {
    handleLogin: function() {
      if (this.$refs.loginForm.validate()) {
        let loginData = {};
        loginData.username = this.loginInfo.username;
        loginData.token = btoa(
          this.loginInfo.username + ":" + this.loginInfo.password
        );
        this.$emit("validate-login", loginData);
      }
    }
  }
};
</script>

<style scoped>
.center {
  text-align: center;
}
</style>
