import {expect, test} from "vitest";
import { JoinParty, StartGame } from "../messages/commands";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { PartyRoomService, type PartyRoomValues } from "./PartyRoom-service";

test("Players", ()=> {
  const service = new PartyRoomService({
    players: ["a", "b", "c"],
  } as PartyRoomValues, null);

  expect(service.players).to.deep.equal(["a", "b", "c"]);
});

test("Has Player joined", ()=> {
  const service = new PartyRoomService({
    hasPlayerJoined: true
  } as PartyRoomValues, null);

  expect(service.hasPlayerJoined).to.be.true;
});

test("Start Game", ()=> {
  const dispatcher = new DispatcherMock();
  const service = new PartyRoomService({} as PartyRoomValues, dispatcher);

  service.startGame();
  expect(dispatcher.receivedMessage).to.deep.equal(new StartGame());
});

test("Join Game", ()=> {
  const dispatcher = new DispatcherMock();
  const service = new PartyRoomService({} as PartyRoomValues, dispatcher);

  service.joinParty("name");
  expect(dispatcher.receivedMessage).to.deep.equal(new JoinParty("name"));
});

test("Can start game", ()=> {
  const dispatcher = new DispatcherMock();
  let service = new PartyRoomService({
    players: ["1", "2", "3", "4"]
  } as PartyRoomValues, dispatcher);

  expect(service.canStartGame).to.be.false;
  
  service = new PartyRoomService({
    players: ["1", "2", "3", "4", "5"]
  } as PartyRoomValues, dispatcher);

  expect(service.canStartGame).to.be.true;
});