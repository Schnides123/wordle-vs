import { useEffect } from "react"
import { useSelector } from "react-redux"
import { useParams, Navigate, useLocation } from "react-router-dom"
import { selectGame } from "../features/game/gameSlice"

function Matchmaking() {
    const location = useLocation()
    const game = useSelector(selectGame)
    let { id } = useParams() // Unpacking and retrieve id

    useEffect(() => {
        JoinMatchmaking()
    }, [])

    if (!game) {
        return <h1>Loading</h1>
    }

    return <Navigate
        to={`/play/${game.id}`}
        replace
        state={{ location }}
    />
}

export default Matchmaking