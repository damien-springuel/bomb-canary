<script lang="ts">
import Axios from "axios";
import type {Message} from "./messagebus";
import {CreatePartyClicked} from "./messagebus";
export let dispatcher: {Dispatch: (message: Message) => void};
export let store: any;

let name: string;

function createParty() {
  dispatcher.Dispatch(new CreatePartyClicked(name))
  // Axios.post("http://localhost:44324/party/create", {name: name}, {withCredentials: true}).then(r => console.log(r));
}

</script>

<div class="flex flex-col items-center h-full bg-gray-900 text-blue-500 px-6 overflow-auto">
  <h1 class="text-5xl">
    Bomb Canary
  </h1>
  <h2 class="text-2xl text-center text-blue-400">
    A real-time app to play "The Resistance" board game with friends.
  </h2>
  <h3>like {$store.name}</h3>

  <div class="flex-grow flex flex-col justify-evenly w-full text-2xl items-stretch space-y-10 mb-4">
    <input type="text" placeholder="Name" class="bg-blue-400 focus:bg-blue-200 hover:bg-blue-300 placeholder-gray-700 rounded-lg px-2 text-gray-900 outline-none" bind:value={name}>
    <div class="border-solid border-t-2 border-blue-500"></div>
    <input type="text" placeholder="Party Code" class="bg-blue-400 focus:bg-blue-200 hover:bg-blue-300 placeholder-gray-700 rounded-lg px-2 text-gray-900 outline-none">
    <button class="rounded-full px-6 py-4 bg-blue-500 hover:bg-blue-300 text-gray-900 focus:outline-none">Join Party</button>
    <div class="text-center">or</div>
    <button class="rounded-full px-6 py-4 bg-blue-500 hover:bg-blue-300 text-gray-900 focus:outline-none" on:click={createParty}>Create Party</button>
  </div>
</div>