import { DndContext, DragEndEvent } from "@dnd-kit/core"
import { useEffect, useState } from "react"
import { Peg, Piece, PieceDragRow, PieceSlot } from "./Piece"
import { BoardRow, Colour } from "../types/mastermind"
import { OpIn, OpOut, Packet, ParamsSubmit } from "../types/websocket"
import "./Game.css"
import Timer from "./Timer"
import Button from "./Button"

interface GameProps {
  ws: WebSocket
  opponent: string
  opCombi: number[]
}

const blankRow = [
  Colour.None,
  Colour.None,
  Colour.None,
  Colour.None,
  Colour.None,
]

function reverse<T>(a: T[]) {
  return a.map((_, idx) => a[a.length - 1 - idx])
}

function futureDate(seconds: number) {
  return new Date(new Date().getTime() + seconds * 1000)
}

export default function Game({ ws, opponent, opCombi }: GameProps) {
  const [board, setBoard] = useState<BoardRow[]>([])
  const [combi, setCombi] = useState<Colour[]>(blankRow)
  const [round, setRound] = useState(1)
  const [submitted, setSubmitted] = useState(false)
  const [oppStatus, setOppStatus] = useState<Colour[][]>([])
  const [gameStatus, setGameStatus] = useState("")
  const [timerEnd, setTimerEnd] = useState<Date>(futureDate(30))

  useEffect(() => {
    ws.addEventListener("message", (e) => {
      const packet: Packet<any> = JSON.parse(e.data)

      switch (packet.op) {
        case OpIn.RoundStart:
          setSubmitted(false)
          setRound(packet.d.roundNum)
          setTimerEnd(futureDate(30))
          break
        case OpIn.BoardUpdate:
          setBoard(packet.d)
          break
        case OpIn.OpponentStatus:
          setOppStatus((oppStatus) => [...oppStatus, packet.d.status])
          break
        case OpIn.GameEnd:
          if (packet.d.win === true) {
            setGameStatus("won!")
          } else if (packet.d.win === false) {
            setGameStatus("lost.")
          } else {
            setGameStatus("drew.")
          }
          break
        case OpIn.ForceSubmit:
          submitCombi()
      }
    })
  }, [])

  function onDragEnd(event: DragEndEvent) {
    if (event.over) {
      const slotIndex = Number.parseInt(event.over.id as string)
      const colourId = Number.parseInt(event.active.id as string)
      setCombi(combi.map((c, i) => (i === slotIndex ? colourId : c)))
    }
  }

  function submitCombi() {
    setCombi((combi) => {
      const packet: Packet<ParamsSubmit> = {
        op: OpOut.Submit,
        d: { combi: combi },
      }
      ws.send(JSON.stringify(packet))
      return combi
    })
    setCombi(blankRow)
    setSubmitted(true)
  }

  return (
    <DndContext onDragEnd={onDragEnd}>
      <PieceDragRow />
      <hr />
      <div className="player-status">
        <div className="player-status-row">
          <h2>
            {gameStatus == "" ? (
              <>
                round {round} - <Timer end={timerEnd} /> left
              </>
            ) : (
              `you ${gameStatus}`
            )}
          </h2>
          <div className="tooltip">
            i
            <span className="tooltip-text">
              <span>black peg:</span> right colour and in the right position
              <br />
              <span>white peg:</span> right colour but in the wrong position
            </span>
          </div>
        </div>
        <div className="row">
          <div className="piece-slots">
            {[...Array(5)].map((_, n) => (
              <PieceSlot key={n} id={n.toString()} displayCol={combi[n]} />
            ))}
          </div>
          {gameStatus == "" ? (
            !submitted ? (
              <Button active={combi.every((c) => c > 0)} onClick={submitCombi}>
                submit
              </Button>
            ) : (
              <p>submitted</p>
            )
          ) : (
            <button onClick={() => location.reload()}>play again</button>
          )}
        </div>
        <div className="board">
          {reverse(board).map((row, i) => (
            <div key={i} className="row">
              {row.pieces.map((piece, i) => (
                <Piece key={i} colour={piece} />
              ))}
              <div className="row-status">
                <h2>{String(board.length - i).padStart(2, "0")}</h2>
                {row.check
                  .filter((p) => p > 0)
                  .sort((a, b) => b - a)
                  .map((piece, i) => (
                    <Peg key={i} colour={piece} opponent={false} />
                  ))}
              </div>
            </div>
          ))}
        </div>
      </div>
      <hr />
      <div className="opponent-status">
        <h2>
          opponent <span className="subtext">- {opponent}</span>
        </h2>
        <div className="row">
          {opCombi.map((piece, i) => (
            <Piece key={i} colour={piece} />
          ))}
        </div>
        <div className="opponent-pegs">
          {reverse(oppStatus).map((row, i) => (
            <div key={i} className="row">
              {row
                .filter((p) => p > 0)
                .sort((a, b) => b - a)
                .map((c, i) => (
                  <Peg key={i} colour={c} opponent={true} />
                ))}
            </div>
          ))}
        </div>
      </div>
    </DndContext>
  )
}
