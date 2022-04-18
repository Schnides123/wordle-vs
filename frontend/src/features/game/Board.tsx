import { Box, Grid, Stack } from "@mui/material"
import GuessCard from "./GuessCard"
import WordCard from "./WordCard"
import GameInput from "./GameInput"
import { Guess, GameWord, getGuessedLetters } from "./gameSlice"
import React, { useState } from "react"

interface BoardProps {
	word: GameWord
	guessed: Guess[]
	winner: boolean
}

class Board extends React.Component<BoardProps> {
	state = { guess: null as any as string[], first: React.createRef<HTMLTextAreaElement>() }

	constructor(props: BoardProps) {
		super(props)
		if (this.state.guess == null || this.state.guess.length != props.word.length) {
			this.state.guess = Array<string>(props.word.length).fill("")
		}
	}

	render() {
		const { word, guessed, winner } = this.props

		let inputs: JSX.Element[] = []
		for (let i = 0; i < word.length; i++) {
			inputs.push(
				<div key={i}>
					<textarea
						ref={i == 0 ? this.state.first : undefined}
						value={this.state.guess[i]}
						autoFocus={i == 0}
						rows={1}
						cols={1}
						maxLength={1}
						id={`guess-${i}`}
						style={{
							resize: "none",
							fontSize: 45,
							textAlign: "center",
						}}
						onChange={(e) => {
							const { maxLength, value, id } = e.target
							const [fieldName, fieldIndex] = id.split("-")

							console.log(value)
							let fieldIntIndex = i
							this.state.guess[i] = value.toUpperCase()

							// Check if no of char in field == maxlength
							if (value.length >= 1) {
								// It should not be last input field
								if (fieldIntIndex < word.length) {
									// Get the next input field using it's name
									const nextfield = document.querySelector(
										`textarea#guess-${fieldIntIndex + 1}`
									)
									this.state.guess[fieldIntIndex + 1] = ""

									// If found, focus the next field
									if (nextfield !== null) {
										(nextfield as any).focus()
									}
								}
							}
							this.setState(this.state)
						}}
						onKeyPress={(e) => {
							if (e.key === "Enter" && !e.shiftKey) {
								e.preventDefault()
								// Collect submission
								let s = ""
								for (let i = 0; i < word.length; i++) {
									if (this.state.guess[i].length == 0) {
									} else {
										s = s + this.state.guess[i]
									}
								}

								if (s.length === word.length) {
									SubmitGuess(s)
									this.state.guess = Array(word.length).fill("")
									this.setState(this.state)
									this.state.first.current?.focus()
								} else {
									console.error("a field is empty!")
								}
							}
						}}
						onKeyDown={(e) => {
							if (e.key === "0" ||
								e.key === "1" ||
								e.key === "2" ||
								e.key === "3" ||
								e.key === "4" ||
								e.key === "5" ||
								e.key === "6" ||
								e.key === "7" ||
								e.key === "8" ||
								e.key === "9") {
								console.log("intercept")
								e.preventDefault()
								return false
							}

							if (
								i > 0 &&
								e.key == "Backspace" &&
								e.currentTarget.value.length == 0
							) {
								e.preventDefault()

								const nextfield = document.querySelector(
									`textarea#guess-${i - 1}`
								)

								// If found, focus the next field
								if (nextfield !== null) {
									this.state.guess[i - 1] = ""
									this.setState(this.state);
									(nextfield as any).focus()
								}
							} else if (
								e.key !== "Enter" &&
								!e.shiftKey &&
								e.key.length == 1 &&
								!e.metaKey &&
								!e.altKey &&
								!e.ctrlKey
							) {
								e.currentTarget.value = ""
							}
						}}
					/>
				</div>
			)
		}

		return (
			<Stack direction="column" maxHeight="100%" alignItems="center">
				<WordCard style={{ flex: "0 1 auto" }} word={getGuessedLetters(word, guessed)} />
				<div style={{ flex: "1 1 100%", overflow: "auto", display: "flex", flexDirection: "column-reverse" }}>
					{guessed.slice(0).reverse().map((guess, index) => (
						<GuessCard key={guessed.length - index} {...guess} />
					))}
				</div>
				<Stack style={{ flex: "0 1 auto" }} direction="row">{inputs}</Stack>
			</Stack>
		)
	}
}

export default Board
