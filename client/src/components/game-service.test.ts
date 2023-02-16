import { expect, test } from "vitest";
import { GameService } from "./game-service";

test("Is Current Mission", () => {
  const gameService: GameService = new GameService({missionResults: [], missionRequirements:[]});
  expect(gameService.isCurrentMission(0)).to.be.true;
  expect(gameService.isCurrentMission(1)).to.be.false;
});

test("Should mission tag have no border", ()=> {
  const gameService: GameService = new GameService({
    missionResults: [
      {nbFails: 1, success: false}, 
    ], missionRequirements:[]});
  expect(gameService.shouldMissionTagHaveNoBorder(0)).to.be.true;
  expect(gameService.shouldMissionTagHaveNoBorder(1)).to.be.true;
  expect(gameService.shouldMissionTagHaveNoBorder(2)).to.be.false;
});

test("Should mission tag text be gray", ()=> {
  const gameService: GameService = new GameService({
    missionResults: [
      {nbFails: 1, success: false}, 
    ], missionRequirements:[]});
  expect(gameService.shouldMissionTagTextBeGray(0)).to.be.true;
  expect(gameService.shouldMissionTagTextBeGray(1)).to.be.true;
  expect(gameService.shouldMissionTagTextBeGray(2)).to.be.false;
});

test("Should mission tag show success", ()=> {
  const gameService: GameService = new GameService({
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
  const gameService: GameService = new GameService({
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
  const gameService: GameService = new GameService({
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
  const gameService: GameService = new GameService({
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
  const gameService: GameService = new GameService({
    missionResults: [], 
    missionRequirements:[]});
  expect(gameService.missions).to.deep.equal([0,1,2,3,4]);
});
