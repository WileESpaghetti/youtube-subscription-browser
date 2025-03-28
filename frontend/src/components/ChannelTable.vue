<script setup lang="ts">
import { ref, shallowRef, watch } from 'vue'
import {AgGridVue, GridApi} from "ag-grid-vue3";
import vue from "@vitejs/plugin-vue";
import { AllCommunityModule, ModuleRegistry } from 'ag-grid-community';

// Register all Community features
ModuleRegistry.registerModules([AllCommunityModule]);

const props = defineProps<{
  channels
}>();

const columnDefs = ref([]);
const rowData = ref([]);
const gridApi = shallowRef<GridApi | null>(null);
const onGridReady = (params: GridReadyEvent) => {
  gridApi.value = params.api;
};
watch( ()=> props.channels, (newVal, oldVal)=> {

  if (typeof newVal === 'string') {
    newVal = JSON.parse(newVal);
  }

  console.log('watcher called');
  console.log('new: %o, old: %o', newVal, oldVal);

columnDefs.value = newVal.items && newVal.items.length ? Object.keys(newVal.items[0]).map(k => {return {field: k}}) : [];
rowData.value = newVal.items || [];
// if (gridApi.value!) {
//   gridApi.value!.redrawCells();
//   gridApi.value!.redrawRows();
// }
console.log(columnDefs.value);
console.log(rowData.value);
console.log('end watcher');

}, {immediate:true, deep: true});
</script>

<template>
  <div style="width:100%;">
    <ag-grid-vue
      class="ag-theme-alpine"
      :columnDefs="columnDefs"
       style="height: 500px; width: 100%"
      :rowData="rowData"
      :paginationAutoPageSize="true"
      :pagination="true"
    ></ag-grid-vue>
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
