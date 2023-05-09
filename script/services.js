export class UrlService {
  static url = "http://127.0.0.1:17000/?cmd=";

  static create = (options) => {
    let newUrl = this.url;
    for (const action in options) {
      newUrl += `${action}`;
      const arrayOfArgs = options[action];
      if (arrayOfArgs) newUrl += " " + arrayOfArgs.join(" ");
      newUrl += ",";
    }
    return newUrl.slice(0, -1);
  };
  static parseCommandString = (commandString) => {
    const commands = commandString.split(",").map((c) => c.trim());
    const result = {};

    commands.forEach((c) => {
      const [command, ...args] = c.split(" ");
      if (command) {
        result[command] = args.length ? args : null;
      }
    });

    return result;
  };
}
