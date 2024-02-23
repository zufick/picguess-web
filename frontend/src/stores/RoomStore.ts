import { defineStore } from 'pinia'

export const useRoomStore = defineStore('room', {
    state: () => ({ joinRoomId: undefined }),
    actions: {
      joinRoom() {
        alert(this.joinRoomId)
      },
    },
  })