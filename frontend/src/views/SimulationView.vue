<template>
  <div class="simulation-view">
    <h1 class="page-title">Simulation</h1>

    <div v-if="loading" class="loading">Loading...</div>

    <div v-else-if="!state.leagueState.fixturesCreated" class="card text-center">
      <p>No fixtures generated yet.</p>
      <router-link to="/" class="btn btn-primary">Go to Teams</router-link>
    </div>

    <div v-else>
      <div class="simulation-grid">
        <!-- League Table -->
        <div class="card">
          <div class="card-header">League Table</div>
          <table>
            <thead>
              <tr>
                <th>Team Name</th>
                <th class="text-center">P</th>
                <th class="text-center">W</th>
                <th class="text-center">D</th>
                <th class="text-center">L</th>
                <th class="text-center">GD</th>
                <th class="text-center">PTS</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="team in state.standings" :key="team.teamId">
                <td>{{ team.teamName }}</td>
                <td class="text-center">{{ team.played }}</td>
                <td class="text-center">{{ team.won }}</td>
                <td class="text-center">{{ team.drawn }}</td>
                <td class="text-center">{{ team.lost }}</td>
                <td class="text-center">{{ team.goalDifference }}</td>
                <td class="text-center">
                  <strong>{{ team.points }}</strong>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Current Week Results -->
        <div class="card">
          <div class="card-header">
            {{ currentWeek > 0 ? `Week ${currentWeek} Results` : 'Match Results' }}
          </div>
          <div v-if="currentWeek === 0" class="text-center" style="padding: 20px; color: #666">
            Click "Play Next Week" to start the simulation
          </div>
          <div v-else>
            <div v-for="(match, index) in currentWeekMatches" :key="index" class="match-item">
              <span class="team-name">{{ match.homeTeamName }}</span>
              <span class="score">{{ match.homeScore }} - {{ match.awayScore }}</span>
              <span class="team-name">{{ match.awayTeamName }}</span>
            </div>
          </div>

          <!-- All Played Weeks Summary -->
          <div
            v-if="currentWeek > 1"
            style="margin-top: 20px; border-top: 1px solid #dee2e6; padding-top: 15px"
          >
            <details v-for="week in previousWeeks" :key="week">
              <summary style="cursor: pointer; padding: 5px; font-weight: 500">
                Week {{ week }} Results
              </summary>
              <div
                v-for="(match, index) in getWeekMatches(week)"
                :key="index"
                class="match-item"
                style="padding-left: 20px"
              >
                <span class="team-name">{{ match.homeTeamName }}</span>
                <span class="score">{{ match.homeScore }} - {{ match.awayScore }}</span>
                <span class="team-name">{{ match.awayTeamName }}</span>
              </div>
            </details>
          </div>
        </div>

        <!-- Championship Predictions -->
        <div class="card">
          <div class="card-header">Championship Predictions</div>
          <table>
            <thead>
              <tr>
                <th>Team</th>
                <th class="text-right">%</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="prediction in state.predictions" :key="prediction.teamId">
                <td>{{ prediction.teamName }}</td>
                <td class="text-right">
                  <strong>{{ prediction.percentage }}%</strong>
                </td>
              </tr>
            </tbody>
          </table>
          <div
            v-if="currentWeek < 4 && currentWeek > 0"
            class="text-center"
            style="padding: 10px; color: #666; font-size: 0.9rem"
          >
            Predictions start from Week 4
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="btn-group">
        <button
          class="btn btn-accent"
          :disabled="actionLoading || isCompleted"
          @click="handlePlayAllWeeks"
        >
          Play All Weeks
        </button>
        <button
          class="btn btn-primary"
          :disabled="actionLoading || isCompleted"
          @click="handlePlayNextWeek"
        >
          Play Next Week
        </button>
        <button class="btn btn-danger" :disabled="actionLoading" @click="handleReset">
          Reset Data
        </button>
      </div>

      <div v-if="isCompleted" class="card text-center" style="margin-top: 20px">
        <h3>🏆 League Completed!</h3>
        <p v-if="state.standings.length">
          Champion: <strong>{{ state.standings[0].teamName }}</strong> with
          {{ state.standings[0].points }} points!
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getSimulationState, playNextWeek, playAllWeeks, resetSimulation } from '../api'

const router = useRouter()
const loading = ref(true)
const actionLoading = ref(false)
const state = ref({
  leagueState: { fixturesCreated: false, currentWeek: 0, completed: false, totalWeeks: 6 },
  standings: [],
  currentWeekResults: [],
  allMatches: {},
  predictions: [],
})

onMounted(async () => {
  await loadState()
})

const loadState = async () => {
  loading.value = true
  try {
    const response = await getSimulationState()
    state.value = response.data.data
  } catch (error) {
    console.error('Failed to load state:', error)
  } finally {
    loading.value = false
  }
}

const currentWeek = computed(() => state.value.leagueState.currentWeek)
const isCompleted = computed(() => state.value.leagueState.completed)

const currentWeekMatches = computed(() => {
  if (!state.value.allMatches || !currentWeek.value) return []
  return state.value.allMatches[currentWeek.value] || []
})

const previousWeeks = computed(() => {
  const weeks = []
  for (let i = 1; i < currentWeek.value; i++) {
    weeks.push(i)
  }
  return weeks
})

const getWeekMatches = week => {
  if (!state.value.allMatches) return []
  return state.value.allMatches[week] || []
}

const handlePlayNextWeek = async () => {
  actionLoading.value = true
  try {
    const response = await playNextWeek()
    state.value = response.data.data
  } catch (error) {
    console.error('Failed to play next week:', error)
    alert(error.response?.data?.message || 'Failed to play next week')
  } finally {
    actionLoading.value = false
  }
}

const handlePlayAllWeeks = async () => {
  actionLoading.value = true
  try {
    const response = await playAllWeeks()
    state.value = response.data.data
  } catch (error) {
    console.error('Failed to play all weeks:', error)
    alert(error.response?.data?.message || 'Failed to play all weeks')
  } finally {
    actionLoading.value = false
  }
}

const handleReset = async () => {
  if (
    !confirm(
      'Are you sure you want to reset the simulation? This will clear all match results and fixtures.'
    )
  ) {
    return
  }

  actionLoading.value = true
  try {
    await resetSimulation()
    router.push('/')
  } catch (error) {
    console.error('Failed to reset:', error)
    alert(error.response?.data?.message || 'Failed to reset simulation')
  } finally {
    actionLoading.value = false
  }
}
</script>

<style scoped>
.simulation-view {
  max-width: 1400px;
  margin: 0 auto;
}

.match-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  border-bottom: 1px solid var(--border-color);
}

.match-item:last-child {
  border-bottom: none;
}

details summary {
  color: var(--primary-color);
}

details summary:hover {
  color: var(--secondary-color);
}
</style>
