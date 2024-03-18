<script setup lang="ts">
import { useGameStore } from '@/stores/GameStore';
const gameStore = useGameStore();

function getUnansweredCount() {
    if (gameStore.gameState.player.answerResults == undefined)
        return gameStore.gameState.player.wordPool.length;
    return gameStore.gameState.player.wordPool.length - gameStore.gameState.player.answerResults.length;
}

</script>

<template>
    <div class="container">
        <div class="middle-container">
            <div v-for="answer in gameStore.gameState.player.answerResults" class="answer border-2 border-stone-700"
            :class="{ 'correct': answer.isCorrect, 'wrong': !answer.isCorrect }"
            ></div>
            <!--<div v-for="n in getUnansweredCount()" class="answer border-2 border-stone-700"></div>-->
        </div>
    </div>
</template>

<style scoped>
    .container {
        position: absolute;
        width: 100%;
        height: 105%;
        left: 0;
        top: 0;
        pointer-events: none;
        display: flex;
        justify-content: center;
        align-items: flex-end;
    }
    .middle-container {
        margin-top: -20px;
        width: 100%;
        height: 25px;
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 10px;
    }
    .answer {
        width: 25px;
        height: 25px;
        border-radius: 100px;
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