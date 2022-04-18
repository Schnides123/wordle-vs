import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit"
import { RootState, AppThunk, store } from "../../app/store"

export interface Guess {
    word: string
    result: number[]
    player: string
    timestamp: string
}

export interface GameWord {
    length: number
}

export interface GameState {
	guessed: Guess[]
    id: string
    options: object
    players: string[]
    winner: boolean
    word: GameWord
}

export const getLastGuess = (state: GameState) => {
	return state.guessed[state.guessed.length - 1];
};

export const isTurn = (state: GameState, player: String) => {
	if (state.guessed.length === 0) {
		return true;
	}
	const timeSinceLastGuess =
		Date.now() - Date.parse(getLastGuess(state).timestamp);
	const lastGuess = getLastGuess(state);
	return (
		!state.winner &&
		state.players.length > 1 &&
		(lastGuess.player !== player || timeSinceLastGuess >= 30000)
	);
};

export const getGuessedLetters = (word: GameWord, guessed: Guess[]) => {
	let out = "";
	for (let i = 0; i < word.length; i++) {
		out += checkLetter(i, guessed);
	}
	return out
};

const checkLetter = (index: number, guessed: Guess[]) => {
	for (let j = 0; j < guessed.length; j++) {
		if (guessed[j].result[index] === 2) {
			return guessed[j].word[index];
		}
	}
	return "⬛";
};

export const canSteal = (state: GameState) => {
	if (state.guessed.length === 0) {
		return true;
	}
	return Date.now() - Date.parse(getLastGuess(state).timestamp) >= 30000;
};

const initialState: GameState | null = null

export const gameSlice = createSlice({
    name: "session",
    initialState: initialState as GameState | null,
    reducers: {
        update: function (state, action: PayloadAction<GameState>) {
            return action.payload
        },
    },
})

export const { update } = gameSlice.actions
export const selectGame = (state: RootState) => state.game

export default gameSlice.reducer
