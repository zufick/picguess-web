<script setup lang="ts">

const model = defineModel({default: "#ffffff"})

</script>

<template>
    <input v-model="model" value="#ffffff" type="color">
</template>

<style scoped>
input[type="color"] {
  -webkit-appearance: none;
  width: 100%;
  height: 44px;
  background: #fff;
  border-width: 2px;
  border-color: rgb(68 64 60);
  border-radius: .375rem;
  overflow: hidden;
  outline: none;
  cursor: pointer;
}

/* -webkit */
input[type="color"]::-webkit-color-swatch-wrapper {
  padding: 6px;
  background-color: #1c1917;
}

input[type="color"]::-webkit-color-swatch {
  border: 2px solid white;  
  border-radius: .2rem;
}


/* firefox */
input[type=color]::-moz-focus-inner {
  border: none;
  padding: 4px;
  background-color: #1c1917;
  border-radius: .1rem;
}

input[type=color]::-moz-color-swatch {
  border: none;
  border-radius: .1rem;
  height: 140px;
}
</style>