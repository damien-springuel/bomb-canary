import test from "ava";
import { PartyCreated, ServerConnectionClosed } from "../messages/events";
import { PageManager, RoomStore } from "./page";

test(`Page Manager - show lobby on server connection closed`, t => {
  let lobbyShowed = false
  const pageMgr = new PageManager({showLobby: ()=> {lobbyShowed = true}} as RoomStore);
  pageMgr.consume(new ServerConnectionClosed());
  t.true(lobbyShowed);
});

test(`Page Manager - show party room on party created`, t => {
  let receivedPartyCode: string;
  const pageMgr = new PageManager({showPartyRoom: code => {receivedPartyCode = code}} as RoomStore);
  pageMgr.consume(new PartyCreated("testCode"));
  t.deepEqual(receivedPartyCode, "testCode");
});
