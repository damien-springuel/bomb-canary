import { CloseIdentity } from "../messages/commands";
import type { Dispatcher } from "../messages/dispatcher";
import type { StoreValues } from "../store/store";

export class IdentityService{
  constructor(readonly dispatcher: Dispatcher, readonly storeValues: StoreValues){}

  closeIdentity(): void {
    this.dispatcher.dispatch(new CloseIdentity());
  }

  isPlayerIsASpy(): boolean {
    return this.storeValues.revealedSpies.has(this.storeValues.player);
  }

  otherSpies(): string[] {
    const otherSpies = new Set<string>(this.storeValues.revealedSpies);
    otherSpies.delete(this.storeValues.player);
    return Array.from(otherSpies);
  }
}