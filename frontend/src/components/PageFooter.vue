<template>
  <v-footer absolute padless color="#424242">
    <div class="page-footer-text font-weight-light">
      Version v{{ $store.state.appInfo.semver }} |
      {{ humanizedBuildTime }} ({{ humanizedRelativeBuildTime }}) |
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
import * as moment from "moment";

export default Vue.extend({
  data: function() {
    return {};
  },
  computed: {
    humanizedBuildTime: function() {
      let build_time = moment(
        moment.unix(this.$store.state.appInfo.build_time)
      ).format("L");

      return build_time;
    },
    humanizedRelativeBuildTime: function() {
      let build_time = moment(
        moment.unix(this.$store.state.appInfo.build_time)
      ).fromNow();

      return build_time;
    }
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
