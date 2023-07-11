import React, { useEffect, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Link, useNavigate } from "react-router-dom";
import Wrapper from "./components/Wrapper";
import Container from "./components/Container";
import Button from "./components/Button";
import Input from "./components/Input";
import Navbar from "./components/Navbar";
import useWebSocket from "react-use-websocket";
import styled from "styled-components";
import { text, accent, container } from "./components/vars";
const WEBSOCKET_URL = "ws://localhost:8080/game";

const Game = styled.div`
  padding-top: 32px;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  height: inherit;
  .table {
    display: flex;
    flex-direction: column;
    height: 100%;
    width: 100%;
    .dices {
      display: flex;
      flex-direction: row;
    }
    .dice {
      width: 64px;
      height: 64px;
      background-color: ${accent};
      border-radius: 8px;
      display: flex;
      justify-content: center;
      align-items: center;
      color: ${container};
      font-size: 28px;
      font-weight: bold;
    }
  }
  .info {
    height: 100%;
    width: 25%;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    padding-left: 18px;
    border-left: 2px solid ${text};
    font-size: 18px;
    .players {
      width: 100%;
      height: 50%;
      border-bottom: 2px solid ${text};
      margin-bottom: 12px;
      div {
        margin-bottom: 8px;
      }
    }
    .logs {
      width: 100%;
      height: 50%;
      overflow: scroll;
    }
  }
`;

export default () => {
  const gameInfo = useSelector((state) => state.game);

  const navigate = useNavigate();
  const [con, setCon] = useState(0);
  const [players, setPlayers] = useState([]);
  const [logs, setLogs] = useState([]);
  const [dice, setDice] = useState(100);
  useEffect(() => {
    if (gameInfo.room_id === -1) {
      navigate("/");
    }
  }, []);

  const { sendMessage } = useWebSocket(
    WEBSOCKET_URL +
      "?room_id=" +
      gameInfo.room_id +
      "&player_id=" +
      gameInfo.player_id,
    {
      onOpen: () => {
        console.log("[WS] connected");
        setCon(1);
      },
      onClose: () => {
        console.log("[WS] disconnected");
        setCon(0);
      },
      onMessage: (msg) => {
        let data = JSON.parse(msg.data);
        if (Object.keys(data).includes("action")) {
          let action = data.action;
          switch (action) {
            case "connect":
              playerConnect(data.data);
              setLogs((logs) => [
                ...logs,
                data.player_id + " joined the room.",
              ]);
              break;
            case "disconnect":
              playerConnect(data.data);
              setLogs((logs) => [...logs, data.player_id + " left the room."]);
              break;
            case "roll":
              console.log(data);
              let dices = data.data.slice(1, -1).split(",");
              handleRoll(dices),
                setLogs((logs) => [
                  ...logs,
                  data.player_id + " rolled " + dices[dices.length - 1],
                ]);
              break;
            default:
              console.log("incorrect action");
          }
        } else {
          console.log("bad response");
        }
      },
      onError: (msg) => {
        console.log(msg);
        setCon(2);
      },
    }
  );

  const handleRoll = (dices) => {
    let i = 0;
    setInterval(() => {
      if (i >= dices.length) {
        return dice;
      }
      setDice(dices[i]);
      i++;
    }, 200);
  };
  const playerConnect = (p) => {
    p = JSON.parse(p).players;
    p = p.slice(1, p.length - 1).split(", ");
    console.log(p);
    setPlayers(p);
  };

  const roll = () => {
    sendMessage(
      JSON.stringify({
        action: "roll",
        data: {
          dice: "1d100",
        },
      })
    );
  };
  return (
    <Wrapper>
      <Container>
        <Navbar
          room_id={gameInfo.room_id}
          player_id={gameInfo.player_id}
          status={con}
        />
        <Game>
          <div className="table">
            <div className="dices">
              <div className="dice" onClick={roll}>
                {dice}
              </div>
            </div>
          </div>
          <div className="info">
            <div className="players">
              <div>Players: [{players.length}]</div>

              {players.map((players) => (
                <div>{players}</div>
              ))}
            </div>
            <div className="logs">
              Actions:
              {logs.map((logs) => (
                <div>{logs}</div>
              ))}
            </div>
          </div>
        </Game>
      </Container>
    </Wrapper>
  );
};
