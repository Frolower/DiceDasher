import React, { useEffect, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { Link, useNavigate } from "react-router-dom";
import Wrapper from "./components/Wrapper";
import Container from "./components/Container";
import Button from "./components/Button";
import Input from "./components/Input";
import "./App.css";
import { setPlayerId, setRoomId } from "./store/gameSlice";

export default () => {
  const [roomid, set_room_id] = useState("");
  const navigate = useNavigate();
  const gameInfo = useSelector((state) => state.game);
  const dispatch = useDispatch();
  useEffect(() => {
    console.log(gameInfo);
  }, []);

  const createHandler = () => {
    fetch("http://localhost:8080/newroom", {
      method: "GET",
    })
      .then((response) => {
        return response.json();
      })
      .then((json) => {
        console.log(json);

        dispatch(setRoomId({ room_id: json.room_id }));

        dispatch(setPlayerId({ player_id: json.player_id }));

        navigate("/room");
      });
  };

  const handleInputChange = (e) => {
    set_room_id(e.target.value);
  };

  const handleJoinButton = () => {
    console.log(roomid);

    dispatch(setRoomId({ room_id: roomid }));
    dispatch(setPlayerId({ player_id: "" }));

    navigate("/room");
  };
  return (
    <Wrapper>
      <Container>
        <div className="mainPage">
          <div className="logo">DiceDasher</div>
          <div className="controls">
            <div className="new">
              <Button text="Create room" handler={createHandler} />
            </div>
            <div className="join">
              <Input
                placeholder="room_id"
                value={roomid}
                onChange={handleInputChange}
              />
              <Button text="Join" handler={handleJoinButton} />
            </div>
          </div>
        </div>
      </Container>
    </Wrapper>
  );
};
