import test from "ava";
import type { AxiosResponse } from "axios";
import { HttpPostMock } from "../http/post.test-utils";
import { CreateParty, JoinParty } from "../messages/commands";
import { AsyncDispatcherMock } from "../messages/dispatcher.test-utils";
import { CreatePartySucceeded, JoinPartySucceeded} from "../messages/events";
import { Party } from "./party";

test(`Create Party`, async t => {
  const http = new HttpPostMock<{}>(Promise.resolve({data:{}} as AxiosResponse<{}>));
  const dispatcher = new AsyncDispatcherMock();
  
  const party = new Party(http, dispatcher);
  party.consume(new CreateParty("testName"));
  
  await dispatcher.isDone;
  
  t.deepEqual(http.givenUrl, "/party/create");
  t.deepEqual(http.givenData, {name: "testName"});
  t.deepEqual(dispatcher.receivedMessage, new CreatePartySucceeded());
});

test(`Join Party`, async t => {
  const http = new HttpPostMock<{}>(Promise.resolve({data:{}} as AxiosResponse<{}>));
  const dispatcher = new AsyncDispatcherMock();
  
  const party = new Party(http, dispatcher);
  party.consume(new JoinParty("testName", "testCode"));
  
  await dispatcher.isDone;
  
  t.deepEqual(http.givenUrl, "/party/join");
  t.deepEqual(http.givenData, {name: "testName", code: "testCode"});
  t.deepEqual(dispatcher.receivedMessage, new JoinPartySucceeded());
});
