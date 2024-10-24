import { expect, test } from "vitest";
import { LastMissionResultService, type LastMissionResultValues } from "./LastMissionResult-service";

test(`last mission success`, () => {
  let service = new LastMissionResultService({success: false} as LastMissionResultValues);
  expect(service.success).to.be.false;
  
  service = new LastMissionResultService({success: true} as LastMissionResultValues);
  expect(service.success).to.be.true;
});