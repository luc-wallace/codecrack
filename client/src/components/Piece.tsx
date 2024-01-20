import { useDroppable, useDraggable } from "@dnd-kit/core"
import { ReactNode } from "react"
import { Colour } from "../types/mastermind"

interface PieceSlotProps {
  id: string
  children?: ReactNode
  displayCol: number
}

export function PieceSlot({ id, children, displayCol }: PieceSlotProps) {
  const { isOver, setNodeRef } = useDroppable({ id })

  return (
    <div
      ref={setNodeRef}
      className="piece"
      style={{
        backgroundColor: `var(--col-${displayCol})`,
        opacity: isOver ? 0.2 : 1,
      }}>
      {children}
    </div>
  )
}

interface PieceDragProps {
  children?: ReactNode
  id: string
}

export function PieceDragRow() {
  return (
    <div>
      <h2>colours</h2>
      <div className="colour-container">
        {[...Array(8)].map((_, n) => (
          <PieceDrag key={n} id={(n + 1).toString()}>
            <div
              className="piece drag"
              style={{ backgroundColor: `var(--col-${n + 1})` }}></div>
          </PieceDrag>
        ))}
      </div>
    </div>
  )
}

function PieceDrag({ children, id }: PieceDragProps) {
  const { attributes, listeners, setNodeRef, transform } = useDraggable({ id })
  const style = transform
    ? {
        transform: `translate3d(${transform.x}px, ${transform.y}px, 0)`,
      }
    : undefined

  return (
    <div ref={setNodeRef} style={style} {...listeners} {...attributes}>
      {children}
    </div>
  )
}

interface PieceProps {
  colour: Colour
}

export function Piece({ colour }: PieceProps) {
  return (
    <div
      className="piece"
      style={{ backgroundColor: `var(--col-${colour})` }}></div>
  )
}

interface PegProps {
  colour: Colour
  opponent: boolean
}

export function Peg({ colour, opponent }: PegProps) {
  return (
    <div
      className={`peg ${opponent ? "opponent" : ""}`}
      style={{
        backgroundColor: `var(--col-${colour})`,
      }}></div>
  )
}
