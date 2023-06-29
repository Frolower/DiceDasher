import React, { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { setRoomId, setPlayerId } from "./store/gameSlice";
import { useSelector, useDispatch } from "react-redux";
export default () => {
  const navigate = useNavigate();
  const gameInfo = useSelector((state) => state.game.game);
  const dispatch = useDispatch();
  useEffect(() => {
    fetch("http://localhost:8080/newroom")
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
  }, []);
  return <div>create room</div>;
};
