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
header {
    line-height: 1.5;
}

.logo {
    display: block;
    margin: 0 auto 2rem;
}

@media (min-width: 1024px) {
    header {
        display: flex;
        place-items: center;
        padding-right: calc(var(--section-gap) / 2);
    }

    .logo {
        margin: 0 2rem 0 0;
    }

    header .wrapper {
        display: flex;
        place-items: flex-start;
        flex-wrap: wrap;
    }
}
</style>
