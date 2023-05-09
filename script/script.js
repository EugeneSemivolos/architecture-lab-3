import {serverUrl, cmd, moveFigure } from "./processing-functions.js";
import { sendHTTPRequest } from "./request-api.js";

const form = document.querySelector('form');
const textarea = document.querySelector('textarea');
const greenFrameBtn = document.querySelector(".gf-script");
const diagonalMoveBtn = document.querySelector(".dm-script");

greenFrameBtn.addEventListener("click", () => {

  sendHTTPRequest(serverUrl, cmd.setGreenFrame())
    .then((response) => console.log(response))
    .catch((error) => console.error(error));
});

form.addEventListener("submit", (e) => {
  e.preventDefault();
  const commandString = textarea.value.trim();

  sendHTTPRequest(serverUrl, commandString)
    .then((response) => console.log(response))
    .catch((error) => console.error(error));
});

diagonalMoveBtn.addEventListener("click", async () => {
  diagonalMoveBtn.disabled = true;
  await moveFigure();
  diagonalMoveBtn.disabled = false;
});
