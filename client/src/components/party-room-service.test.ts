import {expect, test} from "vitest";
import { StartGame } from "../messages/commands";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { PartyRoomService } from "./party-room-service";

test("Party Code", ()=> {
  const service = new PartyRoomService({
    partyCode: "code",
    players: [],
  }, null);

  expect(service.partyCode).to.equal("code");
});

test("Players", ()=> {
  const service = new PartyRoomService({
    partyCode: "code",
    players: ["a", "b", "c"],
  }, null);

  expect(service.players).to.deep.equal(["a", "b", "c"]);
});

test("Start Game", ()=> {
  const dispatcher = new DispatcherMock();
  const service = new PartyRoomService({
    partyCode: "code",
    players: ["a", "b", "c"],
  }, dispatcher);

  service.startGame();
  expect(dispatcher.receivedMessage).to.be.instanceof(StartGame);
  expect(dispatcher.receivedMessage).to.deep.equal(new StartGame());
});