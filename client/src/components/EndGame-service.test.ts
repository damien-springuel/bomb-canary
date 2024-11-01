import {expect, test} from "vitest";
import { EndGameService, type EndGameValues } from "./EndGame-service";
import { Allegiance } from "../types/types";

test("Spies as String", () => {
  const service = new EndGameService({
    spies: new Set<string>(["a", "b", "c"]),
  } as EndGameValues);
  expect(service.spiesAsString).to.equal("a, b, c");
});

test("Spies have won", () => {
  let service = new EndGameService({winner: Allegiance.Spies} as EndGameValues);
  expect(service.spiesHaveWon).to.be.true;
  
  service = new EndGameService({winner: Allegiance.Resistance} as EndGameValues);
  expect(service.spiesHaveWon).to.be.false;
});

test("Player has won", () => {
  let service = new EndGameService({
    winner: Allegiance.Spies,
    spies: new Set<string>(["a", "b", "c"]),
    player: "a"
  } as EndGameValues);
  expect(service.playerHasWon).to.be.true;
  
  service = new EndGameService({
    winner: Allegiance.Resistance,
    spies: new Set<string>(["a", "b", "c"]),
    player: "a"
  } as EndGameValues);
  expect(service.playerHasWon).to.be.false;
  
  service = new EndGameService({
    winner: Allegiance.Spies,
    spies: new Set<string>(["a", "b", "c"]),
    player: "d"
  } as EndGameValues);
  expect(service.playerHasWon).to.be.false;
  
  service = new EndGameService({
    winner: Allegiance.Resistance,
    spies: new Set<string>(["a", "b", "c"]),
    player: "d"
  } as EndGameValues);
  expect(service.playerHasWon).to.be.true;
});