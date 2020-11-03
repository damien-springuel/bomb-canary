import { 
  EventsReplayEnded, 
  EventsReplayStarted, 
  LeaderConfirmedTeam, 
  LeaderDeselectedMember, 
  LeaderSelectedMember, 
  LeaderStartedToSelectMembers, 
  PartyCreated, 
  PlayerConnected, 
  PlayerDisconnected, 
  PlayerJoined, 
  ServerConnectionClosed, 
  ServerConnectionErrorOccured, 
  SpiesRevealed 
} from "../messages/events";
import type { Message } from "../messages/messagebus";
import type { ServerEvent } from "./server-event";

export class Handler {

  constructor(private readonly dispatcher: {dispatch: (m: Message) => void}){}

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
    if (event.EventsReplayEnded) {
      this.dispatcher.dispatch(new EventsReplayEnded());
    }
    else if (event.PartyCreated) {
      this.dispatcher.dispatch(new PartyCreated(event.PartyCreated.Code));
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
  }
}