import type { Dispatcher } from "../messages/dispatcher";
import { 
  AllPlayerVotedOnTeam,
  EventsReplayEnded, 
  EventsReplayStarted, 
  GameEnded, 
  GameStarted, 
  LeaderConfirmedTeam, 
  LeaderDeselectedMember, 
  LeaderSelectedMember, 
  LeaderStartedToSelectMembers, 
  MissionCompleted, 
  MissionStarted, 
  PlayerConnected, 
  PlayerDisconnected, 
  PlayerJoined, 
  PlayerVotedOnTeam, 
  PlayerWorkedOnMission, 
  ServerConnectionClosed, 
  ServerConnectionErrorOccured, 
  SpiesRevealed 
} from "../messages/events";
import { Allegiance } from "../types/types";
import type { ServerEvent } from "./server-event";

export class Handler {

  constructor(private readonly dispatcher: Dispatcher){}

  onClose(): void {
    this.dispatcher.dispatch(new ServerConnectionClosed());
  }
  onError(): void {
    this.dispatcher.dispatch(new ServerConnectionErrorOccured());
  }
  onEvent(event: ServerEvent) {
    if (event.EventsReplayStarted) {
      this.dispatcher.dispatch(new EventsReplayStarted(event.EventsReplayStarted.Player));
    }
    else if (event.EventsReplayEnded) {
      this.dispatcher.dispatch(new EventsReplayEnded());
    }
    else if (event.PlayerConnected) {
      this.dispatcher.dispatch(new PlayerConnected(event.PlayerConnected.Name));
    }
    else if (event.PlayerDisconnected) {
      this.dispatcher.dispatch(new PlayerDisconnected(event.PlayerDisconnected.Name));
    }
    else if (event.PlayerJoined) {
      this.dispatcher.dispatch(new PlayerJoined(event.PlayerJoined.Name));
    }
    else if (event.GameStarted) {
      const req = event.GameStarted.MissionRequirements
        .map(r => ({nbPeopleOnMission: r.NbPeopleOnMission, nbFailuresRequiredToFail: r.NbFailuresRequiredToFail}))
      this.dispatcher.dispatch(new GameStarted(req));
    }
    else if (event.SpiesRevealed) {
      this.dispatcher.dispatch(new SpiesRevealed(new Set<string>(Object.keys(event.SpiesRevealed.Spies || []))));
    }
    else if (event.LeaderStartedToSelectMembers) {
      this.dispatcher.dispatch(new LeaderStartedToSelectMembers(event.LeaderStartedToSelectMembers.Leader));
    }
    else if (event.LeaderSelectedMember) {
      this.dispatcher.dispatch(new LeaderSelectedMember(event.LeaderSelectedMember.SelectedMember));
    }
    else if (event.LeaderDeselectedMember) {
      this.dispatcher.dispatch(new LeaderDeselectedMember(event.LeaderDeselectedMember.DeselectedMember));
    }
    else if (event.LeaderConfirmedSelection) {
      this.dispatcher.dispatch(new LeaderConfirmedTeam());
    }
    else if (event.PlayerVotedOnTeam) {
      const approved = typeof event.PlayerVotedOnTeam.Approved === 'boolean' ? event.PlayerVotedOnTeam.Approved : null
      this.dispatcher.dispatch(new PlayerVotedOnTeam(event.PlayerVotedOnTeam.Player, approved));
    }
    else if (event.AllPlayerVotedOnTeam) {
      const playerVote = new Map<string,boolean>(Object.keys(event.AllPlayerVotedOnTeam.PlayerVotes).map(k => [k, event.AllPlayerVotedOnTeam.PlayerVotes[k]]));
      this.dispatcher.dispatch(new AllPlayerVotedOnTeam(event.AllPlayerVotedOnTeam.Approved, playerVote));
    }
    else if (event.MissionStarted) {
      this.dispatcher.dispatch(new MissionStarted());
    }
    else if (event.PlayerWorkedOnMission) {
      const success = typeof event.PlayerWorkedOnMission.Success === 'boolean' ? event.PlayerWorkedOnMission.Success : null
      this.dispatcher.dispatch(new PlayerWorkedOnMission(event.PlayerWorkedOnMission.Player, success));
    }
    else if (event.MissionCompleted) {
      this.dispatcher.dispatch(new MissionCompleted(event.MissionCompleted.Success, event.MissionCompleted.NbFails));
    }
    else if (event.GameEnded) {
      const winner = event.GameEnded.Winner === "spy" ? Allegiance.Spies : Allegiance.Resistance
      this.dispatcher.dispatch(new GameEnded(winner, new Set<string>(event.GameEnded.Spies)));
    }
  }
}