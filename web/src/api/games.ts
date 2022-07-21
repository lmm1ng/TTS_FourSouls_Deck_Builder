import { IGame, IGameCreatePayload } from '@/models/Game'

export default {
    gameList (): Promise<Array<IGame>> {
        return fetch('/games').then(response => response.json().then(out => out.data))
    },
    createGame (data: IGameCreatePayload) {
        return fetch('games',{
            method: 'POST',
            body: JSON.stringify(data)
        })
    }
}
