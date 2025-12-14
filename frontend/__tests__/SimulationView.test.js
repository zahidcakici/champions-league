import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import SimulationView from '../src/views/SimulationView.vue'

// Mock the API
vi.mock('../src/api', () => ({
  getSimulationState: vi.fn(() =>
    Promise.resolve({
      data: {
        data: {
          leagueState: {
            fixturesCreated: true,
            currentWeek: 2,
            completed: false,
            totalWeeks: 6,
          },
          standings: [
            {
              teamId: 1,
              teamName: 'Chelsea',
              played: 2,
              won: 2,
              drawn: 0,
              lost: 0,
              goalDifference: 3,
              points: 6,
            },
            {
              teamId: 2,
              teamName: 'Arsenal',
              played: 2,
              won: 1,
              drawn: 1,
              lost: 0,
              goalDifference: 1,
              points: 4,
            },
            {
              teamId: 3,
              teamName: 'Man City',
              played: 2,
              won: 0,
              drawn: 1,
              lost: 1,
              goalDifference: -1,
              points: 1,
            },
            {
              teamId: 4,
              teamName: 'Liverpool',
              played: 2,
              won: 0,
              drawn: 0,
              lost: 2,
              goalDifference: -3,
              points: 0,
            },
          ],
          predictions: [
            { teamId: 1, teamName: 'Chelsea', percentage: 0 },
            { teamId: 2, teamName: 'Arsenal', percentage: 0 },
            { teamId: 3, teamName: 'Man City', percentage: 0 },
            { teamId: 4, teamName: 'Liverpool', percentage: 0 },
          ],
          allMatches: {
            1: [
              {
                homeTeamName: 'Chelsea',
                awayTeamName: 'Arsenal',
                homeScore: 2,
                awayScore: 1,
              },
            ],
            2: [
              {
                homeTeamName: 'Man City',
                awayTeamName: 'Liverpool',
                homeScore: 1,
                awayScore: 0,
              },
            ],
          },
          currentWeekResults: [],
        },
      },
    })
  ),
  playNextWeek: vi.fn(() => Promise.resolve({ data: { data: {} } })),
  playAllWeeks: vi.fn(() => Promise.resolve({ data: { data: {} } })),
  resetSimulation: vi.fn(() => Promise.resolve({ data: { success: true } })),
}))

describe('SimulationView', () => {
  let router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div>Teams</div>' } },
        { path: '/simulation', component: SimulationView },
      ],
    })
  })

  it('renders the simulation page', async () => {
    const wrapper = mount(SimulationView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('h1').text()).toBe('Simulation')
  })

  it('displays league table', async () => {
    const wrapper = mount(SimulationView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    // Check for standings table
    expect(wrapper.text()).toContain('League Table')
  })

  it('shows action buttons', async () => {
    const wrapper = mount(SimulationView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    const buttons = wrapper.findAll('button')
    expect(buttons.length).toBeGreaterThan(0)

    // Check for specific buttons
    const buttonTexts = buttons.map(b => b.text())
    expect(buttonTexts.some(t => t.includes('Play'))).toBe(true)
    expect(buttonTexts.some(t => t.includes('Reset'))).toBe(true)
  })

  it('displays predictions section', async () => {
    const wrapper = mount(SimulationView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('Championship Predictions')
  })
})
