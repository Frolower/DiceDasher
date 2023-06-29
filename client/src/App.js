import React, { useEffect } from "react";
import { useSelector, useDispatch } from "react-redux";
export default () => {
  const gameInfo = useSelector((state) => state.game);
  useEffect(() => {
    console.log(gameInfo);
  }, []);
  return (
    <div className="App">
      <h1>dicedasher</h1>
      <a href="/createRoom">create room</a>
      <hr />
      <a href="/join">join</a>
    </div>
  );
};
