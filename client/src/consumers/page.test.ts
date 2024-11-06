import { expect, test } from "vitest";
import { CloseDialog, ViewIdentity, ViewMissionDetails } from "../messages/commands";
import { ServerConnectionClosed, SpiesRevealed } from "../messages/events";
import { PageConsumer, type RoomStore } from "./page";

test(`Page Manager - show party room on server connection closed`, () => {
  let lobbyShown = false
  const pageConsumer = new PageConsumer({showPartyRoom: ()=> {lobbyShown = true}} as RoomStore);
  pageConsumer.consume(new ServerConnectionClosed());
  expect(lobbyShown).to.be.true;
});

test(`Page Manager - show game room on spies revealed`, () => {
  let gameShown = false;
  let identityShown = false;
  const pageConsumer = new PageConsumer({
    showGameRoom: () => {gameShown = true},
    showIdentity: () => {identityShown = true}
  } as RoomStore);
  pageConsumer.consume(new SpiesRevealed(null));
  expect(gameShown).to.be.true;
  expect(identityShown).to.be.true;
});

test(`Page Manager - show identity on view identity`, () => {
  let identityShown = false;
  const pageConsumer = new PageConsumer({showIdentity: () => {identityShown = true}} as RoomStore);
  pageConsumer.consume(new ViewIdentity());
  expect(identityShown).to.be.true;
});

test(`Page Manager - show proper mission details on view mission details`, () => {
  let missionDetailsShown = false;
  let missionGiven: number = null;
  const pageConsumer = new PageConsumer({
    showMissionDetails: (mission: number) => {
      missionDetailsShown = true;
      missionGiven = mission
    }
  } as RoomStore);
  pageConsumer.consume(new ViewMissionDetails(3));
  expect(missionDetailsShown).to.be.true;
  expect(missionGiven).to.equal(3);
});

test(`Page Manager - close dialog on close dialog`, () => {
  let dialogClosed = false;
  const pageConsumer = new PageConsumer({closeDialog: () => {dialogClosed = true}} as RoomStore);
  pageConsumer.consume(new CloseDialog());
  expect(dialogClosed).to.be.true;
});
