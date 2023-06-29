import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import CreateRoom from "./CreateRoom";
import JoinRoom from "./JoinRoom";
import Room from "./Room";
import { Provider } from "react-redux";
import store from "./store/index";
const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
  // <React.StrictMode>
  <Provider store={store}>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/createRoom" element={<CreateRoom />} />
        <Route path="/join" element={<JoinRoom />} />
        <Route path="/room" element={<Room />} />
      </Routes>
    </BrowserRouter>
  </Provider>
  // </React.StrictMode>
);
