<script setup lang="ts">
import JoinRoomWindow from './components/JoinRoomWindow.vue'
import DrawingCanvas from './components/DrawingCanvas.vue'
import { useWebsocketStore } from './stores/WebsocketStore';

const wsStore = useWebsocketStore();

</script>

<template>
    <main>
        <h1 v-if="wsStore.connectionError">Error: connection failed</h1>
        <div v-else>
            <JoinRoomWindow v-if="!wsStore.isJoinedRoom" />
            <div v-else>
                <p>Joined room <br> <b>{{ wsStore.joinedRoomId  }}</b></p>
                <DrawingCanvas />
            </div>
        </div>
    </main>
</template>

<style scoped>
</style>
