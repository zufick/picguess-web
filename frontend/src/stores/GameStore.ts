import { defineStore } from 'pinia'


export const useGameStore = defineStore('game', {
    state: () => ({ 
        gameState: {}
    }),
  })
  