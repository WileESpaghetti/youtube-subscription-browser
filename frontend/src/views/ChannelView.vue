<template>
  <div class="about">
    <template v-if="isFetching || isFetchingStats">
      <h1 >Loading data...</h1>
    </template>

    <template v-else>
      <h1>{{channel?.title}}</h1>
      <p>{{channel?.description}}</p>

      <h2>Details</h2>
      <ul v-if="statData">
        <li>Total Videos: {{channel?.video_count}}</li>
        <li>{{statData?.has_archive ? `Archived: ${statData?.total_videos_archived}` : 'Not Archived' }}</li>
        <li>Last Upload Date: {{lastUploadDate}}</li>
      </ul>

      <h2>Upload statistics</h2>
      <p>All Time</p>
      <ul v-if="statData">
        <li>Upload Frequency: </li>
      </ul>

      <h2>Upload statistics</h2>
      <p>Last 12 Months</p>
      <ul v-if="statData">
        <li>Upload Frequency: </li>
      </ul>

      <VideoChart :channelId="Number(route.params.id)"></VideoChart>
    </template>
  </div>
</template>
<script setup lang="ts">
import {computed, onUnmounted, ref} from 'vue';
import {useFetch} from "@vueuse/core";
import {useRoute} from 'vue-router'
import VideoChart from "@/components/VideoChart.vue";

const route = useRoute();
const channelId = route.params.id;

const { data, error, isFetching, abort } = useFetch(`/api/channels/${channelId}`).json();
const { data: statData, isFetching: isFetchingStats } = useFetch(`/api/channels/${channelId}/video_stats`).json();

const channel = computed(() => {
  if (!data.value) {
    return null;
  }

  return data.value.item;
});

const lastUploadDate = computed(() => {
  if (!statData?.value?.latest_video_upload_date) {
    return null;
  }

  return new Date(statData.value.latest_video_upload_date * 1000).toLocaleDateString()
});

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
