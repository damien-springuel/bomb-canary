import { expect, test } from "vitest";
import { MissionDetailsService } from "./MissionDetails-service";

test("Get Mission", ()=>{
  const service = new MissionDetailsService({mission: 2, teamVotes: null});
  expect(service.mission).to.equal(2);
});

test("Get Team votes", ()=>{
  const service = new MissionDetailsService({
    mission: 2, 
    teamVotes: {
      votes: [{
        team: new Set<string>(["p1", "p2"]),
        approved: true, 
        playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])
      }]
    }
  });
  expect(service.teamVotes).to.deep.equal({votes: [{
    team: new Set<string>(["p1", "p2"]),
    approved: true, 
    playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])
  }]});
});

test("Get mission team by vote", ()=>{
  const service = new MissionDetailsService({
    mission: 2, 
    teamVotes: {
      votes: [{
        team: new Set<string>(["p1", "p2"]),
        approved: true, 
        playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])
      },
      {
        team: new Set<string>(["p1", "p3"]),
        approved: true, 
        playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])
      }]
    }
  });
  expect(service.getTeamFromVote(0)).to.equal("p1, p2");
  expect(service.getTeamFromVote(1)).to.equal("p1, p3");
});