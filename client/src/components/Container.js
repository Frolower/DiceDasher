import styled from "styled-components";
import { accent, container } from "./vars";
export default styled.div`
  width: 80%;
  height: 80%;
  margin: auto auto;
  border-radius: 8px;
  padding: 24px;
  background-color: ${container};

  .mainPage {
    height: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    .logo {
      margin-top: 48px;
      text-align: center;
      font-size: 64px;
    }
    .controls {
      height: 100%;
      display: flex;
      flex-direction: row;
      justify-content: space-between;
      align-items: center;
      div {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      .new {
        * {
          width: 70%;
        }
      }
      .join {
        flex-direction: column;
        * {
          margin: 24px 0;
          width: 70%;
          box-sizing: border-box;
        }
      }
    }
  }
`;
