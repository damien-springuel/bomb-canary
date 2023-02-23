export interface IdentityValues {
  readonly player:string,
  readonly revealedSpies: Set<string>,
}

export class IdentityService{
  constructor(readonly identityValues: IdentityValues){}

  isPlayerIsASpy(): boolean {
    return this.identityValues.revealedSpies.has(this.identityValues.player);
  }

  otherSpies(): string {
    const otherSpies = new Set<string>(this.identityValues.revealedSpies);
    otherSpies.delete(this.identityValues.player);
    return Array.from(otherSpies).join(", ");
  }
}