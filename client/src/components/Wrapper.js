import styled from "styled-components";
import { background, text } from "./vars";

export default styled.div`
  box-sizing: border-box;
  width: 100%;
  height: 100%;
  background-color: ${background};
  color: ${text};
  font-family: "Trebuchet MS", "Lucida Sans Unicode", "Lucida Grande",
    "Lucida Sans", Arial, sans-serif;
  display: flex;
  justify-content: center;
  align-items: center;
`;
