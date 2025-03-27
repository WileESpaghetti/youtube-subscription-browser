<script setup lang="ts">
import { ref, shallowRef, watch } from 'vue'
import {AgGridVue, GridApi} from "ag-grid-vue3";
import vue from "@vitejs/plugin-vue";
import { AllCommunityModule, ModuleRegistry } from 'ag-grid-community';

// Register all Community features
ModuleRegistry.registerModules([AllCommunityModule]);

const props = defineProps<{
  channels: array
}>();

const columnDefs = ref([]);
const rowData = ref([]);
const gridApi = shallowRef<GridApi | null>(null);
const onGridReady = (params: GridReadyEvent) => {
  gridApi.value = params.api;
};
watch( ()=> props.channels, (newVal, oldVal)=> {
  console.log(newVal, oldVal);

columnDefs.value = newVal && newVal.length ? Object.keys(newVal[0]).map(k => {return {field: k}}) : [];
rowData.value = newVal || [];
// if (gridApi.value!) {
//   gridApi.value!.redrawCells();
//   gridApi.value!.redrawRows();
// }
console.log(columnDefs.value);
  console.log(rowData.value);

}, {immediate:true, deep: true});
</script>

<template>
  <div>
    <ag-grid-vue v-if="channels.length"
      class="ag-theme-alpine"
      :columnDefs="columnDefs"
       domLayout="autoHeight"
      :rowData="rowData"
    ></ag-grid-vue>
    @grid-ready="gridReady"
  </div>
</template>

<style scoped>
h1 {
  font-weight: 500;
  font-size: 2.6rem;
  position: relative;
  top: -10px;
}

h3 {
  font-size: 1.2rem;
}

.greetings h1,
.greetings h3 {
  text-align: center;
}

@media (min-width: 1024px) {
  .greetings h1,
  .greetings h3 {
    text-align: left;
  }
}
</style>
