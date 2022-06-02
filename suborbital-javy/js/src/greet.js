import "fastestsmallesttextencoderdecoder-encodeinto/EncoderDecoderTogether.min.js";
import { setup, runnable } from "@suborbital/runnable";

const decoder = new TextDecoder();

function run_e(payload, ident) {
    // Imports will be injected by the runtime
    setup(this.imports, ident);

    const input = decoder.decode(payload);

    const decodedJSON = JSON.parse(input); // Decode the JSON because Suborbital's Javy does not do this for us
    const result = { greeting: "Hello " + decodedJSON.name + "!" };
    const encodedJSON = JSON.stringify(result); // Encode the JSON because Suborbital's Javy does not do this for us

    runnable.returnResult(encodedJSON);
}

export { run_e }; // The exported must be called run_e in order for Suborbital's Javy to pick it up