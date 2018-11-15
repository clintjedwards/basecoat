<template>
  <v-toolbar color="#424242" height="150">
    <img class="icon" src="images/paintbrush.svg">
    <v-toolbar-title class="page-header-text-title font-weight-thin display-3">
      <slot></slot>
    </v-toolbar-title>
    <v-spacer></v-spacer>
    <h3
      class="page-header-text-username title font-weight-light text-capitalize"
    >{{ $store.state.username }}</h3>
  </v-toolbar>
</template>

<script>
export default {};
</script>

<style scoped>
.page-header-text-title {
  color: #e0e0e0;
  margin-left: 0px !important;
}

.page-header-text-username {
  color: #e0e0e0;
}
</style>
