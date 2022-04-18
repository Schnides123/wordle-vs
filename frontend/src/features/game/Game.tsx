import { Stack, Grid, Container } from '@mui/material'
import PlayerCard from './PlayerCard'
import Board from './Board'
import { GameState, isTurn } from './gameSlice'


function Game({ game }: { game: GameState }) {
	return (
		<Stack direction="row" spacing={2} alignItems="flex-start" justifyContent="center" height="100%">
			<PlayerCard
				player={game.players[0]}
				isTurn={isTurn(game, game.players[0])}
			/>
			<Board word={game.word} guessed={game.guessed} winner={game.winner} />
			<PlayerCard
				player={game.players[1]}
				isTurn={isTurn(game, game.players[1])}
			/>
		</Stack>
	)
}

export default Game
