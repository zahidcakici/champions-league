import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import FixturesView from '../src/views/FixturesView.vue'

// Mock the API
vi.mock('../src/api', () => ({
  getFixtures: vi.fn(() =>
    Promise.resolve({
      data: {
        data: [
          { id: 1, week: 1, homeTeam: { name: 'Chelsea' }, awayTeam: { name: 'Arsenal' } },
          { id: 2, week: 1, homeTeam: { name: 'Man City' }, awayTeam: { name: 'Liverpool' } },
          { id: 3, week: 2, homeTeam: { name: 'Arsenal' }, awayTeam: { name: 'Man City' } },
          { id: 4, week: 2, homeTeam: { name: 'Liverpool' }, awayTeam: { name: 'Chelsea' } },
        ],
      },
    })
  ),
}))

describe('FixturesView', () => {
  let router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: { template: '<div>Teams</div>' } },
        { path: '/fixtures', component: FixturesView },
        { path: '/simulation', component: { template: '<div>Simulation</div>' } },
      ],
    })
  })

  it('renders the fixtures page', async () => {
    const wrapper = mount(FixturesView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('h1').text()).toBe('Generated Fixtures')
  })

  it('shows start simulation button when fixtures exist', async () => {
    const wrapper = mount(FixturesView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    // Check for start simulation button
    const buttons = wrapper.findAll('button')
    const startButton = buttons.find(b => b.text().includes('Start Simulation'))
    expect(startButton).toBeDefined()
  })

  it('groups matches by week', async () => {
    const wrapper = mount(FixturesView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    // Should have week cards
    const weekCards = wrapper.findAll('.week-card')
    expect(weekCards.length).toBeGreaterThan(0)
  })
})
