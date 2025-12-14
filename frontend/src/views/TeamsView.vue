<template>
  <div class="teams-view">
    <h1 class="page-title">Tournament Teams</h1>

    <!-- Add Team Form -->
    <div v-if="!fixturesGenerated" class="add-team-form">
      <h3>Add New Team</h3>
      <form @submit.prevent="handleAddTeam">
        <div class="form-row">
          <input
            v-model="newTeam.name"
            type="text"
            placeholder="Team Name"
            class="form-input"
            required
          />
          <input
            v-model.number="newTeam.power"
            type="number"
            min="1"
            max="100"
            placeholder="Power (1-100)"
            class="form-input power-input"
            required
          />
          <button type="submit" class="btn btn-success" :disabled="loading">Add Team</button>
        </div>
        <p v-if="formError" class="error-text">{{ formError }}</p>
      </form>
    </div>

    <div class="table-container">
      <table>
        <thead>
          <tr>
            <th>Team Name</th>
            <th class="text-center">Power Rating</th>
            <th v-if="!fixturesGenerated" class="text-center">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="team in teams" :key="team.id">
            <td>{{ team.name }}</td>
            <td class="text-center">{{ team.power }}</td>
            <td v-if="!fixturesGenerated" class="text-center">
              <button
                class="btn btn-danger btn-sm"
                :disabled="loading"
                @click="handleDeleteTeam(team.id)"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="btn-group">
      <button
        class="btn btn-primary"
        :disabled="loading || fixturesGenerated || teams.length < 2"
        @click="handleGenerateFixtures"
      >
        {{ fixturesGenerated ? 'Fixtures Generated' : 'Generate Fixtures' }}
      </button>
    </div>
    <p v-if="teams.length < 2 && !fixturesGenerated" class="hint-text">
      Add at least 2 teams to generate fixtures
    </p>

    <div v-if="loading" class="loading">Loading...</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getTeams, createTeam, deleteTeam, generateFixtures, getSimulationState } from '../api'

const router = useRouter()
const teams = ref([])
const loading = ref(false)
const fixturesGenerated = ref(false)
const formError = ref('')
const newTeam = ref({
  name: '',
  power: 50,
})

onMounted(async () => {
  await loadTeams()
  await checkFixturesStatus()
})

const loadTeams = async () => {
  try {
    const response = await getTeams()
    teams.value = response.data.data
  } catch (error) {
    console.error('Failed to load teams:', error)
  }
}

const checkFixturesStatus = async () => {
  try {
    const response = await getSimulationState()
    fixturesGenerated.value = response.data.data.leagueState.fixturesCreated
  } catch (_error) {
    // State might not exist yet
    fixturesGenerated.value = false
  }
}

const handleAddTeam = async () => {
  formError.value = ''

  if (!newTeam.value.name.trim()) {
    formError.value = 'Team name is required'
    return
  }
  if (newTeam.value.power < 1 || newTeam.value.power > 100) {
    formError.value = 'Power must be between 1 and 100'
    return
  }

  loading.value = true
  try {
    await createTeam(newTeam.value.name.trim(), newTeam.value.power)
    newTeam.value = { name: '', power: 50 }
    await loadTeams()
  } catch (error) {
    formError.value = error.response?.data?.error || 'Failed to create team'
    console.error('Failed to create team:', error)
  } finally {
    loading.value = false
  }
}

const handleDeleteTeam = async id => {
  loading.value = true
  try {
    await deleteTeam(id)
    teams.value = teams.value.filter(t => t.id !== id)
  } catch (error) {
    console.error('Failed to delete team:', error)
  } finally {
    loading.value = false
  }
}

const handleGenerateFixtures = async () => {
  loading.value = true
  try {
    await generateFixtures()
    fixturesGenerated.value = true
    router.push('/fixtures')
  } catch (error) {
    console.error('Failed to generate fixtures:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.teams-view {
  max-width: 600px;
  margin: 0 auto;
}

.add-team-form {
  background: var(--card-bg);
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  border: 1px solid var(--border-color);
}

.add-team-form h3 {
  margin: 0 0 1rem 0;
  font-size: 1rem;
  color: var(--text-secondary);
}

.form-row {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.form-input {
  padding: 0.5rem 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-color);
  color: var(--text-color);
  font-size: 0.9rem;
}

.form-input:focus {
  outline: none;
  border-color: var(--primary-color);
}

.form-input:first-child {
  flex: 1;
}

.power-input {
  width: 120px;
}

.btn-success {
  background-color: #10b981;
}

.btn-success:hover:not(:disabled) {
  background-color: #059669;
}

.error-text {
  color: #ef4444;
  font-size: 0.85rem;
  margin: 0.5rem 0 0 0;
}

.hint-text {
  text-align: center;
  color: var(--text-secondary);
  font-size: 0.85rem;
  margin-top: 0.5rem;
}

.btn-danger {
  background-color: #ef4444;
}

.btn-danger:hover:not(:disabled) {
  background-color: #dc2626;
}

.btn-sm {
  padding: 0.25rem 0.5rem;
  font-size: 0.8rem;
}
</style>
