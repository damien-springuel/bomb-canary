import {expect, test} from "vitest";
import { TeamSelectionService } from "./team-selection-service";

test("Is given player in current team", ()=> {
  const service = new TeamSelectionService({
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b"],
    });

  expect(service.isGivenPlayerInTeam("a")).to.be.true;
  expect(service.isGivenPlayerInTeam("c")).to.be.false;
});

test("Is the player the leader", ()=> {
  let service = new TeamSelectionService({
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b"],
    });
  expect(service.isPlayerTheLeader).to.be.true;

  service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "b",
      currentTeamVoteNb: 1,
      players: ["a", "b"],
    });

  expect(service.isPlayerTheLeader).to.be.false;
});