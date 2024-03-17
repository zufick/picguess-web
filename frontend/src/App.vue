<script setup lang="ts">
import JoinRoomWindow from './components/JoinRoomWindow.vue'
import DrawingCanvas from './components/DrawingCanvas.vue'
import UserList from './components/UserList.vue'
import GameWindow from './components/GameWindow.vue'

import { useWebsocketStore } from './stores/WebsocketStore';
import { useGameStore } from './stores/GameStore';
import { RouterView, useRoute } from 'vue-router'
import { onMounted } from 'vue';

const wsStore = useWebsocketStore();
const gameStore = useGameStore();
const route = useRoute();

function devAutoJoin() {
    // @ts-ignore
    const isDev = import.meta.env.DEV

    setTimeout(() => {
        if (isDev) {
            if (route.params.room == undefined) {
                wsStore.createRoom({username: Math.random()*1000 + ""})
                setTimeout(() => {
                    // @ts-ignore
                    window.open(window.location.href, '_blank');
                }, 200)

            } else {
                wsStore.startGame()
            }
        }
    }, 300)
}

onMounted(() => {
    //devAutoJoin();
})

</script>

<template>
    <main>
        <h1 v-if="wsStore.connectionError">Error: connection failed</h1>
        <div v-if="wsStore.isConnected">
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
