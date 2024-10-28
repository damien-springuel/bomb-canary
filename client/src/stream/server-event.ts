export interface ServerEvent {
  PartyCreated?: {
    Code: string
  },
  
  PlayerConnected?: {
    Name: string
  },
  
  PlayerDisconnected?: {
    Name: string
  },
  
  PlayerJoined?: {
    Name: string
  },
  
  GameStarted?: {
    MissionRequirements: {
      NbPeopleOnMission: number,
      NbFailuresRequiredToFail: number,
    }[]
  },
  
  SpiesRevealed?: {
    Spies: {[name:string]:{}}
  },
  
  LeaderStartedToSelectMembers?: {
    Leader: string,
  },
  
  LeaderSelectedMember?: {
    SelectedMember: string,
  },
  
  LeaderDeselectedMember?: {
    DeselectedMember: string,
  },
  
  LeaderConfirmedSelection?: {},
  
  PlayerVotedOnTeam?: {
    Player: string,
    Approved?: boolean,
  },
  
  AllPlayerVotedOnTeam?: {
    Approved: boolean,
    VoteFailures: number,
    PlayerVotes: {[name:string]:boolean}
  },

  MissionStarted?: {},

  PlayerWorkedOnMission?: {
    Player: string,
    Success?: boolean,
  }

  MissionCompleted?: {
    Success: boolean,
    NbFails: number,
  }

  GameEnded?: {
    Winner: string,
    Spies: string[],
  }

  EventsReplayStarted?: {
    Player: string,
  }
  
  EventsReplayEnded?: {}
}