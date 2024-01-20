import { useState, useEffect } from "react"

interface TimerProps {
  end: Date
}

function calcTimeLeft(end: Date) {
  return Math.floor((+end - +new Date()) / 1000)
}

export default function Timer({ end }: TimerProps) {
  const [timeLeft, setTimeLeft] = useState(calcTimeLeft(end))

  useEffect(() => {
    setTimeLeft(calcTimeLeft(end))
  }, [end])

  useEffect(() => {
    timeLeft > 0 &&
      setTimeout(() => {
        setTimeLeft(calcTimeLeft(end))
      }, 1000)
  }, [timeLeft])

  return <>{timeLeft}s</>
}
