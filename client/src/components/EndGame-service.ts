import { Allegiance } from "../types/types";

export interface EndGameValues {
  readonly player: string,
  readonly winner: Allegiance,
  readonly spies: Set<string>,
}

export class EndGameService {
  constructor(readonly endGameValues: EndGameValues) {}
  
  get spiesAsString(): string {
    return Array.from(this.endGameValues.spies).join(", "); 
  }

  get spiesHaveWon(): boolean {
    return this.endGameValues.winner === Allegiance.Spies;
  }
  
  get playerHasWon(): boolean {
    return this.playerAllegiance === this.endGameValues.winner;
  }

  private get playerAllegiance(): Allegiance {
    return this.endGameValues.spies.has(this.endGameValues.player) ? 
      Allegiance.Spies : 
      Allegiance.Resistance;
  }
}