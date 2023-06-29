import React, { useEffect, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Link, useNavigate } from "react-router-dom";
export default () => {
  const gameInfo = useSelector((state) => state.game);
  const navigate = useNavigate();
  const [dice, setDice] = useState("");
  const socket = new WebSocket(
    "ws://localhost:8080/game?player_id=" + gameInfo.player_id
  );
  useEffect(() => {
    console.log(gameInfo);
    if (gameInfo.room_id == -1) {
      navigate("/");
    }

    socket.addEventListener("message", (event) => {
      console.log(event.data);
      let data = JSON.parse(event.data);
      if (data.action === "roll") {
        renderRoll(data.data.result);
      }
    });
  }, []);
  const renderRoll = (arr) => {
    let results = arr.slice(1, -1).split(",");
    let i = 0;
    setInterval(() => {
      if (i >= results.length) {
        return;
      } else {
        setDice(results[i]);
      }
      i++;
    }, 250);
  };
  const roll = () => {
    socket.send(
      JSON.stringify({
        room_id: gameInfo.room_id,
        player_id: gameInfo.player_id,
        action: "roll",
        data: {},
      })
    );
  };
  return (
    <div>
      <div>
        room_id: {gameInfo.room_id}
        <br />
      </div>
      <div>
        player_id: {gameInfo.player_id}
        <br />
      </div>
      <div>{dice}</div>
      <button onClick={roll}>roll 1d100</button>
    </div>
  );
};
