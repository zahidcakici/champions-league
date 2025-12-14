<template>
  <div class="app">
    <nav class="nav">
      <div class="nav-container">
        <router-link to="/">Teams</router-link>
        <router-link to="/fixtures">Fixtures</router-link>
        <router-link to="/simulation">Simulation</router-link>
      </div>
    </nav>
    <main class="container">
      <router-view />
    </main>
  </div>
</template>

<script setup></script>

<style>
.app {
  min-height: 100vh;
}
</style>
