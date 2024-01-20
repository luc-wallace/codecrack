export interface Packet<T> {
  op: Opcode
  d: T
}

export type Opcode = OpIn | OpOut

export enum OpIn {
  Ping = 0,
  Matchmaking,
  GameStart,
  GameEnd,
  RoundStart,
  RoundEnd,
  BoardUpdate,
  OpponentStatus,
  ForceSubmit
}

export enum OpOut {
  Pong = 0,
  FindGame,
  Submit,
}

export interface FindGameParams {
  username: string
  combi: number[]
}

export interface ParamsSubmit {
  combi: number[]
}

export interface ParamsRoundEnd {
  combi: number[]
}
