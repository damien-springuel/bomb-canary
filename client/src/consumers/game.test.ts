import test from "ava";
import { LeaderConfirmedTeam, LeaderDeselectedMember, LeaderSelectedMember, LeaderStartedToSelectMembers } from "../messages/events";
import { GameManager, GameStore } from "./game";

test(`Game Manager - LeaderStartedToSelectMembers`, t => {
  let receivedLeader: string
  const gameMgr = new GameManager({assignLeader: leader => {receivedLeader = leader}} as GameStore);
  gameMgr.consume(new LeaderStartedToSelectMembers("testLeader"));
  t.deepEqual(receivedLeader, "testLeader");
});

test(`Game Manager - LeaderSelectedMember`, t => {
  let receivedMember: string
  const gameMgr = new GameManager({selectPlayer: member => {receivedMember = member}} as GameStore);
  gameMgr.consume(new LeaderSelectedMember("member"));
  t.deepEqual(receivedMember, "member");
});

test(`Game Manager - LeaderDeselectedMember`, t => {
  let receivedMember: string
  const gameMgr = new GameManager({deselectPlayer: member => {receivedMember = member}} as GameStore);
  gameMgr.consume(new LeaderDeselectedMember("member"));
  t.deepEqual(receivedMember, "member");
});
