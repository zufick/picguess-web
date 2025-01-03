import type { GameState } from '@/types/GameState'
import { defineStore } from 'pinia'


export const useGameStore = defineStore('game', {
    state: () => ({ 
        gameState: {} as GameState
    }),
  })
  