import test from "ava";
import { LeaderStartedToSelectMembers } from "../messages/events";
import { GameManager } from "./game";

test(`Game Manager - LeaderStartedToSelectMembers`, t => {
  let receivedLeader: string
  const gameMgr = new GameManager({assignLeader: leader => receivedLeader = leader});
  gameMgr.consume(new LeaderStartedToSelectMembers("testLeader"));
  t.deepEqual(receivedLeader, "testLeader");
});
