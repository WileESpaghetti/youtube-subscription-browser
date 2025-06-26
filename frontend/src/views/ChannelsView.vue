<template>
  <div class="about">
    <h1>This is a channels page</h1>
    <h2 v-if="isFetching">Loading channels...</h2>
    <channel-table v-if="data" :channels="data || {items: []}"></channel-table>
  </div>
</template>

<script setup lang="ts">
import ChannelTable from "@/components/ChannelTable.vue";
import {useFetch} from "@vueuse/core";
import {onUnmounted} from "vue";

const { data, error, isFetching, abort } = useFetch(`/api/channels`);

onUnmounted(() => {
  abort()
})
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
