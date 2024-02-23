<script setup lang="ts">
import { inject, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useRoomStore } from '@/stores/RoomStore'

const roomStore = useRoomStore();
const { joinRoomId } = storeToRefs(roomStore)
const { joinRoom } = roomStore

</script>

<template>
    <div class="mb-6">
        <p>Join room</p>
        <div class="flex">
            <input v-model="joinRoomId" type="text" placeholder="Room id">
            <button @click="joinRoom()">Join</button>
        </div>
    </div>
    <button>Create room</button>
</template>
