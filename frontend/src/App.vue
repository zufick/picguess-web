<script setup lang="ts">
import JoinRoomWindow from './components/JoinRoomWindow.vue'
import DrawingCanvas from './components/DrawingCanvas.vue'
import UserList from './components/UserList.vue'
import GameWindow from './components/GameWindow.vue'

import { useWebsocketStore } from './stores/WebsocketStore';
import { useGameStore } from './stores/GameStore';

const wsStore = useWebsocketStore();
const gameStore = useGameStore();

</script>

<template>
    <main>
        <h1 v-if="wsStore.connectionError">Error: connection failed</h1>
        <div v-else>
            <JoinRoomWindow v-if="!wsStore.isJoinedRoom" />
            <div v-else>
                <p>Joined room <br> <b>{{ wsStore.joinedRoomId  }}</b></p>
                <UserList/>
                
                <GameWindow v-if="gameStore.gameState.player"/>
                <button v-else @click="wsStore.startGame()">Start game</button>
            </div>
        </div>
    </main>
</template>

<style scoped>
</style>
