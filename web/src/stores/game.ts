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
        fetchCreateGame (game: IGameCreatePayload): Promise<any> {
            this.isGameCreating = true
            return gamesApi.createGame(game).finally(() => this.isGameCreating = false)
        },
        fetchGames (): Promise<any> {
            return gamesApi.gameList().then((games) => {
                this.games = games
            })
        }
    }
})
