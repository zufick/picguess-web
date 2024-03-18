<script setup lang="ts">
import { ref, watch } from 'vue';
import { useWebsocketStore } from '../stores/WebsocketStore';
const wsStore = useWebsocketStore();
const props = defineProps(['user', 'score'])

const scoreGained = ref(false)
const scoreLost = ref(false)
const scoreDiff = ref(0)
const scoreDiffAnim = ref(false)


function isLocalUser(id: string) {
  return wsStore.localUserId == id
}

watch(() => props.score, (newValue: string, oldValue: string) => {
    scoreGained.value = parseInt(newValue) > parseInt(oldValue)
    scoreLost.value = parseInt(newValue) < parseInt(oldValue)
    scoreDiff.value = parseInt(newValue) - parseInt(oldValue)

    if (newValue != oldValue) {
        scoreDiffAnim.value = true
        setTimeout(() => {
            scoreDiffAnim.value = false
    }, 1000)
    }
})

</script>

<template>
    <div class="flex p-2 gap-1 border-2 rounded-lg  bg-stone-800" 
        :class="{'border-blue-400 text-blue-400 bg-blue-900/25': isLocalUser(user.id)}">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 6a3.75 3.75 0 1 1-7.5 0 3.75 3.75 0 0 1 7.5 0ZM4.501 20.118a7.5 7.5 0 0 1 14.998 0A17.933 17.933 0 0 1 12 21.75c-2.676 0-5.216-.584-7.499-1.632Z" />
        </svg>
        <div>
            <span>{{ user.username }}</span>
            <span class="ml-4 score-text relative" :class="{'score-gained': scoreGained, 'score-lost': scoreLost}">
                Score: {{props.score}}
                <span class="absolute score-diff px-4" :class="{'score-diff-slide-anim' : scoreDiffAnim, 'text-rose-600' : scoreLost, 'text-emerald-400': scoreGained}">{{ scoreDiff }}</span>
            </span>
        </div>
    </div>
</template>

<style scoped>

.list-move, /* apply transition to moving elements */
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateY(30px);
}
.list-leave-active {
  position: absolute;
}

.score-text {
  border: 2px;
    padding: 3px 15px;
    border-radius: 0.5rem;
    background: #ffffff14;
    color: #ffffffa1;
}

.score-diff {
    width: 100%;
    left: 0;
    text-align: right;
    top: -30px;
    opacity: 0;
    pointer-events: none;
}

@keyframes score-gained-anim {
    from {
        color: rgb(52 211 153);
    }
    to {
        color: #ffffffa1;
    }
}

@keyframes score-lost-anim {
    from {
        color: rgb(225 29 72);
    }
    to {
        color: #ffffffa1;
    }
}

.score-gained {
    animation: 1s ease score-gained-anim;
}

.score-lost {
    animation: 1s ease score-lost-anim;
}

@keyframes score-diff-slide {
    0% {
        opacity: 1;
        transform: translateY(0);
    }
    100% {
        opacity: 0;
        transform: translateY(-10px);
    }
}

.score-diff-slide-anim {
    animation: 1s linear score-diff-slide;
}
</style>