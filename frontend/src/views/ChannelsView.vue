<template>
  <div class="about">
    <h1>This is a channels page</h1>
    <channel-table :channels="post || []"></channel-table>
    {{post}}
  </div>
</template>
<script setup>
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getChannels } from '../api'
import ChannelTable from "@/components/ChannelTable.vue";

const route = useRoute()

const loading = ref(false)
const post = ref(null)
const error = ref(null)

// watch the params of the route to fetch the data again
watch(() => route.params.id, fetchData, { immediate: true })

async function fetchData() {
  error.value = post.value = null
  loading.value = true

  try {
    // replace `getPost` with your data fetching util / API wrapper
    post.value = await getChannels()
  } catch (err) {
    error.value = err.toString()
  } finally {
    loading.value = false
  }
}
</script>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
