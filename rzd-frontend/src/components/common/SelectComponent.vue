<script setup lang="ts">
  import { ref } from 'vue';

  const selected = ref('A')
</script>

<template>
  <div class="selector">
    <div>Selected: {{ selected }}</div>

    <select v-model="selected">
      <option disabled value="">Please select one</option>
      <option>A</option>
      <option>B</option>
      <option>C</option>
    </select>
  </div>
</template>

<style scoped>
  .rzd-wrapper {
    width: 240;
    height: 120px;
    background-color: gray;
    color: black;
  }
</style>