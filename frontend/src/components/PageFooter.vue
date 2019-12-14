<template>
  <v-footer absolute padless color="#424242">
    <div class="page-footer-text font-weight-light">
      Version {{ $store.state.appInfo.version }}
      {{ $store.state.appInfo.commit }}
    </div>
    <template v-if="$store.state.appInfo.debug_enabled">
      <v-spacer></v-spacer>
      <div class="page-footer-text font-weight-light">Debug Enabled</div>
    </template>
  </v-footer>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  data: function() {
    return {};
  }
});
</script>

<style scoped>
.page-footer-text {
  color: #e0e0e0;
  padding-left: 5px;
  padding-right: 5px;
}
</style>
