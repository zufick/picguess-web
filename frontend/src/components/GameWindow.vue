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

            <svg data-name="Слой 2" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 359.97 512.38" width="24" height="24" fill="currentColor">
                <g>
                    <g>
                    <path class="cls-1" d="m294.19,0h-54.86s-.04,0-.05,0H43.62C19.57,0,0,19.57,0,43.62v349.67c0,24.05,19.57,43.62,43.62,43.62h11.37l3.16,60.37c.34,6.5,4.53,12.08,10.68,14.21,5.72,1.98,11.96.58,16.25-3.56l23.65-17.69,23.65,17.69c3,2.89,6.95,4.45,11,4.45,1.75,0,3.53-.29,5.25-.89,6.15-2.14,10.34-7.71,10.68-14.21l3.16-60.37h131.71c36.27,0,65.77-29.51,65.77-65.77V65.78C359.97,29.51,330.46,0,294.19,0Zm-155.13,486.99l-20.79-15.55c-5.61-4.2-13.46-4.2-19.08,0l-20.78,15.55-2.62-50.08h65.89l-2.62,50.08Zm-95.44-70.85c-12.6,0-22.85-10.25-22.85-22.85V43.62c0-12.6,10.25-22.85,22.85-22.85h195.67s.05,0,.05,0c24.79.03,44.95,20.21,44.95,45v305.36c0,24.82-20.19,45-45,45H43.62Zm295.58-45c0,24.82-20.19,45-45,45h-7.01c11.07-11.78,17.87-27.61,17.87-45V65.77c0-17.4-6.8-33.23-17.87-45h7c24.82,0,45,20.19,45,45v305.36Z"/>
                    <path class="cls-1" d="m89.67,240.27c5.3,2.19,11.37-.32,13.57-5.62l16.4-39.6h65.76l16.4,39.6c1.66,4,5.52,6.41,9.6,6.41,1.32,0,2.67-.25,3.97-.79,5.3-2.2,7.81-8.27,5.62-13.57l-19.07-46.03s-.02-.06-.03-.08l-34.66-83.68c-2.47-5.97-8.25-9.83-14.71-9.83s-12.24,3.86-14.71,9.83l-34.66,83.68s-.02.06-.03.08l-19.07,46.03c-2.19,5.3.32,11.37,5.62,13.57Zm62.85-124.61l24.28,58.62h-48.56l24.28-58.62Z"/>
                    <path class="cls-1" d="m229.65,279.39H75.41c-5.74,0-10.39,4.65-10.39,10.39s4.65,10.39,10.39,10.39h154.25c5.74,0,10.39-4.65,10.39-10.39s-4.65-10.39-10.39-10.39Z"/>
                    <path class="cls-1" d="m229.65,329.06H75.41c-5.74,0-10.39,4.65-10.39,10.39s4.65,10.39,10.39,10.39h154.25c5.74,0,10.39-4.65,10.39-10.39s-4.65-10.39-10.39-10.39Z"/>
                    </g>
                </g>
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