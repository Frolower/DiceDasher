import styled from "styled-components";
import React from "react";
import { accent, text } from "./vars";
const Input = styled.input`
  background: none;
  color: ${text};
  border: 2px solid ${text};
  padding: 8px 24px;
  font-size: 18px;
  border-radius: 8px;
  outline: none;
  &:focus {
    border: 2px solid ${accent};
  }
`;

export default (props) => {
  return (
    <Input
      placeholder={props.placeholder}
      onChange={props.onChange}
      value={props.value}
    />
  );
};
