import {expect,test} from "vitest";
import { ViewIdentity } from "../messages/commands";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { Dialog, GamePhase } from "../types/types";
import { GameService, type GameValues } from "./Game-service";

test("View Identity", ()=> {
  const dispatcher = new DispatcherMock();
  const service = new GameService({} as GameValues, dispatcher);

  service.viewIdentity();
  expect(dispatcher.receivedMessage).to.be.instanceof(ViewIdentity);
  expect(dispatcher.receivedMessage).to.deep.equal(new ViewIdentity());
});

test("Is game phase team selection", ()=> {
  let service = new GameService(
    {currentGamePhase: GamePhase.TeamSelection} as GameValues, 
    null);

  expect(service.isTeamSelectionPhase).to.be.true;
  
  service = new GameService(
    {currentGamePhase: GamePhase.TeamVote} as GameValues, 
    null);

  expect(service.isTeamSelectionPhase).to.be.false;
});

test("Is game phase team vote", ()=> {
  let service = new GameService(
    {currentGamePhase: GamePhase.TeamVote} as GameValues, 
    null);

  expect(service.isTeamVotePhase).to.be.true;

  service = new GameService(
    {currentGamePhase: GamePhase.TeamSelection} as GameValues, 
    null);

  expect(service.isTeamVotePhase).to.be.false;
});

test("Is game phase Mission", ()=> {
  let service = new GameService(
    {currentGamePhase: GamePhase.Mission} as GameValues, 
    null);

  expect(service.isMissionConductingPhase).to.be.true;
  
  service = new GameService(
    {currentGamePhase: GamePhase.TeamSelection} as GameValues, 
    null);

  expect(service.isMissionConductingPhase).to.be.false;
});

test("Is dialog shown Identity", ()=> {
  let service = new GameService(
    {dialogShown: Dialog.Identity} as GameValues, 
    null);

  expect(service.isDialogShownIdentity).to.be.true;
  
  service = new GameService(
    {dialogShown: null} as GameValues, 
    null);

  expect(service.isDialogShownIdentity).to.be.false;
});

test("Is dialog shown Mission Details", ()=> {
  let service = new GameService(
    {dialogShown: Dialog.MissionDetails} as GameValues, 
    null);

  expect(service.isDialogShownMissionDetails).to.be.true;
  
  service = new GameService(
    {dialogShown: null} as GameValues, 
    null);

  expect(service.isDialogShownMissionDetails).to.be.false;
});

test("Is dialog shown last mission result", ()=> {
  let service = new GameService(
    {dialogShown: Dialog.LastMissionResult} as GameValues, 
    null);

  expect(service.isDialogShownLastMissionResult).to.be.true;
  
  service = new GameService(
    {dialogShown: null} as GameValues, 
    null);

  expect(service.isDialogShownLastMissionResult).to.be.false;
});