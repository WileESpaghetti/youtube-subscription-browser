<script setup lang="ts">
import { VisXYContainer, VisGroupedBar, VisAxis } from '@unovis/vue';
import {useFetch} from "@vueuse/core";
import {BarChart} from 'vue-chart-3';
import {Chart, registerables} from "chart.js";
import {computed} from "vue";
Chart.register(...registerables);

const props = defineProps<{
  channelId: number
}>()

// TODO make date range configuratble to allow people to zoom in/out
const months = [...Array(12).keys()].map(key => new Date(0, key).toLocaleString('en', { month: 'long' }));
const today = new Date();
const thisYear = today.getUTCFullYear();
const thisMonth = today.getUTCMonth();

const last12Months = months.slice(thisMonth + 1, months.length).concat(months.slice(0, thisMonth + 1));

const from = Date.UTC(thisYear - 1, thisMonth + 1  ,1) / 1000;
console.log('from: ', from);
const { data, error, isFetching, execute, abort } = useFetch(`/api/videos?channel_id=${props.channelId}&from=${from}`).json();

const chartData = computed(() => {
  return {
    labels: last12Months,
    datasets: [{
      label: 'Videos by Month', data: data?.value?.items?.reduce((a, v) => {
        console.log('hello');
        const month = new Date(v.timestamp * 1000).toLocaleString('default', {month: 'long'});
        if (!a[month]) {
          a[month] = 1;
        } else {
          a[month]++;
        }
        return a;
      }, {}), backgroundColor: ['#77CEFF', '#0079AF', '#123E6B', '#97B0C4', '#A5C8ED'],
    }]
  };
});





///////////////
</script>

<template>
  <BarChart v-if="!isFetching"  :chartData="chartData" />
</template>

<!--<template>-->
<!--  <div v-if="isFetching">-->
<!--    Getting video data...-->
<!--  </div>-->
<!--  <div v-else>-->
<!--    <VisXYContainer :data="r">-->
<!--      <VisGroupedBar :x="(d) => d.x" :y="(d) => d.map(d => (c) => c.y)" />-->
<!--    </VisXYContainer>-->
<!--  </div>-->
<!--</template>-->
