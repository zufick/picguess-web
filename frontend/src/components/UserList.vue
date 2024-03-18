<script setup lang="ts">
import { useGameStore } from '@/stores/GameStore';
import { useWebsocketStore } from '../stores/WebsocketStore';
import UserListCard from './UserListCard.vue';
const wsStore = useWebsocketStore();
const gameStore = useGameStore();

function getPlayerScore(id: string) {

  if (!gameStore.gameState.players) {
    return "0"
  }

  for(let i = 0; i < gameStore.gameState.players.length; i++) {
    if (gameStore.gameState.players[i].id + "" == id) {
      return gameStore.gameState.players[i].score;
    }
  }

  return "0";
}

</script>

<template>
        <transition-group name="list" tag="div" class="flex gap-2">
          <UserListCard v-if="wsStore.roomState.clients" v-for="user in wsStore.roomState.clients.sort((a,b) => parseInt(a.id) - parseInt(b.id))" :key="user.id" :user="user" :score="getPlayerScore(user.id)" />
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
</style>