import {expect, test} from "vitest";
import { AppService } from "./App-service";
import { DispatcherMock } from "./messages/dispatcher.test-utils";
import { AppLoaded } from "./messages/events";
import { Page } from "./types/types";

test("App mounted", ()=> {
  const dispatcher = new DispatcherMock();
  const service = new AppService(null, dispatcher);
  service.appMounted();
  expect(dispatcher.receivedMessage).to.be.instanceof(AppLoaded);
  expect(dispatcher.receivedMessage).to.deep.equal(new AppLoaded());
});

test("Is page Lobby", ()=> {
  let service = new AppService({
    pageToShow: Page.Lobby,
  }, null);
  expect(service.isPageLobby).to.be.true;

  service = new AppService({
    pageToShow: Page.Game,
  }, null);
  expect(service.isPageLobby).to.be.false;
});

test("Is page Party Room", ()=> {
  let service = new AppService({
    pageToShow: Page.PartyRoom,
  }, null);
  expect(service.isPagePartyRoom).to.be.true;

  service = new AppService({
    pageToShow: Page.Game,
  }, null);
  expect(service.isPagePartyRoom).to.be.false;
});

test("Is page Party Room", ()=> {
  let service = new AppService({
    pageToShow: Page.Game,
  }, null);
  expect(service.isPageGame).to.be.true;

  service = new AppService({
    pageToShow: Page.Lobby,
  }, null);
  expect(service.isPageGame).to.be.false;
});