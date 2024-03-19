<script setup lang="ts">
const props = defineProps(["size", "width"])


</script>

<template>
    <svg v-if="props.size == 'small'" :width="props.width" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 18 18" fill="currentColor">
        <g>
            <path class="cls-1" d="m9,18C4.04,18,0,13.96,0,9S4.04,0,9,0s9,4.04,9,9-4.04,9-9,9Zm0-16.5C4.86,1.5,1.5,4.86,1.5,9s3.36,7.5,7.5,7.5,7.5-3.36,7.5-7.5S13.14,1.5,9,1.5Z"/>
            <circle class="cls-1" cx="9" cy="9" r="3"/>
        </g>
    </svg>
    <svg v-if="props.size == 'medium'" :width="props.width" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 18 18" fill="currentColor">
        <g>
            <path class="cls-1" d="m9,18C4.04,18,0,13.96,0,9S4.04,0,9,0s9,4.04,9,9-4.04,9-9,9Zm0-16.5C4.86,1.5,1.5,4.86,1.5,9s3.36,7.5,7.5,7.5,7.5-3.36,7.5-7.5S13.14,1.5,9,1.5Z"/>
            <circle class="cls-1" cx="9" cy="9" r="4.5"/>
        </g>
    </svg>
    <svg v-if="props.size == 'large'" :width="props.width" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 18 18" fill="currentColor">
        <g>
            <path class="cls-1" d="m9,18C4.04,18,0,13.96,0,9S4.04,0,9,0s9,4.04,9,9-4.04,9-9,9Zm0-16.5C4.86,1.5,1.5,4.86,1.5,9s3.36,7.5,7.5,7.5,7.5-3.36,7.5-7.5S13.14,1.5,9,1.5Z"/>
            <circle class="cls-1" cx="9" cy="9" r="6"/>
        </g>
    </svg>
</template>

<style>

</style>