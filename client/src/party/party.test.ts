import test from "ava";
import type { AxiosResponse } from "axios";
import { HttpPostMock } from "../http/post.test-utils";
import { CreateParty } from "../messages/commands";
import { AsyncDispatcherMock } from "../messages/dispatcher.test-utils";
import { PartyCreated } from "../messages/events";
import type { CreatePartyResponse } from "./party";
import { Party } from "./party";

test(`Create Party`, async t => {
  const http = new HttpPostMock<CreatePartyResponse>(Promise.resolve({data:{code: "testCode"}} as AxiosResponse<CreatePartyResponse>));
  const dispatcher = new AsyncDispatcherMock();
  
  const party = new Party(http, dispatcher);
  party.consume(new CreateParty("testName"));
  
  await dispatcher.isDone;
  
  t.deepEqual(http.givenUrl, "/party/create");
  t.deepEqual(http.givenData, {name: "testName"});
  t.deepEqual(dispatcher.receivedMessage, new PartyCreated("testCode"));
});
