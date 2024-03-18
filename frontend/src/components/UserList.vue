<script setup lang="ts">
import { useGameStore } from '@/stores/GameStore';
import { useWebsocketStore } from '../stores/WebsocketStore';
const wsStore = useWebsocketStore();
const gameStore = useGameStore();

function isLocalUser(id: string) {
  return wsStore.localUserId == id
}

function getPlayerScore(id: string) {

  if (!gameStore.gameState.players) {
    return "-1"
  }

  for(let i = 0; i < gameStore.gameState.players.length; i++) {
    if (gameStore.gameState.players[i].id + "" == id) {
      return gameStore.gameState.players[i].score;
    }
  }

  return "-1";
}

</script>

<template>
        <transition-group name="list" tag="div" class="flex gap-2">
            <div v-if="wsStore.roomState.clients" v-for="user in wsStore.roomState.clients.sort((a,b) => parseInt(a.id) - parseInt(b.id))" :key="user.id" class="flex p-2 gap-1 border-2 rounded-lg  bg-stone-800" 
              :class="{'border-blue-400 text-blue-400 bg-blue-900/25': isLocalUser(user.id)}">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
                </svg>
                <div>
                  <span>{{ user.username }}</span>
                  <span class="ml-4 score-text">Score: {{getPlayerScore(user.id)}}</span>
                </div>
            </div>
        </transition-group>
</template>

<style scoped>

.list-move, /* apply transition to moving elements */
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(30px);
}
.list-leave-active {
  position: absolute;
}

.score-text {
  border: 2px;
    padding: 3px 15px;
    border-radius: 0.5rem;
    background: #ffffff14;
    color: #ffffffa1;
}

</style>