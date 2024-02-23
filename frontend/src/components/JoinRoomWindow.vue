<script setup lang="ts">
import { ref } from 'vue'
import { useWebsocketStore } from '@/stores/WebsocketStore';


const wsStore = useWebsocketStore();
const enteredRoomId = ref("")
</script>

<template>
    <div class="mb-6">
        <p>Join room</p>
        <div class="flex items-center">
            <input v-model="enteredRoomId" type="text" placeholder="Room id">
            <button @click="wsStore.joinRoom(enteredRoomId)">Join</button>
        </div>
    </div>
    <button @click="wsStore.createRoom()">Create room</button>
</template>
