import { expect, test } from "vitest";
import { MissionTrackerService } from "./mission-tracker-service";

test("Is Current Mission", () => {
  const gameService: MissionTrackerService = new MissionTrackerService({missionResults: [], missionRequirements:[]});
  expect(gameService.isCurrentMission(0)).to.be.true;
  expect(gameService.isCurrentMission(1)).to.be.false;
});

test("Should mission tag show success", ()=> {
  const gameService: MissionTrackerService = new MissionTrackerService({
    missionResults: [
      {nbFails: 1, success: false}, 
      {nbFails: 0, success: true}, 
    ], missionRequirements:[]});
  expect(gameService.shouldMissionTagShowSuccess(0)).to.be.false;
  expect(gameService.shouldMissionTagShowSuccess(1)).to.be.true;
  expect(gameService.shouldMissionTagShowSuccess(2)).to.be.false;
  expect(gameService.shouldMissionTagShowSuccess(3)).to.be.false;
});

test("Should mission tag show failure", ()=> {
  const gameService: MissionTrackerService = new MissionTrackerService({
    missionResults: [
      {nbFails: 1, success: false}, 
      {nbFails: 0, success: true}, 
    ], missionRequirements:[]});
  expect(gameService.shouldMissionTagShowFailure(0)).to.be.true;
  expect(gameService.shouldMissionTagShowFailure(1)).to.be.false;
  expect(gameService.shouldMissionTagShowFailure(2)).to.be.false;
  expect(gameService.shouldMissionTagShowFailure(3)).to.be.false;
});

test("Should mission tag show nb of people on mission", ()=> {
  const gameService: MissionTrackerService = new MissionTrackerService({
    missionResults: [
      {nbFails: 1, success: false}, 
      {nbFails: 0, success: true}, 
    ], missionRequirements:[]});
  expect(gameService.shouldMissionTagShowNbOfPeopleOnMission(0)).to.be.false;
  expect(gameService.shouldMissionTagShowNbOfPeopleOnMission(1)).to.be.false;
  expect(gameService.shouldMissionTagShowNbOfPeopleOnMission(2)).to.be.true;
  expect(gameService.shouldMissionTagShowNbOfPeopleOnMission(3)).to.be.true;
});

test("Get number of people on mission", ()=> {
  const gameService: MissionTrackerService = new MissionTrackerService({
    missionResults: [
      {nbFails: 1, success: false}, 
      {nbFails: 0, success: true}], 
    missionRequirements:[
      {nbPeopleOnMission: 2, nbFailuresRequiredToFail: 1},
      {nbPeopleOnMission: 3, nbFailuresRequiredToFail: 2},
      {nbPeopleOnMission: 4, nbFailuresRequiredToFail: 3},
      {nbPeopleOnMission: 5, nbFailuresRequiredToFail: 4},
    ]});
  expect(gameService.getNumberPeopleOnMission(0)).to.equal(2);
  expect(gameService.getNumberPeopleOnMission(1)).to.equal(3);
  expect(gameService.getNumberPeopleOnMission(2)).to.equal(4);
  expect(gameService.getNumberPeopleOnMission(3)).to.equal(5);
});

test("Get missions", ()=> {
  const gameService: MissionTrackerService = new MissionTrackerService({
    missionResults: [], 
    missionRequirements:[]});
  expect(gameService.missions).to.deep.equal([0,1,2,3,4]);
});

test("Mission has more than one fail required", ()=> {
  const gameService: MissionTrackerService = new MissionTrackerService({
    missionResults: [
      {nbFails: 1, success: false}, 
      {nbFails: 0, success: true}], 
    missionRequirements:[
      {nbPeopleOnMission: 2, nbFailuresRequiredToFail: 1},
      {nbPeopleOnMission: 3, nbFailuresRequiredToFail: 2},
      {nbPeopleOnMission: 4, nbFailuresRequiredToFail: 1},
      {nbPeopleOnMission: 5, nbFailuresRequiredToFail: 4},
    ]});
  expect(gameService.doesMissionNeedMoreThanOneFail(0)).to.be.false;
  expect(gameService.doesMissionNeedMoreThanOneFail(1)).to.be.true;
  expect(gameService.doesMissionNeedMoreThanOneFail(2)).to.be.false;
  expect(gameService.doesMissionNeedMoreThanOneFail(3)).to.be.true;
});

test("Get nb people required on current mission", ()=> {
  const gameService: MissionTrackerService = new MissionTrackerService({
    missionResults: [
      {nbFails: 1, success: false}, 
      {nbFails: 0, success: true}], 
    missionRequirements:[
      {nbPeopleOnMission: 2, nbFailuresRequiredToFail: 1},
      {nbPeopleOnMission: 3, nbFailuresRequiredToFail: 2},
      {nbPeopleOnMission: 4, nbFailuresRequiredToFail: 1},
      {nbPeopleOnMission: 5, nbFailuresRequiredToFail: 4},
    ]});
  expect(gameService.nbPeopleRequiredOnCurrentMission).to.equal(4);
});
