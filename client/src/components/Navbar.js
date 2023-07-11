import React from "react";
import styled from "styled-components";
import { text } from "./vars";

export default (props) => {
  let color = "red";
  if (props.status === 0) {
    color = "red";
  } else if (props.status == 1) {
    color = "green";
  } else {
    color = "yellow";
  }
  const Navbar = styled.div`
    width: 100%;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    padding-top: 12px;
    padding-bottom: 18px;
    border-bottom: 2px solid ${text};
    .info {
      display: flex;
      flex-direction: row;
      align-items: center;
      span {
        margin: 0 12px;
      }
    }
    .logo {
      font-size: 24px;
    }
    .status {
      margin-left: 12px;
      width: 16px;
      height: 16px;
      border-radius: 50%;
      background-color: ${color};
    }
  `;

  return (
    <Navbar>
      <div className="logo">DiceDasher</div>
      <div className="info">
        <span>room_id: {props.room_id}</span>
        <span>player_id: {props.player_id} </span>

        <div className="status"></div>
      </div>
    </Navbar>
  );
};
