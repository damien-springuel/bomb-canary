import test from "ava";
import { Store, StoreValues } from "./store";
import {get} from "svelte/store";

test(`Store`, t => {
  const store = new Store();
  const storeValues: StoreValues = get(store);
  t.deepEqual(storeValues, {} as StoreValues);
});

test(`Store - setname`, t => {
  const store = new Store();
  store.setName("test");
  const storeValues: StoreValues = get(store);
  t.is(storeValues.name, "test");
});