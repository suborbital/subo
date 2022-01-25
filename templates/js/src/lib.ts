import { Env, LogLevel } from "./env";

export const env = new Env();

export const run = (input: string, ident: number): string => {
  env.logMsg("Hello " + input, LogLevel.Info, ident);

  return "Hello " + input;
};
