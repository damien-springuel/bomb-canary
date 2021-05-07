import test from "ava";
import { CloseDialog, ViewIdentity } from "../messages/commands";
import { PartyCreated, ServerConnectionClosed, SpiesRevealed } from "../messages/events";
import { PageManager, RoomStore } from "./page";

test(`Page Manager - show lobby on server connection closed`, t => {
  let lobbyShown = false
  const pageMgr = new PageManager({showLobby: ()=> {lobbyShown = true}} as RoomStore);
  pageMgr.consume(new ServerConnectionClosed());
  t.true(lobbyShown);
});

test(`Page Manager - show party room on party created`, t => {
  let receivedPartyCode: string;
  const pageMgr = new PageManager({showPartyRoom: code => {receivedPartyCode = code}} as RoomStore);
  pageMgr.consume(new PartyCreated("testCode"));
  t.deepEqual(receivedPartyCode, "testCode");
});

test(`Page Manager - show game room on spies revealed`, t => {
  let gameShown = false;
  let identityShown = false;
  const pageMgr = new PageManager({
    showGameRoom: () => {gameShown = true},
    showIdentity: () => {identityShown = true}
  } as RoomStore);
  pageMgr.consume(new SpiesRevealed(null));
  t.true(gameShown);
  t.true(identityShown);
});

test(`Page Manager - show identity on view identity`, t => {
  let identityShown = false;
  const pageMgr = new PageManager({showIdentity: () => {identityShown = true}} as RoomStore);
  pageMgr.consume(new ViewIdentity());
  t.true(identityShown);
});

test(`Page Manager - close dialog on close dialog`, t => {
  let dialogClosed = false;
  const pageMgr = new PageManager({closeDialog: () => {dialogClosed = true}} as RoomStore);
  pageMgr.consume(new CloseDialog());
  t.true(dialogClosed);
});
