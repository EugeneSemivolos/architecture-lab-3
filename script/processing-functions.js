import { sendHTTPRequest } from "./request-api.js";

const serverUrl = "http://127.0.0.1:17000";


const startPos = 0.1;
const fps = 0.05;

const delayTime = 200;

const cmd = {
  setGreenFrame: () => "green\nbgrect 0.25 0.25 0.75 0.75\nupdate",
  reset: () => "reset",
  createFigure: (x, y) => `figure ${x} ${y}\nupdate`,
  move: (x, y) => `move ${x} ${y}\nupdate`,
}

const moveFigure = async () => {
  const command = cmd.createFigure(startPos, startPos);
  await sendHTTPRequest(serverUrl, command);
  
  for (let i = 1; i < 1/fps; i++) {
    const command = cmd.move(fps * i, fps * i);
    await sendHTTPRequest(serverUrl, command);

    await delay();
  }

  for (let i = 1/fps - 1; i > 0; i--) {
    const command = cmd.move(fps * i, fps * i);
    await sendHTTPRequest(serverUrl, command);

    await delay();
  }
};
  
const delay = async() => new Promise((resolve) => setTimeout(resolve, delayTime));

export { serverUrl, cmd, moveFigure };
