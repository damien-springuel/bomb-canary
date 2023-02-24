import { expect, test } from "vitest";
import { CloseDialog, ViewIdentity } from "../messages/commands";
import { PartyCreated, ServerConnectionClosed, SpiesRevealed } from "../messages/events";
import { PageConsumer, type RoomStore } from "./page";

test(`Page Manager - show lobby on server connection closed`, t => {
  let lobbyShown = false
  const pageConsumer = new PageConsumer({showLobby: ()=> {lobbyShown = true}} as RoomStore);
  pageConsumer.consume(new ServerConnectionClosed());
  expect(lobbyShown).to.be.true;
});

test(`Page Manager - show party room on party created`, t => {
  let receivedPartyCode: string;
  const pageConsumer = new PageConsumer({showPartyRoom: code => {receivedPartyCode = code}} as RoomStore);
  pageConsumer.consume(new PartyCreated("testCode"));
  expect(receivedPartyCode).to.equal("testCode");
});

test(`Page Manager - show game room on spies revealed`, t => {
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

test(`Page Manager - show identity on view identity`, t => {
  let identityShown = false;
  const pageConsumer = new PageConsumer({showIdentity: () => {identityShown = true}} as RoomStore);
  pageConsumer.consume(new ViewIdentity());
  expect(identityShown).to.be.true;
});

test(`Page Manager - close dialog on close dialog`, t => {
  let dialogClosed = false;
  const pageConsumer = new PageConsumer({closeDialog: () => {dialogClosed = true}} as RoomStore);
  pageConsumer.consume(new CloseDialog());
  expect(dialogClosed).to.be.true;
});
