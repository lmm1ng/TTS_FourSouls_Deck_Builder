import { defineStore } from 'pinia'
import { IGame, IGameCreatePayload } from '@/models/Game'
import gamesApi from '@/api/games'

type GameState = {
    games: Array<IGame>,
    isGameCreating: boolean
}

export default defineStore({
    id: 'gameStore',
    state: () => ({ games: [], isGameCreating: false } as GameState),
    actions: {
        fetchCreateGame (game: IGameCreatePayload) {
            this.isGameCreating = true
            gamesApi.createGame(game).finally(() => this.isGameCreating = false)
        },
        fetchGames () {
            gamesApi.gameList().then((response) => {
                this.games = response
                console.log(this.games)
            })
            // console.log(this.games)
        }
    }
})
