import React from "react"
import { Routes, Route } from "react-router-dom"
import { Button, Container, Stack } from "@mui/material"
import Game from "./features/game/Game"

import Welcome from "./routes/Welcome"
import Play from "./routes/Play"
import Matchmaking from "./routes/Matchmaking"
import Start from "./routes/Start"
import { store } from "./app/store"
import { GameState, update } from "./features/game/gameSlice"
import { updatePlayerID } from "./features/session/sessionSlice"

declare class Go { }

interface AppState {
	wasmInit: boolean
}

global.OnGameError = (s: string) => {
	alert(s)
}

global.OnGameState = (s: GameState) => {
	console.log(s)
	store.dispatch(update(s))
}

global.OnPlayerID = (id: string) => {
	console.log("new player id", id)
	store.dispatch(updatePlayerID(id))
}

class App extends React.Component<{}, AppState> {
	state = {
		wasmInit: false,
	};

	componentDidMount() {
		console.log("Begin loading wasm...")

		// polyfill
		if (!WebAssembly.instantiateStreaming) {
			WebAssembly.instantiateStreaming = async (resp, importObject) => {
				const source = await (await resp).arrayBuffer()
				return await WebAssembly.instantiate(source, importObject)
			}
		}

		const go: any = new Go()
		WebAssembly.instantiateStreaming(fetch("/main.wasm"), go.importObject).then(
			(result) => {
				console.log("Loaded wasm")
				go.run(result.instance)
				window.wasm = result.instance as any
				this.setState({ wasmInit: true })
			}
		)
	}

	render() {
		if (!this.state.wasmInit) {
			return <h3>Loading...</h3>
		}

		return (
			<Container style={{ height: "100vh", padding: "1em" }}>
				<Routes>
					<Route path="/" element={<Welcome />} />
					<Route path="start" element={<Start />} />
					<Route path="play" element={<Matchmaking />} />
					<Route path="play/:id" element={<Play />} />
				</Routes>
			</Container>
		)
	}
}

export default App
