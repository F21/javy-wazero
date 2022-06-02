module.exports = {
    mode: "production",
    target: "es2019",
    optimization: {
        sideEffects: true,
    },
    entry: "./src/greet.js",
    output: {
        globalObject: "this",
        filename: "greet.js",
        path: "/tmp",
        libraryTarget: "umd",
        library: "Suborbital", // We need so that the Suborbital javy compiler can see the run_e function under the Suborbital namespace
        chunkFormat: "array-push",
    },
    performance: {
        hints: false,
    },
};