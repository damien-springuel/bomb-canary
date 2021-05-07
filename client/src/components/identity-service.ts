import type { StoreValues } from "../store/store";

export class IdentityService{
  constructor(readonly storeValues: StoreValues){}

  isPlayerIsASpy(): boolean {
    return this.storeValues.revealedSpies.has(this.storeValues.player);
  }

  otherSpies(): string[] {
    const otherSpies = new Set<string>(this.storeValues.revealedSpies);
    otherSpies.delete(this.storeValues.player);
    return Array.from(otherSpies);
  }
}