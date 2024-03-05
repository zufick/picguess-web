<script setup lang="ts">
import { useWebsocketStore } from '../stores/WebsocketStore';
const wsStore = useWebsocketStore();
</script>

<template>
        <transition-group name="list" tag="div" class="flex gap-2">
            <div v-for="user in wsStore.gameState.clients" :key="user.id" class="flex p-2 gap-1 border-2 rounded-lg border-blue-400 bg-stone-800 text-blue-400	bg-blue-900/25">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
                </svg>
                <span>{{ user.username }}</span>
            </div>
        </transition-group>
</template>

<style>

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