
<template>
  <v-data-table
    :headers="headers"
    :items="stringifyFormulaData"
    :search="$store.state.formulaTableSearchTerm"
    hide-actions
  >
    <template v-slot:items="props">
      <tr
        style="cursor: pointer;"
        :ripple="{ center: true }"
        v-on:click="$store.commit('showManageFormulaModal', props.index)"
      >
        <td class="text-capitalize">{{ props.item.name }}</td>
        <td>{{ props.item.number }}</td>
        <td class="text-capitalize">{{ props.item.base }}</td>
        <td>{{ props.item.colorants }}</td>
        <td>{{ props.item.created }}</td>
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
          text: "Number",
          value: "number",
          sortable: false
        },
        {
          text: "Base",
          value: "base"
        },
        {
          text: "Colorants",
          value: "colorants"
        },
        {
          text: "Created",
          value: "created"
        }
      ]
    };
  },
  computed: {
    // This makes it so that the formula table is sortable.
    // the formula table sorts based on the data structure that
    // you pass it. So you have to pass it a data structure with
    // correct types in order of it to sort properly
    stringifyFormulaData: function() {
      let modifiedFormulaList = JSON.parse(
        JSON.stringify(this.$store.state.formulaData)
      );
      let formula = "";

      for (formula of modifiedFormulaList) {
        if (formula.base == null) {
          formula.base = { "": "" };
        }

        if (formula.colorants == null) {
          formula.colorants = {};
        }
        formula.base = String(Object.keys(formula.base)[0]);
        formula.colorants = Object.keys(formula.colorants).length;
        formula.created = parseInt(formula.created);
      }

      return modifiedFormulaList;
    }
  }
};
</script>

<style scoped>
</style>
