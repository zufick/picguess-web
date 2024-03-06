<script setup lang="ts">
import DrawingCanvas from './DrawingCanvas.vue'
import { useGameStore } from '../stores/GameStore';

const gameStore = useGameStore();


</script>

<template>
    <div class="window">
        <div class="word-container">
            <button class="p-2 px-4 bg-stone-900 rounded-md border-2 border-stone-700 w-40 h-20 text-center flex items-center justify-center" v-for="opoonentWord in gameStore.gameState.player.wordPool">
                <span>{{ opoonentWord }}</span>
            </button>
        </div>

    </div>
</template>

<style scoped>
    .window {
        position: absolute;
        top: 0;
        left: 0;
        z-index: 20;
        width: 100%;
        height: 100%;

        backdrop-filter: blur(10px);
        border: 1px solid white;
        border-radius: 10px;
        background: #1c191773;
    }

    .word-container {
        width: 100%;
        height: 100%;
        display: grid;
        grid-template-columns: 20% 20% 20% 20%;
        grid-gap: 10px;
        align-items: center;
        justify-content: center;
        justify-items: center;
    }

</style>