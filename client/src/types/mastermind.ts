export enum Colour {
  None = 0,
  White,
  Black,
  Red,
  Green,
  Blue,
  Yellow,
  Orange,
  Brown,
}

export interface BoardRow {
  pieces: Colour[]
  check: Colour[]
}
