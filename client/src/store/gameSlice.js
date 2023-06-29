import { createSlice } from "@reduxjs/toolkit";

export const gameSlice = createSlice({
  name: "game",
  initialState: {
    player_id: -1,
    room_id: -1,
  },
  reducers: {
    setRoomId: (state, action) => {
      state.room_id = action.payload.room_id;
    },
    setPlayerId: (state, action) => {
      state.player_id = action.payload.player_id;
    },
  },
});

export const { setRoomId, setPlayerId } = gameSlice.actions;

export default gameSlice.reducer;
