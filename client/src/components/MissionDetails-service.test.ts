import { expect, test } from "vitest";
import { MissionDetailsService } from "./MissionDetails-service";

test("Get Mission", ()=>{
  const service = new MissionDetailsService({mission: 2, teamVotes: null});
  expect(service.mission).to.equal(2);
});

test("Get Mission", ()=>{
  const service = new MissionDetailsService({
    mission: 2, 
    teamVotes: {votes: [{approved: true, playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])}]}
  });
  expect(service.teamVotes).to.deep.equal({votes: [{approved: true, playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])}]});
});