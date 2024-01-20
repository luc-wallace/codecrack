import { ReactNode } from "react"

interface ButtonProps {
  active: boolean
  onClick: () => void
  children?: ReactNode
}

export default function Button({ active, onClick, children }: ButtonProps) {
  const style = active ? { opacity: 1 } : { opacity: 0.5 }

  return (
    <button style={style} onClick={active ? onClick : () => {}}>
      {children}
    </button>
  )
}
