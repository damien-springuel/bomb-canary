import test from "ava";
import { HttpPostMock } from "../http/post.test-utils";
import { CreateParty } from "../messages/commands";
import { Party } from "./party";

test(`Create Party`, t => {
  const httpMock = new HttpPostMock();
  const party = new Party(httpMock);
  party.consume(new CreateParty("testName"));
  t.deepEqual(httpMock.givenUrl, "/party/create");
  t.deepEqual(httpMock.givenData, {name: "testName"});
});
