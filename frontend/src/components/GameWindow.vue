<script setup lang="ts">
import { ref } from 'vue'
import DrawingCanvas from './DrawingCanvas.vue'
import GameWordPoolWindow from './GameWordPoolWindow.vue'
import GameAnswers from './GameAnswers.vue'
import GameWinner from './GameWinner.vue'
import { useGameStore } from '../stores/GameStore';

const gameStore = useGameStore();
const isWordPoolOpened = ref(false)

function wordsClicked() {
    isWordPoolOpened.value = !isWordPoolOpened.value
}

function isWordCorrectEquals(word: string, checkForValue: boolean) {

    if (gameStore.gameState.player.opponentResults == undefined) {
        return false;
    }

    let isCorrect = false
    gameStore.gameState.player.opponentResults.forEach(i => {
        if (i.answer == word) {
            isCorrect = i.isCorrect == checkForValue
            return;
        }
    })
    return isCorrect
}

function getWrongOpponentWords() {
    if (gameStore.gameState.player.opponentResults == undefined) {
        return [];
    }

    return gameStore.gameState.player.opponentResults.filter(r => !r.isCorrect);
}

</script>

<template>
    <div class="relative">
        <div class="flex gap-2 word-container my-2">
            <div class="p-2 px-4 bg-stone-900 rounded-md border-2 border-stone-700" 
            v-for="opponentWord in gameStore.gameState.player.opponentWinWords"
            :class="{ 
                'correct': isWordCorrectEquals(opponentWord, true), 
                'wrong': isWordCorrectEquals(opponentWord, false) 
            }">{{ opponentWord }}</div>
        </div>
        <div class="flex gap-2 word-container-wrong">
            <div class="p-2 px-4 bg-stone-900 rounded-md border-2 border-stone-700 wrong" 
            v-for="opponentWord in getWrongOpponentWords()">{{ opponentWord.answer }}</div>
        </div>
        <DrawingCanvas />
        <button class="wordpool-btn" @click="wordsClicked">
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" d="m10.5 21 5.25-11.25L21 21m-9-3h7.5M3 5.621a48.474 48.474 0 0 1 6-.371m0 0c1.12 0 2.233.038 3.334.114M9 5.25V3m3.334 2.364C11.176 10.658 7.69 15.08 3 17.502m9.334-12.138c.896.061 1.785.147 2.666.257m-4.589 8.495a18.023 18.023 0 0 1-3.827-5.802" />
            </svg>
        </button>
        <GameAnswers/>
        <GameWordPoolWindow v-if="isWordPoolOpened"/>
        <GameWinner v-if="gameStore.gameState.winner" :winner="gameStore.gameState.winner"/>
    </div>
</template>

<style scoped>
    .word-container {
        pointer-events: none;
    }

    .word-container-wrong {
        position: absolute;
        padding: 20px;
        pointer-events: none;
        z-index: 10;
        bottom: 0;
        opacity: .5;
    }

    .container > * {
        pointer-events: none;
    }

    
    .wordpool-btn {
        position: absolute;
        bottom: 0px;
        right: 0px;
        z-index: 30;
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