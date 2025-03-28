<template>
  <div class="about">
    <h1>This is a channels page</h1>
    <h2 v-if="isFetching">Loading data...</h2>
    <channel-table v-if="data" :channels="data || {items: []}"></channel-table>
  </div>
</template>
<script setup lang="ts">
import ChannelTable from "@/components/ChannelTable.vue";
import {useFetch} from "@vueuse/core";

const { data, error, isFetching, execute, abort } = useFetch(`/api/channels`);
</script>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
    flex-direction: column;
    justify-content: center;
  }
}
</style>
