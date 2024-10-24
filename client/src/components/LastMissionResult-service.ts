export interface LastMissionResultValues {
  readonly success: boolean
}

export class LastMissionResultService {
  constructor(readonly lastMissionResultValues: LastMissionResultValues) {}

  get success(): boolean {
    return this.lastMissionResultValues.success;
  }
}