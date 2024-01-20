import { DndContext, DragEndEvent } from "@dnd-kit/core"
import "./App.css"
import { useState, useRef, useReducer } from "react"
import { PieceDragRow, PieceSlot } from "./components/Piece"
import { Colour } from "./types/mastermind"
import Game from "./components/Game"
import { FindGameParams, OpIn, OpOut, Packet } from "./types/websocket"
import Symbol from "./crack_symbol.svg"
import Button from "./components/Button"

function App() {
  const [foundGame, setFoundGame] = useState(false)
  const [matchmaking, setMatchmaking] = useState(false)
  const ws = useRef<WebSocket | null>(null)
  const [opponent, setOpponent] = useState("")

  const [name, setName] = useReducer((_prev: string, cur: string) => {
    localStorage.setItem("name", cur)
    return cur
  }, localStorage.getItem("name") ?? "")

  const [combi, setCombi] = useState<Colour[]>([
    Colour.None,
    Colour.None,
    Colour.None,
    Colour.None,
    Colour.None,
  ])

  function onDragEnd(event: DragEndEvent) {
    if (event.over) {
      const slotIndex = Number.parseInt(event.over.id as string)
      const colourId = Number.parseInt(event.active.id as string)
      setCombi(combi.map((c, i) => (i === slotIndex ? colourId : c)))
    }
  }

  function matchmake() {
    setMatchmaking(true)

    let wsUrl = ""
    if (import.meta.env.DEV) {
      wsUrl = import.meta.env.VITE_WS_URL ?? "ws://localhost:9999/ws"
    } else {
      wsUrl = `ws://${window.location.host}/ws`
    }

    ws.current = new WebSocket(wsUrl)

    ws.current.addEventListener("open", () => {
      const packet: Packet<FindGameParams> = {
        op: OpOut.FindGame,
        d: { username: name, combi },
      }
      ws.current?.send(JSON.stringify(packet))
    })

    ws.current.addEventListener("message", (e) => {
      const packet: Packet<any> = JSON.parse(e.data)
      switch (packet.op) {
        case OpIn.Matchmaking:
          setMatchmaking(true)
          break
        case OpIn.GameStart:
          setMatchmaking(false)
          setFoundGame(true)
          setOpponent(packet.d.opponent)
          ws.current?.send(JSON.stringify({ op: OpOut.Pong }))
      }
    })
  }

  return (
    <>
      <div className="logo">
        <h1>code</h1>
        <img src={Symbol}></img>
        <h1>crack</h1>
      </div>
      <div className="game-container">
        {foundGame ? (
          <Game ws={ws.current!} opponent={opponent} opCombi={combi} />
        ) : (
          <DndContext onDragEnd={onDragEnd}>
            <PieceDragRow />
            <hr />
            <div className="game-form col">
              <h2>welcome player</h2>
              <h3>choose a username:</h3>
              <input
                value={name}
                spellCheck={false}
                placeholder="type here..."
                onChange={(e) => setName(e.target.value)}
              />
              <h3>choose your code:</h3>
              <div className="piece-slots">
                {[...Array(5)].map((_, n) => (
                  <PieceSlot key={n} id={n.toString()} displayCol={combi[n]} />
                ))}
              </div>
              {matchmaking ? (
                <h3>finding game...</h3>
              ) : (
                <Button
                  active={combi.every((c) => c > 0) && name.length > 0}
                  onClick={matchmake}>
                  find game
                </Button>
              )}
            </div>
          </DndContext>
        )}
      </div>
    </>
  )
}

export default App
