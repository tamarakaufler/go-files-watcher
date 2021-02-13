# Synopsis

Go implementation of a daemon for montoring file changes and running a command when change is detected.

The configurable options are:

|                |                  default                          |
-----------------|:-------------------------------------------------:|
|  BasePath      |     directory where the watcher daemon is started |
|  Extension     |     .go (currently only one)                      |
|  Command       |     tree (command to run upon detected change)    |
|  Excluded      |     none ( a list of strings/regexes)             |
|  Frequency     |     5sec (repeat of the check)                    |

# Implementation

There are 3 progressive implementations, from the initial one using directly filepath.Walk,
an intemediate one as a preparation for the third parallelized third implementation. First two
versions are commented out ((*Daemon).Watch method).

## Parellelized implementation

Uses filepath.Walk to collect files that are under watch, ie excluding those that:
  * dont have the required extension
  * are configured to be excluded

The resulting files are looped through, each processed in a separate goroutine to check whether the file
has been modified since the last check. When a change has been detected a message is sent to a channel to
stop further looping and to stop already spun up gouroutines.

# TODO

* tests
* Makefile
* golanci-lint

* implement version 2 suitable for running in a Docker container
