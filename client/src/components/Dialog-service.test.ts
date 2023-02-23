import {expect, test} from "vitest";
import { CloseDialog } from "../messages/commands";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { DialogService } from "./Dialog-service";

test("Close dialog", () => {
  let dispatcherMock = new DispatcherMock();
  const service = new DialogService(dispatcherMock);

  service.closeDialog();
  expect(dispatcherMock.receivedMessage).to.be.instanceof(CloseDialog);
  expect(dispatcherMock.receivedMessage).to.deep.equal(new CloseDialog());
});