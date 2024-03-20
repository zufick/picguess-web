<script setup lang="ts">
import {ref} from 'vue'
import BrushIcon from './BrushIcon.vue';

const menuOpened = ref(false)
const model = defineModel({default: 25})

let brushSizes: { [key: string]: number; }  = {
    'small': 10,
    'medium': 25,
    'large': 50
}

function selectBrush(size: number) {
    model.value = size
    menuOpened.value = false
}

function getSizeName(size: number) {
    let name:string = "";
    Object.keys(brushSizes).forEach((brushSize: string) => {
        if (brushSizes[brushSize] == size) {
            name = brushSize
        }
    })
    return name;
}

</script>

<template>
        <div>
        <button v-if="!menuOpened" @click="menuOpened = true">
            <BrushIcon :size="getSizeName(model)" width="20"/>
        </button>
            <div v-else class="flex flex-col items-center gap-4 appear-anim rounded-md border-2 border-stone-700 p-2 bg-stone-900 vote-anim" >
                <button v-for="item in brushSizes" @click="selectBrush(item)" class="w-6 flex flex-col items-center justify-center" :class="{'success-btn': model == item }">
                    <BrushIcon :size="getSizeName(item)" width="20"/>
                </button>
            </div>
    </div>
</template>

<style scoped>

</style>