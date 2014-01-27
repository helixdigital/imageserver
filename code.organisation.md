There are three directories and the main binary file.

* `/entities` basic data elements
* `/core` use cases
* `/plugin` concrete implementions of the interfaces the use cases employ 
* `main.go` instantiates plugins and starts the server.

The organisation of the code is such that no code knows anything about other modules that are lower on this list than itself.
