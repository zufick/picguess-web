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
        <div class="flex gap-2 word-container">
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
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 20.25c4.97 0 9-3.694 9-8.25s-4.03-8.25-9-8.25S3 7.444 3 12c0 2.104.859 4.023 2.273 5.48.432.447.74 1.04.586 1.641a4.483 4.483 0 0 1-.923 1.785A5.969 5.969 0 0 0 6 21c1.282 0 2.47-.402 3.445-1.087.81.22 1.668.337 2.555.337Z" />
            </svg>
        </button>
        <GameAnswers/>
        <GameWordPoolWindow v-if="isWordPoolOpened"/>
        <GameWinner v-if="gameStore.gameState.winner" :winner="gameStore.gameState.winner"/>
    </div>
</template>

<style scoped>
    .word-container {
        position: absolute;
        padding: 20px;
        pointer-events: none;
        z-index: 10;
    }

    .word-container-wrong {
        position: absolute;
        padding: 20px;
        pointer-events: none;
        z-index: 10;
        bottom: 0;
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