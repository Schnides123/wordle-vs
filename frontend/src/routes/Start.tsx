import { useEffect } from "react"
import { useSelector, useDispatch } from "react-redux"
import { GameState, selectGame, update } from "../features/game/gameSlice"
import { Navigate, useLocation } from 'react-router-dom'

function Start() {
    const game = useSelector(selectGame)
    const location = useLocation()

    useEffect(() => {
        NewGame()
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

export default Start
