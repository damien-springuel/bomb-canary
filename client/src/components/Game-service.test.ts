import {expect,test} from "vitest";
import { ViewIdentity } from "../messages/commands";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { Dialog, GamePhase } from "../types/types";
import { GameService } from "./Game-service";

test("View Identity", ()=> {
  const dispatcher = new DispatcherMock();
  const service = new GameService({currentGamePhase: null, dialogShown: null}, dispatcher);

  service.viewIdentity();
  expect(dispatcher.receivedMessage).to.be.instanceof(ViewIdentity);
  expect(dispatcher.receivedMessage).to.deep.equal(new ViewIdentity());
});

test("Is game phase team selection", ()=> {
  let service = new GameService({
    currentGamePhase: GamePhase.TeamSelection,
    dialogShown: null,
  }, null);

  expect(service.isTeamSelectionPhase).to.be.true;
  
  service = new GameService({
    currentGamePhase: GamePhase.TeamVote,
    dialogShown: null,
  }, null);

  expect(service.isTeamSelectionPhase).to.be.false;
});

test("Is game phase team vote", ()=> {
  let service = new GameService({
    currentGamePhase: GamePhase.TeamVote,
    dialogShown: null,
  }, null);

  expect(service.isTeamVotePhase).to.be.true;

  service = new GameService({
    currentGamePhase: GamePhase.TeamSelection,
    dialogShown: Dialog.Identity,
  }, null);

  expect(service.isTeamVotePhase).to.be.false;
});

test("Is game phase Mission", ()=> {
  let service = new GameService({
    currentGamePhase: GamePhase.Mission,
    dialogShown: null,
  }, null);

  expect(service.isMissionConductingPhase).to.be.true;
  
  service = new GameService({
    currentGamePhase: GamePhase.TeamSelection,
    dialogShown: Dialog.Identity,
  }, null);

  expect(service.isMissionConductingPhase).to.be.false;
});

test("Is dialog shown Identity", ()=> {
  let service = new GameService({
    currentGamePhase: GamePhase.Mission,
    dialogShown: Dialog.Identity,
  }, null);

  expect(service.isDialogShownIdentity).to.be.true;
  
  service = new GameService({
    currentGamePhase: GamePhase.TeamSelection,
    dialogShown: null,
  }, null);

  expect(service.isDialogShownIdentity).to.be.false;
});