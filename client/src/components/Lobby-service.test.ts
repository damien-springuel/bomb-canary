import {expect, test} from "vitest";
import { CreateParty, JoinParty } from "../messages/commands";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { LobbyService } from "./Lobby-service";

test("Join Party", ()=> {
  const dispatcher = new DispatcherMock();
  const service = new LobbyService(dispatcher);

  service.joinParty("name", "code");
  expect(dispatcher.receivedMessage).to.be.instanceof(JoinParty);
  expect(dispatcher.receivedMessage).to.deep.equal(new JoinParty("name", "code"));
});

test("Create Party", ()=> {
  const dispatcher = new DispatcherMock();
  const service = new LobbyService(dispatcher);

  service.createParty("name");
  expect(dispatcher.receivedMessage).to.be.instanceof(CreateParty);
  expect(dispatcher.receivedMessage).to.deep.equal(new CreateParty("name"));
});