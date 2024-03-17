<script setup lang="ts">
import { useWebsocketStore } from '@/stores/WebsocketStore';

const wsStore = useWebsocketStore();

    function drawClearVoteClicked(voted: boolean) {
        //virtualCanvas.drawClear("local");
        wsStore.sendDrawClear(voted);
    }

    function drawClearCancelClicked() {
        wsStore.sendDrawClear(false);
    }

    function isVoteSent() {
        for (let i in wsStore.roomState.vote_clearcanvas) {
            if (i == wsStore.localUserId) {
                return true;
            }
        }
        return false;
    }

    function getOwnVoteResult() {
        for (let i in wsStore.roomState.vote_clearcanvas) {
            if (i == wsStore.localUserId) {
                return wsStore.roomState.vote_clearcanvas[i];
            }
        }
    }
</script>


<template>
    <div>
        <button v-if="!Object.keys(wsStore.roomState.vote_clearcanvas).length" @click="drawClearVoteClicked(true)">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
                </svg>
            </button>
            <div v-else class="flex flex-col items-center gap-4 appear-anim rounded-md border-2 border-stone-700 py-2 bg-stone-900 vote-anim" >
                <p style="width: 60px; text-align: center;     font-size: .8rem;" class="break-words">Vote: Clear canvas?</p>
                <button class="danger-btn-onhover w-6 flex flex-col items-center justify-center" @click="drawClearVoteClicked(true)" :class="{'danger-btn': getOwnVoteResult()}">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                    <path stroke-linecap="round" stroke-linejoin="round" d="m4.5 12.75 6 6 9-13.5" />
                    </svg>
                </button>

                <button class="success-btn-onhover w-6 flex flex-col items-center justify-center" @click="drawClearVoteClicked(false)" :class="{'success-btn': isVoteSent() ? !getOwnVoteResult() : false}">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-4 h-4">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18 18 6M6 6l12 12" />
                    </svg>
                </button>

            </div>
    </div>
</template>

<style scoped>
    @keyframes vote {
        0% {
            background: rgba(255, 255, 255, 0.137);
            transform: scale(1.1);
        }

        20% {
            transform: scale(1);
        }
    }


    .vote-anim {
        animation: 2s ease vote infinite;
    }
</style>

<style>

</style>