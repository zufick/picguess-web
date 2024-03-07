<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useWebsocketStore } from '@/stores/WebsocketStore';
import { useRoute } from 'vue-router';


const wsStore = useWebsocketStore();
const enteredName = ref("name");
const enteredRoomId = ref("");
const route = useRoute()


function setRandomName() {
    let words = ['Personal', 'Flour', 'Discover', 'Marble', 'Doggy', 'Bubble', 'Style', 'Carpet', 'Movie' ]
    enteredName.value = words[Math.floor(Math.random()*words.length)] + words[Math.floor(Math.random()*words.length)] + Math.floor(Math.random() * 1000);
}

setRandomName();


function checkEmptyName() {
    if (enteredName.value.length < 1) {
        setRandomName();
    }
}

function joinRoomClick() {
    wsStore.joinRoom(enteredRoomId.value, {username: enteredName.value})
}

function createRoomClick() {
    wsStore.createRoom({username: enteredName.value})
}

onMounted(() => {
    console.log("on mounted")
    if (route.params.room != undefined) {
        enteredRoomId.value = route.params.room + ""
        joinRoomClick();
    }
})

</script>

<template>
    <div class="mt-6 flex items-start gap-8">
        <div class="flex flex-col items-end p-4 border-2 rounded-lg border-white/20 bg-stone-800 w-96 fade-Anim">
            <p class="text-4xl mb-2 text-right">Welcome,<br><b class="break-all">{{ enteredName }}</b></p>
            <input type="text" class="text-right" id="nameInput" v-model="enteredName" @blur="checkEmptyName" placeholder="Enter nickname" maxlength="18" minlength="1">
        </div>
        <div class="p-4">
            <div class="mb-6">
                <label for="joinroomInput">Join room by id</label>
                <div class="flex items-start">
                    <input id="joinroomInput" v-model="enteredRoomId" type="text" placeholder="Room id">
                    <button @click="joinRoomClick">Join</button>
                </div>
            </div>
            <div class="flex justify-between">
                <button @click="createRoomClick">Create room</button>
                <button>Matchmaking</button>
            </div>
        </div>
    </div>
</template>

<style>
    @keyframes gradient {
        0% {
            background-position: 0% 50%;
        }
        50% {
            background-position: 100% 50%;
        }
        100% {
            background-position: 0% 50%;
        }
    }

    .fade-Anim {
        background: linear-gradient(-45deg, #2d107500, #2d107533, #10757500, #10757533, #10753200, #10753233);
        background-size: 400% 400%;
        animation: gradient 15s ease infinite;
    }
</style>