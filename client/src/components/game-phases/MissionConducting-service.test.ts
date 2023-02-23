import {expect, test} from "vitest";
import { FailMission, SucceedMission } from "../../messages/commands";
import type { Dispatcher } from "../../messages/dispatcher";
import type { Message } from "../../messages/messagebus";
import { MissionConductingService } from "./MissionConducting-service";

test("Get Current Team", () => {
  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(),
    peopleThatWorkedOnMission: new Set<string>(),
    player: "",
    playerMissionSuccess: true,
  }, null);
  expect(service.currentTeam).to.deep.equal([]);

  service = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatWorkedOnMission: new Set<string>(),
    player: "",
    playerMissionSuccess: true,
  }, null);
  expect(service.currentTeam).to.deep.equal(["a", "b", "c"]);
});

test("Get Current Team as string", () => {
  let service = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatWorkedOnMission: new Set<string>(),
    player: "",
    playerMissionSuccess: true,
  }, null);
  expect(service.currentTeamAsString).to.deep.equal("a, b and c");
});

test("Has given player worked on mission", ()=> {
  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "",
    playerMissionSuccess: true,
  }, null);
  expect(service.hasGivenPlayerWorkedOnMission("a")).to.be.true;
  expect(service.hasGivenPlayerWorkedOnMission("b")).to.be.false;
});

test("Is player in current mission", ()=> {
  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "a",
    playerMissionSuccess: true,
  }, null);
  expect(service.isPlayerInCurrentMission).to.be.true;

  service = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "c",
    playerMissionSuccess: true,
  }, null);
  expect(service.isPlayerInCurrentMission).to.be.false;
});

test("Has Player worked on mission", ()=> {
  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "a",
    playerMissionSuccess: true,
  }, null);
  expect(service.hasPlayerWorkedOnMission).to.be.true;

  service = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "b",
    playerMissionSuccess: true,
  }, null);
  expect(service.hasPlayerWorkedOnMission).to.be.false;
});

test("Get player mission result", ()=> {
  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "a",
    playerMissionSuccess: true,
  }, null);
  expect(service.playerMissionSuccess).to.be.true;

  service = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "b",
    playerMissionSuccess: false,
  }, null);
  expect(service.playerMissionSuccess).to.be.false;
});

test("Succeed Mission", ()=> {
  let givenMessage: Message;
  let dispatcher: Dispatcher = {dispatch(message){givenMessage = message}};

  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "a",
    playerMissionSuccess: true,
  }, dispatcher);

  service.succeedMission();
  expect(givenMessage).to.be.instanceof(SucceedMission);
  expect(givenMessage).to.deep.equal(new SucceedMission());

});

test("Fail Mission", ()=> {
  let givenMessage: Message;
  let dispatcher: Dispatcher = {dispatch(message){givenMessage = message}};

  let service: MissionConductingService = new MissionConductingService({
    currentTeam: new Set<string>(["a", "b"]),
    peopleThatWorkedOnMission: new Set<string>(["a"]),
    player: "a",
    playerMissionSuccess: true,
  }, dispatcher);

  service.failMission();
  expect(givenMessage).to.be.instanceof(FailMission);
  expect(givenMessage).to.deep.equal(new FailMission());
});
