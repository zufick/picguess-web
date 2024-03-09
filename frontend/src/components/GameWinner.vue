<script setup lang="ts">
import { useWebsocketStore } from '@/stores/WebsocketStore';

const props = defineProps(['winner'])

const wsStore = useWebsocketStore();

</script>

<template>
    <div class="window">
        <div>
            <p class="text-lg">Game ended!</p>
            <p class="text-5xl">Winner:</p>
            <p class="text-5xl">
                <b class="flex items-center gap-8 text-yellow-300">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 18.75h-9m9 0a3 3 0 0 1 3 3h-15a3 3 0 0 1 3-3m9 0v-3.375c0-.621-.503-1.125-1.125-1.125h-.871M7.5 18.75v-3.375c0-.621.504-1.125 1.125-1.125h.872m5.007 0H9.497m5.007 0a7.454 7.454 0 0 1-.982-3.172M9.497 14.25a7.454 7.454 0 0 0 .981-3.172M5.25 4.236c-.982.143-1.954.317-2.916.52A6.003 6.003 0 0 0 7.73 9.728M5.25 4.236V4.5c0 2.108.966 3.99 2.48 5.228M5.25 4.236V2.721C7.456 2.41 9.71 2.25 12 2.25c2.291 0 4.545.16 6.75.47v1.516M7.73 9.728a6.726 6.726 0 0 0 2.748 1.35m8.272-6.842V4.5c0 2.108-.966 3.99-2.48 5.228m2.48-5.492a46.32 46.32 0 0 1 2.916.52 6.003 6.003 0 0 1-5.395 4.972m0 0a6.726 6.726 0 0 1-2.749 1.35m0 0a6.772 6.772 0 0 1-3.044 0" />
                    </svg>
                    {{  props.winner.user_info.username  }}
                </b>
                <p>Score: <span>{{  props.winner.score  }}</span></p>
            </p>
            <button @click="wsStore.startGame()" class="mt-8">Restart</button>
        </div>
    </div>
</template>

<style scoped>
    .window {
        animation: 1s ease appear;
        position: absolute;
        top: 0;
        left: 0;
        z-index: 50;
        width: 100%;
        height: 100%;
        border: 2px solid rgb(68 64 60);
        border-radius: 10px;
        color: rgb(68 64 60);
        background: rgb(28 25 23);
        box-shadow: 0px 20px 20px #0000003d;
        display: flex;
        justify-content: center;
        align-items: center;
        text-align: center;
        color: white;
    }
</style>