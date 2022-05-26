# Javy - Wazero example

To run: `go run main.go`

To get it to work properly, all the code including the instantiation of the wazero runtime needs to be moved into the `for`
loop. Meaning that we need to create a new runtime and compile the module everytime we want to call it.
