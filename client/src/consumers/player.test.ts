import test from "ava";
import { PlayerJoined } from "../messages/events";
import { PlayerManager } from "./player";

test(`PlayerManager - player joined`, t => {
  let playerJoined: string
  const playerMgr = new PlayerManager({joinPlayer: p => playerJoined = p});
  playerMgr.consume(new PlayerJoined("testName"));
  t.deepEqual(playerJoined, "testName");
});
