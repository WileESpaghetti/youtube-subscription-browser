<script setup lang="ts">
import {type Ref, ref, shallowRef, watch} from 'vue'
import {AgGridVue} from "ag-grid-vue3";
import { AllCommunityModule, ModuleRegistry } from 'ag-grid-community';
import {useRouter} from "vue-router";
import type {ColDef, ColGroupDef} from 'ag-grid-community';
import {
  ClientSideRowModelModule,
  DateFilterModule,
  ExternalFilterModule,
  type GridApi,
  type GridOptions,
  type GridReadyEvent,
  type IDateFilterParams,
  type IRowNode,
  type IsExternalFilterPresentParams,
  NumberFilterModule,
  ValidationModule,
} from "ag-grid-community";

// Register all Community features
ModuleRegistry.registerModules([AllCommunityModule]);

const props = defineProps<{
  channels
}>();

const columnDefs: Ref<ColDef[]|ColGroupDef[]|null> = ref(null);
const rowData = ref([]);
const gridApi = shallowRef<GridApi | null>(null);
const onGridReady = (params: GridReadyEvent) => {
  console.log('gridReady');
  gridApi.value = params.api;
};
const router = useRouter();

watch( ()=> props.channels, (newVal, oldVal)=> {

  if (typeof newVal === 'string') {
    newVal = JSON.parse(newVal);
  }

  console.log('watcher called');
  console.log('new: %o, old: %o', newVal, oldVal);

columnDefs.value = newVal.items && newVal.items.length ? Object.keys(newVal.items[0]).map(k => {
  let opts: ColDef = {
    field: k, filter: true, floatingFilter: true
  };
  if (k === 'title') {
    opts.cellRenderer = (params) => {
      const route = {
        name: "viewChannel",
        params: { id: params.data.id }
      };

      const link = document.createElement("a");
      link.href = router.resolve(route).href;
      link.innerText = params.value;
      link.addEventListener("click", e => {
        e.preventDefault();
        router.push(route);
      });
      return link;
    }
  }
  return opts;
}) : null; // Floating filter might be higher level in the ag-grid options
rowData.value = newVal.items || [];
// if (gridApi.value!) {
//   gridApi.value!.redrawCells();
//   gridApi.value!.redrawRows();
// }
console.log(columnDefs.value);
console.log(rowData.value);
console.log('end watcher');

}, {immediate:true, deep: true});

const tagFilter = ref("")
const isExternalFilterPresent: () => boolean = () => {
  // if ageType is not everyone, then we are filtering
  return !!tagFilter.value
};

const doesExternalFilterPass: (node: IRowNode<any>) => boolean = (
  node: IRowNode<any>,
) => {
  return node.data.title.toLowerCase().includes(tagFilter.value.toLowerCase());
};

watch(tagFilter, async (oldVal, newVal) => {
    gridApi.value!.onFilterChanged();
});
</script>

<template>
  <input type="text" name="tag-filter" v-model="tagFilter">
  <div style="width:100%;">
    <ag-grid-vue
      class="ag-theme-alpine"
      :columnDefs="columnDefs"
      @grid-ready="onGridReady"
       style="height: 500px; width: 100%"
      :rowData="rowData"
      :paginationAutoPageSize="true"
      :pagination="true"
      :isExternalFilterPresent="isExternalFilterPresent"
      :doesExternalFilterPass="doesExternalFilterPass"
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
