<script setup lang="ts">
import { ref } from 'vue';
import ConnectionWindow from './components/ConnectionWindow.vue'

const socket = new WebSocket("ws://localhost:8080");

socket.addEventListener("error", (event) => {
    connectionError.value = true
});

socket.addEventListener("close", (event) => {
    connectionError.value = true
});

// Listen for messages
socket.addEventListener("message", (event) => {
    console.log("Message from server ", event.data);
});

const connectionError = ref(false)

</script>

<template>
    <main>
        <h1 v-if="connectionError">ERROR</h1>
        <ConnectionWindow />
    </main>
</template>

<style scoped>
</style>
