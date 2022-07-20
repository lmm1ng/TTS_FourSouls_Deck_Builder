export interface IGame {
    id: string,
    name: string,
    description: string,
    image: string,
    createdAt: string | null,
    updatedAt?: string | null,
}

export interface IGameCreatePayload {
    name: string,
    description: string,
    image: string,
}
