<script setup lang="ts">
import { computed } from 'vue';
import DrawingCanvas from './DrawingCanvas.vue'
import { useGameStore } from '../stores/GameStore';
import { useWebsocketStore } from '@/stores/WebsocketStore';

const gameStore = useGameStore();
const wsStore = useWebsocketStore();

function isWordCorrectEquals(word: string, checkForValue: boolean) {
    if (gameStore.gameState.player.answerResults == undefined) {
        return false;
    }

    let isCorrect = false
    gameStore.gameState.player.answerResults.forEach(i => {
        if (i.answer == word) {
            isCorrect = i.isCorrect == checkForValue
            return;
        }
    })
    return isCorrect
}
</script>

<template>
    <div class="window">
        <div class="word-container">
            <button @click="wsStore.pickAnswer(opponentWord)" 
                class="p-2 px-4 bg-stone-900 rounded-md border-2 border-stone-700 w-40 h-20 text-center flex items-center justify-center" 
                v-for="(opponentWord) in gameStore.gameState.player.wordPool"
                :class="{ 'correct': isWordCorrectEquals(opponentWord, true), 'wrong': isWordCorrectEquals(opponentWord, false) }">
                <span>{{ opponentWord }}</span>
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
        border: 2px solid rgb(68 64 60);
        border-radius: 10px;
        color: rgb(68 64 60);
        background: rgb(28 25 23);
        box-shadow: 0px 20px 20px #0000003d;
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

    .wrong {
        background: rgb(76 5 25);
        color: rgb(225 29 72);
        border: 1px solid rgb(225 29 72);
    }

    .correct {
        background: rgb(2 44 34);
        color: rgb(52 211 153);
        border: 1px solid rgb(52 211 153);
    }

</style>