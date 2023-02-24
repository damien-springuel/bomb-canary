import { expect, test } from "vitest";
import { EventsReplayStarted, PlayerJoined, SpiesRevealed } from "../messages/events";
import { PlayerConsumer, type PlayerStore } from "./player";

test(`PlayerManager - playerConsumer joined`, () => {
  let playerJoined: string
  const playerConsumer = new PlayerConsumer({joinPlayer: p => {playerJoined = p}} as PlayerStore);
  playerConsumer.consume(new PlayerJoined("testName"));
  expect(playerJoined).to.equal("testName");
});

test(`PlayerManager - events replay started`, () => {
  let definedPlayer: string
  const playerConsumer = new PlayerConsumer({definePlayer: p => {definedPlayer = p}} as PlayerStore);
  playerConsumer.consume(new EventsReplayStarted("testName"));
  expect(definedPlayer).to.equal("testName");
});

test(`PlayerManager - spies revealed`, () => {
  let rememberedSpies: Set<string>;
  const playerConsumer = new PlayerConsumer({rememberSpies: s => {rememberedSpies = s}} as PlayerStore);
  playerConsumer.consume(new SpiesRevealed(new Set<string>(["a", "b"])));
  expect(rememberedSpies).to.deep.equal(new Set<string>(["a", "b"]));
});
