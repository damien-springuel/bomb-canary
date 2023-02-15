import { expect, test } from "vitest";
import { EventsReplayStarted, PlayerJoined, SpiesRevealed } from "../messages/events";
import { PlayerManager, type PlayerStore } from "./player";

test(`PlayerManager - player joined`, () => {
  let playerJoined: string
  const playerMgr = new PlayerManager({joinPlayer: p => {playerJoined = p}} as PlayerStore);
  playerMgr.consume(new PlayerJoined("testName"));
  expect(playerJoined).to.equal("testName");
});

test(`PlayerManager - events replay started`, () => {
  let definedPlayer: string
  const playerMgr = new PlayerManager({definePlayer: p => {definedPlayer = p}} as PlayerStore);
  playerMgr.consume(new EventsReplayStarted("testName"));
  expect(definedPlayer).to.equal("testName");
});

test(`PlayerManager - spies revealed`, () => {
  let rememberedSpies: Set<string>;
  const playerMgr = new PlayerManager({rememberSpies: s => {rememberedSpies = s}} as PlayerStore);
  playerMgr.consume(new SpiesRevealed(new Set<string>(["a", "b"])));
  expect(rememberedSpies).to.deep.equal(new Set<string>(["a", "b"]));
});
