import { expect, test } from "vitest";
import { MissionDetailsService, MissionTimeline, type MissionDetailsValues } from "./MissionDetails-service";

test("Get Mission", ()=>{
  const service = new MissionDetailsService({mission: 2} as MissionDetailsValues);
  expect(service.mission).to.equal(2);
});

test("Get Mission Timeline", ()=>{
  const service = new MissionDetailsService({missionTimeline: MissionTimeline.Future} as MissionDetailsValues);
  expect(service.missionTimeLine).to.equal(MissionTimeline.Future);
});

test("Get Team votes", ()=>{
  const service = new MissionDetailsService({
    teamVotes: {
      votes: [{
        team: new Set<string>(["p1", "p2"]),
        approved: true, 
        playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])
      }]
    }
  }  as MissionDetailsValues);
  expect(service.teamVotes).to.deep.equal({votes: [{
    team: new Set<string>(["p1", "p2"]),
    approved: true, 
    playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])
  }]});
});

test("Get team from vote as string", ()=>{
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
  } as MissionDetailsValues);
  expect(service.teamFromVoteAsString(0)).to.equal("p1, p2");
  expect(service.teamFromVoteAsString(1)).to.equal("p1, p3");
});

test("Has mission passed", ()=>{
  const service = new MissionDetailsService({missionResult: true} as MissionDetailsValues);
  expect(service.hasMissionSucceeded).to.be.true;
});

test("Number of successes", ()=>{
  const service = new MissionDetailsService({teamSize: 4, nbFailures: 1} as MissionDetailsValues);
  expect(service.nbSuccesses).to.equal(3);
});

test("Should Show Votes", ()=>{
  let service = new MissionDetailsService({missionTimeline: MissionTimeline.Current} as MissionDetailsValues);
  expect(service.shouldShowVotes).to.be.true;
  
  service = new MissionDetailsService({missionTimeline: MissionTimeline.Past} as MissionDetailsValues);
  expect(service.shouldShowVotes).to.be.true;
  
  service = new MissionDetailsService({missionTimeline: MissionTimeline.Future} as MissionDetailsValues);
  expect(service.shouldShowVotes).to.be.false;
});

test("Should Show Mission Results", ()=>{
  let service = new MissionDetailsService({missionTimeline: MissionTimeline.Current} as MissionDetailsValues);
  expect(service.shouldShowMissionResult).to.be.false;
  
  service = new MissionDetailsService({missionTimeline: MissionTimeline.Past} as MissionDetailsValues);
  expect(service.shouldShowMissionResult).to.be.true;
  
  service = new MissionDetailsService({missionTimeline: MissionTimeline.Future} as MissionDetailsValues);
  expect(service.shouldShowMissionResult).to.be.false;
});