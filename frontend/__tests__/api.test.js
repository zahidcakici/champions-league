import { describe, it, expect, vi } from 'vitest'
import axios from 'axios'

// Mock axios
vi.mock('axios', () => ({
  default: {
    create: vi.fn(() => ({
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
    })),
  },
}))

describe('API Client', () => {
  it('creates axios instance with correct base URL', async () => {
    // Import after mocking to trigger axios.create call
    await import('../src/api')

    expect(axios.create).toHaveBeenCalledWith({
      baseURL: '/api',
      headers: {
        'Content-Type': 'application/json',
      },
    })
  })
})

describe('API Endpoints', () => {
  it('should have teams endpoint', async () => {
    const { getTeams } = await import('../src/api')
    expect(typeof getTeams).toBe('function')
  })

  it('should have fixtures endpoints', async () => {
    const { getFixtures, getFixturesByWeek, generateFixtures } = await import('../src/api')
    expect(typeof getFixtures).toBe('function')
    expect(typeof getFixturesByWeek).toBe('function')
    expect(typeof generateFixtures).toBe('function')
  })

  it('should have simulation endpoints', async () => {
    const { getSimulationState, playNextWeek, playAllWeeks, updateMatchResult, resetSimulation } =
      await import('../src/api')

    expect(typeof getSimulationState).toBe('function')
    expect(typeof playNextWeek).toBe('function')
    expect(typeof playAllWeeks).toBe('function')
    expect(typeof updateMatchResult).toBe('function')
    expect(typeof resetSimulation).toBe('function')
  })

  it('should have standings endpoints', async () => {
    const { getStandings, getPredictions } = await import('../src/api')
    expect(typeof getStandings).toBe('function')
    expect(typeof getPredictions).toBe('function')
  })
})
