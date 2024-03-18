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

            <svg data-name="Слой 2" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 15.76 21" width="24" height="24" fill="currentColor">
                <g>
                    <path class="cls-1" d="m12.76,0H3C1.35,0,0,1.35,0,3v13.25c0,1.59,1.29,2.89,2.89,2.89h.97l.05.89c.02.42.29.78.69.92.36.13.76.04,1.04-.22l1.04-.78,1.04.78c.19.18.44.28.7.28.11,0,.23-.02.34-.06.4-.14.67-.5.69-.91l.05-.89h3.38c1.59,0,2.89-1.29,2.89-2.89V3c0-1.65-1.35-3-3-3ZM2.89,17.64c-.76,0-1.39-.62-1.39-1.39v-.73c.44.26.95.42,1.5.42h.68l.09,1.7h-.89Zm4.41.9c-.37-.27-.87-.27-1.23,0l-.7.53-.17-3.12h2.98l-.18,3.12-.7-.52Zm5.58-.9h-3.3l.1-1.7h3.09c.55,0,1.06-.16,1.5-.42v.73c0,.76-.62,1.39-1.39,1.39Zm-.12-3.2H3c-.83,0-1.5-.67-1.5-1.5V3c0-.83.67-1.5,1.5-1.5h9.75c.83,0,1.5.67,1.5,1.5v9.94c0,.83-.67,1.5-1.5,1.5Z"/>
                    <path class="cls-1" d="m11.73,8.36h-1.25c-.6-1.4-1.14-2.82-1.59-4.26-.14-.44-.55-.74-1.01-.74h0c-.46,0-.87.3-1.01.74-.45,1.44-.99,2.86-1.59,4.26h-1.25c-.41,0-.75.34-.75.75s.34.75.75.75h.58c-.19.4-.38.8-.58,1.2l-.22.44c-.19.37-.04.82.33,1.01.11.05.22.08.34.08.27,0,.54-.15.67-.41l.22-.43c.31-.61.59-1.25.87-1.88h3.28c.29.63.56,1.26.87,1.88l.22.44c.19.37.64.52,1.01.33.37-.19.52-.64.33-1.01l-.22-.44c-.2-.39-.39-.8-.58-1.2h.58c.41,0,.75-.34.75-.75s-.34-.75-.75-.75Zm-4.83,0c.35-.84.68-1.68.98-2.54.3.86.63,1.7.98,2.54h-1.96Z"/>
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