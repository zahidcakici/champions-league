import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createRouter, createWebHistory } from 'vue-router'
import TeamsView from '../src/views/TeamsView.vue'

// Mock the API
vi.mock('../src/api', () => ({
  getTeams: vi.fn(() =>
    Promise.resolve({
      data: {
        data: [
          { id: 1, name: 'Chelsea', power: 85 },
          { id: 2, name: 'Arsenal', power: 80 },
          { id: 3, name: 'Manchester City', power: 90 },
          { id: 4, name: 'Liverpool', power: 82 },
        ],
      },
    })
  ),
  getSimulationState: vi.fn(() =>
    Promise.resolve({
      data: {
        data: {
          leagueState: { fixturesCreated: false },
        },
      },
    })
  ),
  generateFixtures: vi.fn(() => Promise.resolve({ data: { success: true } })),
  createTeam: vi.fn(() =>
    Promise.resolve({ data: { data: { id: 5, name: 'New Team', power: 50 } } })
  ),
  deleteTeam: vi.fn(() => Promise.resolve({ data: { data: { deleted: true } } })),
}))

describe('TeamsView', () => {
  let router

  beforeEach(() => {
    router = createRouter({
      history: createWebHistory(),
      routes: [
        { path: '/', component: TeamsView },
        { path: '/fixtures', component: { template: '<div>Fixtures</div>' } },
      ],
    })
  })

  it('renders the teams table', async () => {
    const wrapper = mount(TeamsView, {
      global: {
        plugins: [router],
      },
    })

    // Wait for async operations
    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('h1').text()).toBe('Tournament Teams')
    expect(wrapper.find('table').exists()).toBe(true)
  })

  it('displays generate fixtures button', async () => {
    const wrapper = mount(TeamsView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    const button = wrapper.find('.btn-primary')
    expect(button.exists()).toBe(true)
    expect(button.text()).toContain('Generate Fixtures')
  })

  it('has correct page title', () => {
    const wrapper = mount(TeamsView, {
      global: {
        plugins: [router],
      },
    })

    expect(wrapper.find('.page-title').exists()).toBe(true)
  })

  it('displays add team form when fixtures not generated', async () => {
    const wrapper = mount(TeamsView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.add-team-form').exists()).toBe(true)
    expect(wrapper.find('input[placeholder="Team Name"]').exists()).toBe(true)
    expect(wrapper.find('.btn-success').text()).toContain('Add Team')
  })

  it('displays delete buttons for each team', async () => {
    const wrapper = mount(TeamsView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    const deleteButtons = wrapper.findAll('.btn-danger')
    expect(deleteButtons.length).toBe(4) // 4 teams = 4 delete buttons
  })

  it('calls createTeam API when adding a team', async () => {
    const { createTeam } = await import('../src/api')

    const wrapper = mount(TeamsView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    // Fill in the form
    const nameInput = wrapper.find('input[placeholder="Team Name"]')
    const powerInput = wrapper.find('.power-input')

    await nameInput.setValue('Test Team')
    await powerInput.setValue(75)

    // Submit the form
    const form = wrapper.find('form')
    await form.trigger('submit.prevent')

    await new Promise(resolve => setTimeout(resolve, 100))

    expect(createTeam).toHaveBeenCalledWith('Test Team', 75)
  })

  it('calls deleteTeam API when clicking delete button', async () => {
    const { deleteTeam } = await import('../src/api')

    const wrapper = mount(TeamsView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    // Click the first delete button
    const deleteButtons = wrapper.findAll('.btn-danger')
    await deleteButtons[0].trigger('click')

    await new Promise(resolve => setTimeout(resolve, 100))

    expect(deleteTeam).toHaveBeenCalledWith(1) // First team ID
  })

  it('navigates to fixtures after generating fixtures', async () => {
    const wrapper = mount(TeamsView, {
      global: {
        plugins: [router],
      },
    })

    await router.isReady()
    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    const generateButton = wrapper.find('.btn-primary')
    await generateButton.trigger('click')

    await new Promise(resolve => setTimeout(resolve, 100))

    expect(router.currentRoute.value.path).toBe('/fixtures')
  })

  it('disables generate button when less than 2 teams', async () => {
    // Override getTeams to return only 1 team
    const { getTeams } = await import('../src/api')
    getTeams.mockResolvedValueOnce({
      data: {
        data: [{ id: 1, name: 'Only Team', power: 80 }],
      },
    })

    const wrapper = mount(TeamsView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    const generateButton = wrapper.find('.btn-primary')
    expect(generateButton.attributes('disabled')).toBeDefined()
  })

  it('shows hint text when less than 2 teams', async () => {
    const { getTeams } = await import('../src/api')
    getTeams.mockResolvedValueOnce({
      data: {
        data: [{ id: 1, name: 'Only Team', power: 80 }],
      },
    })

    const wrapper = mount(TeamsView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('.hint-text').exists()).toBe(true)
    expect(wrapper.find('.hint-text').text()).toContain('at least 2 teams')
  })

  it('handles API error gracefully when creating team fails', async () => {
    const { createTeam } = await import('../src/api')
    createTeam.mockRejectedValueOnce({
      response: { data: { error: 'Team name already exists' } },
    })

    const wrapper = mount(TeamsView, {
      global: {
        plugins: [router],
      },
    })

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    // Fill and submit form
    await wrapper.find('input[placeholder="Team Name"]').setValue('Duplicate')
    await wrapper.find('.power-input').setValue(80)
    await wrapper.find('form').trigger('submit.prevent')

    await new Promise(resolve => setTimeout(resolve, 100))
    await wrapper.vm.$nextTick()

    // Should show error message
    expect(wrapper.find('.error-text').exists()).toBe(true)
  })
})
