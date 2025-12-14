<template>
  <div class="fixtures-view">
    <h1 class="page-title">Generated Fixtures</h1>

    <div v-if="loading" class="loading">Loading fixtures...</div>

    <div v-else-if="!hasFixtures" class="card text-center">
      <p>No fixtures generated yet.</p>
      <router-link to="/" class="btn btn-primary">Go to Teams</router-link>
    </div>

    <div v-else>
      <div class="grid grid-3">
        <div v-for="week in totalWeeks" :key="week" class="week-card">
          <div class="week-header">Week {{ week }}</div>
          <div v-for="match in getMatchesByWeek(week)" :key="match.id" class="match-item">
            <span class="team-name">{{ match.homeTeam.name }}</span>
            <span class="vs">-</span>
            <span class="team-name">{{ match.awayTeam.name }}</span>
          </div>
        </div>
      </div>

      <div class="btn-group">
        <button class="btn btn-primary" @click="startSimulation">Start Simulation</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getFixtures } from '../api'

const router = useRouter()
const fixtures = ref([])
const loading = ref(true)

onMounted(async () => {
  await loadFixtures()
})

const loadFixtures = async () => {
  loading.value = true
  try {
    const response = await getFixtures()
    fixtures.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load fixtures:', error)
    fixtures.value = []
  } finally {
    loading.value = false
  }
}

const hasFixtures = computed(() => fixtures.value.length > 0)

const totalWeeks = computed(() => {
  if (!fixtures.value.length) return 0
  return Math.max(...fixtures.value.map(m => m.week))
})

const getMatchesByWeek = week => {
  return fixtures.value.filter(m => m.week === week)
}

const startSimulation = () => {
  router.push('/simulation')
}
</script>

<style scoped>
.fixtures-view {
  max-width: 1200px;
  margin: 0 auto;
}
</style>
