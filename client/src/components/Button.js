import styled from "styled-components";
import React from "react";
import { accent, text } from "./vars";
const Button = styled.button`
  background: none;
  color: ${text};
  border: 2px solid ${text};
  padding: 8px 24px;
  font-size: 18px;
  border-radius: 8px;
  &:hover {
    border-color: ${accent};
  }
`;

export default (props) => {
  return <Button onClick={props.handler}>{props.text}</Button>;
};
