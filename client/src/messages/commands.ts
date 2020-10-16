import type { Message } from "./messagebus";

export interface Command extends Message {}

export class CreateParty implements Command {
  constructor(readonly name:string){}
}