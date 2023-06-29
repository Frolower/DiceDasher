import React, { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { setRoomId, setPlayerId } from "./store/gameSlice";
import { useSelector, useDispatch } from "react-redux";
export default () => {
  const navigate = useNavigate();
  const gameInfo = useSelector((state) => state.game.game);
  const dispatch = useDispatch();
  const [input, setInput] = useState("");
  useEffect(() => {}, []);
  const join = () => {
    fetch("http://localhost:8080/joinroom?room_id=" + input)
      .then((response) => response.json())
      .then((r) => {
        if (Object.keys(r).includes("room_id")) {
          dispatch(setRoomId({ room_id: r.room_id }));
          dispatch(setPlayerId({ player_id: r.player_id }));
          navigate("/room");
        } else {
          console.log("error");
        }
      });
  };
  return (
    <div>
      <input value={input} onInput={(e) => setInput(e.target.value)} />
      <button onClick={join}>join</button>
    </div>
  );
};
