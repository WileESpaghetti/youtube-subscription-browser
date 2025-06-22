<template>
  <div class="about">
    <template v-if="isFetching">
      <h1 >Loading data...</h1>
    </template>
    <template v-else>
      <h1>{{data?.title}}</h1>
      <p>{{data?.description}}</p>
    </template>
  </div>
</template>
<script setup lang="ts">
import {ref} from 'vue';
import ChannelTable from "@/components/ChannelTable.vue";
import {useFetch} from "@vueuse/core";
import {useRoute} from 'vue-router'

const route = useRoute();
const { data, error, isFetching, execute, abort } = useFetch(`/api/channels/${route.params.id}`).json();
const mydata = ref(data);
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
