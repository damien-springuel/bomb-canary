import type { Message } from "./messagebus";

export class AppLoaded implements Message{}

export class PartyCreated implements Message{
  constructor(readonly partyCode: string){}
}

export class ServerConnectionClosed implements Message {}
export class ServerConnectionErrorOccured implements Message {}

export class PlayerConnected implements Message {
  constructor(readonly name:string){}
}