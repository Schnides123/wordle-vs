import { useEffect } from "react"
import { useSelector } from "react-redux"
import { useParams } from "react-router-dom"
import Game from "../features/game/Game"
import { selectGame } from "../features/game/gameSlice"

function Play() {
    const game = useSelector(selectGame)
    let { id } = useParams() // Unpacking and retrieve id

    useEffect(() => {
        JoinGame(id as string)
    }, [id])

    if (!game || id !== game.id) {
        return <h1>Loading...</h1>
    }

    return <Game game={game} />
}

export default Play