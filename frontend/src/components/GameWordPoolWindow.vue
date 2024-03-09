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
            <p class="text-white/50">Opponent's possible words</p>
            <button @click="wsStore.pickAnswer(opponentWord)" 
                class="p-2 px-4 bg-stone-900 rounded-md border-2 border-stone-700 w-full h-20 text-center flex items-center justify-center" 
                v-for="(opponentWord) in gameStore.gameState.player.wordPool"
                :class="{ 'correct': isWordCorrectEquals(opponentWord, true), 'wrong': isWordCorrectEquals(opponentWord, false) }">
                <span>{{ opponentWord }}</span>
            </button>
        </div>

    </div>
</template>

<style scoped>
    .window {
        animation: .1s ease appear;
        position: absolute;
        top: 0;
        right: 0;
        z-index: 20;
        width: 300px;
        height: calc(100% - 60px);
        border: 2px solid rgb(68 64 60);
        border-radius: 10px;
        color: rgb(68 64 60);
        background: rgba(28, 25, 23, .95);
        box-shadow: 0px 20px 20px #0000003d;
    }

    .word-container {
        width: 100%;
        height: 100%;
        padding: 30px;
        display: flex;
        flex-direction: column;
        gap: 15px;
        align-items: center;
        justify-items: center;
        overflow-y: scroll;
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

    ::-webkit-scrollbar {
    width: 20px;
    }

    /* Track */
    ::-webkit-scrollbar-track {
        background: #282523;
        border-radius: .5rem;
    }
    
    /* Handle */
    ::-webkit-scrollbar-thumb {
        background: #33302E;
    border-radius: .5rem;
    }

    /* Handle on hover */
    ::-webkit-scrollbar-thumb:hover {
        background: #413F3D;
    }

</style>