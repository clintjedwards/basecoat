
<template>
  <v-data-table
    :headers="headers"
    :items="formulaDataToList"
    :search="$store.state.formulaTableSearchTerm"
    hide-actions
  >
    <template v-slot:items="props">
      <tr
        style="cursor: pointer;"
        :ripple="{ center: true }"
        v-on:click="navigateToFormula(props.item.id)"
      >
        <td class="text-capitalize">{{ props.item.name }}</td>
        <td>{{ props.item.number }}</td>
        <td class="text-capitalize">{{ props.item.base }}</td>
        <td>{{ props.item.colorants }}</td>
        <td class="text-capitalize">{{ props.item.created }}</td>
      </tr>
    </template>
  </v-data-table>
</template>

<script lang="ts">
import Vue from "vue";
import { Formula } from "../basecoat_pb";
import * as moment from "moment";

interface modifiedFormula {
  id: string;
  name: string;
  number: string;
  base: string;
  colorants: number;
  created: string;
}

export default Vue.extend({
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
    formulaDataToList(): modifiedFormula[] {
      let formulaDataMap: { [key: string]: Formula } = this.$store.state
        .formulaData;
      let formulaDataList: Formula[] = [];

      for (const [key, value] of Object.entries(formulaDataMap)) {
        formulaDataList.push(value);
      }

      let modifiedFormulaList: modifiedFormula[] = [];
      let formula: Formula;

      for (formula of formulaDataList) {
        let modifiedFormula: modifiedFormula = {
          id: "",
          name: "",
          number: "",
          base: "",
          colorants: 0,
          created: ""
        };

        modifiedFormula.id = formula.getId();
        modifiedFormula.colorants = formula.getColorantsList().length;
        modifiedFormula.name = formula.getName();
        modifiedFormula.number = formula.getNumber();
        modifiedFormula.created = moment(
          moment.unix(formula.getCreated())
        ).format("L");
        modifiedFormula.base = "None";

        if (formula.getBasesList().length != 0) {
          modifiedFormula.base = formula.getBasesList()[0].getName();
        }
        modifiedFormulaList.push(modifiedFormula);
      }

      return modifiedFormulaList;
    }
  },
  methods: {
    navigateToFormula: function(formulaID: string) {
      this.$router.push("/formulas/" + formulaID);
    }
  }
});
</script>

<style scoped>
</style>
