import { GameState } from "./features/game/gameSlice"

// Type definitions for wasm bindings
interface WasmGlobal extends WebAssembly.Instance {
    exports: WasmExports
}

interface WasmExports extends WebAssembly.Exports {
    NewGame(): void
}

declare global {
    var wasm: WasmGlobal

    function OnGameState(s: GameState)
    function OnGameError(s: string)
    function OnPlayerID(s: string)

    function NewGame()
    function JoinMatchmaking()
    function JoinGame(id: string)
    function SubmitGuess(s: string)
}