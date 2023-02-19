export interface TeamSelectionValues {
  readonly currentTeam: Set<string>,
  readonly player: string,
  readonly leader: string,
  readonly currentTeamVoteNb: number,
  readonly players: string[],
}

export class TeamSelectionService {
  constructor(readonly values: TeamSelectionValues){}

  isGivenPlayerInTeam(player: string): boolean {
    return this.values.currentTeam.has(player);
  }

  get isPlayerTheLeader(): boolean {
    return this.values.player === this.values.leader;
  }

  get leader(): string {
    return this.values.leader;
  }

  get currentTeamVoteNb(): number {
    return this.values.currentTeamVoteNb;
  }
  
  get players(): string[] {
    return this.values.players;
  }
}