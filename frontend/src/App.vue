<script setup lang="ts">
import { ref, provide } from 'vue';
import JoinRoomWindow from './components/JoinRoomWindow.vue'
import { useWebsocketStore } from './stores/WebsocketStore';

const wsStore = useWebsocketStore();

</script>

<template>
    <main>
        <h1 v-if="wsStore.connectionError">Error: connection failed</h1>
        <div v-else>
            <JoinRoomWindow v-if="!wsStore.isJoinedRoom" />
            <div v-else>Joined room <br> <b>{{ wsStore.joinedRoomId  }}</b> </div>
        </div>
    </main>
</template>

<style scoped>
</style>
