import test from "ava";
import { EventsReplayStarted, PlayerJoined } from "../messages/events";
import { PlayerManager, PlayerStore } from "./player";

test(`PlayerManager - player joined`, t => {
  let playerJoined: string
  const playerMgr = new PlayerManager({joinPlayer: p => {playerJoined = p}} as PlayerStore);
  playerMgr.consume(new PlayerJoined("testName"));
  t.deepEqual(playerJoined, "testName");
});

test(`PlayerManager - events replay started`, t => {
  let definedPlayer: string
  const playerMgr = new PlayerManager({definePlayer: p => {definedPlayer = p}} as PlayerStore);
  playerMgr.consume(new EventsReplayStarted("testName"));
  t.deepEqual(definedPlayer, "testName");
});
