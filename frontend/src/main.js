import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import TeamsView from './views/TeamsView.vue'
import FixturesView from './views/FixturesView.vue'
import SimulationView from './views/SimulationView.vue'
import './style.css'

const routes = [
  { path: '/', name: 'Teams', component: TeamsView },
  { path: '/fixtures', name: 'Fixtures', component: FixturesView },
  { path: '/simulation', name: 'Simulation', component: SimulationView },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const app = createApp(App)
app.use(router)
app.mount('#app')
