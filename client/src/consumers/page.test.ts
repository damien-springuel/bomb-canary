import test from "ava";
import { ServerConnectionClosed } from "../messages/events";
import { PageManager } from "./page";

test(`Page Manager - show lobby on server connection closed`, t => {
  let lobbyShowed = false
  const pageMgr = new PageManager({showLobby: ()=> {lobbyShowed = true}});
  pageMgr.consume(new ServerConnectionClosed());
  t.true(lobbyShowed);
});
