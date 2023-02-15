import { expect, test } from "vitest";
import type { AxiosResponse } from "axios";
import { HttpPostMock } from "../http/post.test-utils";
import { CreateParty, JoinParty } from "../messages/commands";
import { AsyncDispatcherMock } from "../messages/dispatcher.test-utils";
import { CreatePartySucceeded, JoinPartySucceeded} from "../messages/events";
import { Party } from "./party";

test(`Create Party`, async () => {
  const http = new HttpPostMock(Promise.resolve({data:{}} as AxiosResponse<{}>));
  const dispatcher = new AsyncDispatcherMock();
  
  const party = new Party(http, dispatcher);
  party.consume(new CreateParty("testName"));
  
  await dispatcher.isDone;
  
  expect(http.givenUrl).to.equal("/party/create");
  expect(http.givenData).to.deep.equal({name: "testName"});
  expect(dispatcher.receivedMessage).to.deep.equal( new CreatePartySucceeded());
});

test(`Join Party`, async () => {
  const http = new HttpPostMock(Promise.resolve({data:{}} as AxiosResponse<{}>));
  const dispatcher = new AsyncDispatcherMock();
  
  const party = new Party(http, dispatcher);
  party.consume(new JoinParty("testName", "testCode"));
  
  await dispatcher.isDone;
  
  expect(http.givenUrl).to.equal("/party/join");
  expect(http.givenData).to.deep.equal({name: "testName", code: "testCode"});
  expect(dispatcher.receivedMessage).to.deep.equal( new JoinPartySucceeded());
});
