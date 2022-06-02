function greet(input) {
    return { greeting: "Hello " + input.name + "!" }; // Javy will automatically convert the output to JSON
}

// In order for the function to be exported properly, we need to put it Shopify.main
Shopify = {
    main: greet,
};