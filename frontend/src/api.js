import axios from 'axios'

// Determine the API base URL from environment variables
// - Empty string: uses proxy (Vite dev server or nginx in Docker)
// - Full URL: direct API calls (for S3/CDN deployments)
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL
  ? `${import.meta.env.VITE_API_BASE_URL}/api`
  : '/api'

const api = axios.create({
  baseURL: apiBaseUrl,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Teams
export const getTeams = () => api.get('/teams')
export const createTeam = (name, power) => api.post('/teams', { name, power })
export const deleteTeam = id => api.delete(`/teams/${id}`)

// Fixtures
export const getFixtures = () => api.get('/fixtures')
export const getFixturesByWeek = week => api.get(`/fixtures/${week}`)
export const generateFixtures = () => api.post('/fixtures/generate')

// Simulation
export const getSimulationState = () => api.get('/simulation/state')
export const playNextWeek = () => api.post('/simulation/play-week')
export const playAllWeeks = () => api.post('/simulation/play-all')
export const updateMatchResult = (matchId, homeScore, awayScore) =>
  api.put(`/simulation/match/${matchId}`, { homeScore, awayScore })
export const resetSimulation = () => api.post('/simulation/reset')

// Standings
export const getStandings = () => api.get('/standings')
export const getPredictions = () => api.get('/predictions')

export default api
