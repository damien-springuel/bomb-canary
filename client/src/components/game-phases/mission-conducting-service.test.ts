import {expect, test} from "vitest";
import { MissionConductingService } from "./mission-conducting-service";

test("Get Current Team", () => {
  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(),
    peopleThatWorkedOnMission: new Set<string>(),
    player: "",
    playerMissionSuccess: true,
  });
  expect(service.currentTeam).to.deep.equal([]);

  service = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatWorkedOnMission: new Set<string>(),
    player: "",
    playerMissionSuccess: true,
  });
  expect(service.currentTeam).to.deep.equal(["a", "b", "c"]);
});

test("Get Current Team as string", () => {
  let service = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatWorkedOnMission: new Set<string>(),
    player: "",
    playerMissionSuccess: true,
  });
  expect(service.currentTeamAsString).to.deep.equal("a, b and c");
});

test("Has given player worked on mission", ()=> {
  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "",
    playerMissionSuccess: true,
  });
  expect(service.hasGivenPlayerWorkedOnMission("a")).to.be.true;
  expect(service.hasGivenPlayerWorkedOnMission("b")).to.be.false;
});

test("Is player in current mission", ()=> {
  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "a",
    playerMissionSuccess: true,
  });
  expect(service.isPlayerInCurrentMission).to.be.true;

  service = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "c",
    playerMissionSuccess: true,
  });
  expect(service.isPlayerInCurrentMission).to.be.false;
});

test("Has Player worked on mission", ()=> {
  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "a",
    playerMissionSuccess: true,
  });
  expect(service.hasPlayerWorkedOnMission).to.be.true;

  service = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "b",
    playerMissionSuccess: true,
  });
  expect(service.hasPlayerWorkedOnMission).to.be.false;
});

test("Get player mission result", ()=> {
  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "a",
    playerMissionSuccess: true,
  });
  expect(service.playerMissionSuccess).to.be.true;

  service = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "b",
    playerMissionSuccess: false,
  });
  expect(service.playerMissionSuccess).to.be.false;
});
