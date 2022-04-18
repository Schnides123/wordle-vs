import { createAsyncThunk, createSlice, PayloadAction } from '@reduxjs/toolkit';
import { RootState, store } from '../../app/store'

export interface SessionState {
    playerID: string | null
}

const initialState: SessionState = {
    playerID: null,
}

export const sessionSlice = createSlice({
    name: 'session',
    initialState,
    reducers: {
        updatePlayerID: function (state, action: PayloadAction<string>) {
            state.playerID = action.payload
        },
    }
})
export const { updatePlayerID } = sessionSlice.actions
export const selectSession = (state: RootState) => state.session

export default sessionSlice.reducer;
