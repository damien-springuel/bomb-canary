import {expect, test} from "vitest";
import { AppService } from "./App-service";
import { DispatcherMock } from "./messages/dispatcher.test-utils";
import { AppLoaded } from "./messages/events";
import { Page } from "./types/types";

test("App mounted", ()=> {
  const dispatcher = new DispatcherMock();
  const service = new AppService(dispatcher);
  service.appMounted();
  expect(dispatcher.receivedMessage).to.be.instanceof(AppLoaded);
  expect(dispatcher.receivedMessage).to.deep.equal(new AppLoaded());
});
