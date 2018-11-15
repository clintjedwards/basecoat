
<template>
  <v-data-table
    :headers="headers"
    :items="this.$store.state.jobData"
    :search="$store.state.jobTableSearchTerm"
    hide-actions
  >
    <template v-slot:items="props">
      <tr
        style="cursor: pointer;"
        :ripple="{ center: true }"
        @click="$store.commit('showManageJobsModal', props.index)"
      >
        <td class="text-capitalize">{{ props.item.name }}</td>
        <td class="text-capitalize">
          {{ props.item.contact_name }}
          <template v-if="props.item.contact_name !== ''">|</template>
          {{ props.item.contact_info }}
        </td>
        <td>{{ props.item.street }} {{ props.item.city }}, {{ props.item.state }} {{ props.item.zipcode }}</td>
      </tr>
    </template>
  </v-data-table>
</template>

<script>
export default {
  data: function() {
    return {
      headers: [
        {
          text: "Name",
          align: "left",
          value: "name"
        },
        {
          text: "Contact Info",
          value: "contact_info"
        },
        {
          text: "Address",
          value: "address"
        }
      ]
    };
  }
};
</script>

<style scoped>
</style>
