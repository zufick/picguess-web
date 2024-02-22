<script setup lang="ts">
import { ref } from 'vue'

const enteredRoomId = ref("");
    
function joinGame() {
    alert("join")
}
</script>

<template>
    <div class="mb-6">
        <p>Join room</p>
        <div class="flex">
            <input
            type="text" placeholder="Room id">
            <button @click="joinGame">Join</button>
        </div>
    </div>
    <button>Create room</button>
</template>
