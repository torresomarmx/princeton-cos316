# COS316, Assignment 1: Socket Programming

## Due: September 24 at 23:00

## Socket Programming

Socket programming is the standard way to write programs
that communicate over a network. While originally developed for Unix computers
programmed in C, the socket abstraction is general and not tied to any specific
operating system or programming language. This allows programmers to use the socket
mental model to write correct network programs in many contexts.

This part of the assignment will give you experience with basic socket programming.
You will write two programs for transmitting text messages over the Internet: a
client for sending and a server for receiving.
Both the client and server must be written in Go.

The client and server programs should meet the following specifications.
Be sure to read these meticulously before and after programming to make sure your
implementation fulfills them:

### Server specification
* The server program should listen on a socket, wait for a client to connect,
  receive a message from the client, print the message to stdout, and then wait
  for the next client indefinitely.
* The server should take one command-line argument: the port number to listen on
  for client connections.
* The server should accept and process client communications in an infinite loop,
  allowing multiple clients to send messages to the same server. The server should
  only exit in response to an external signal (e.g. SIGINT from pressing `ctrl-c`).
* The server should gracefully handle error values potentially returned by socket
  programming library functions (see specifics below). Errors related to handling
  client connections should not cause the server to exit after handling the error;
  all others should.
* The server should maintain a client queue and handle multiple client connection
  attempts sequentially. The good news is, this is the default behavior if you are
  using `net.Listen()`. **You do not need to do anything extra to satisfy this
  requirement.** In real applications, a TCP server would fork a new process to handle
  each client connection concurrently, but that is **not necessary** for this
  assignment.

### Client specification

* The client program should contact a server, read a message from stdin, send the
  message, and exit.
* The client should read and send the message *exactly* as it appears in stdin
  until reaching an EOF (end-of-file).
* The client should take two command-line arguments: the IP address of the server
  and the port number of the server.
* The client must be able to handle arbitrarily large messages by iteratively
  reading and sending chunks of the message, rather than reading the whole message
  into memory first.
* The client should handle partial sends (when a socket only transmits part of
  the data given in the last send operation) by attempting to re-send the rest of
  the data until it has all been sent.
* The client should gracefully handle error values potentially returned by socket
  programming library functions.

### Error Handling

Generally speaking, there are several reasonable actions that a program might take
upon realizing that it has encountered an error, and you may (or may not) sometimes
need to take one or more of these actions.

* **Attempt to recover:**
  Some errors may arise due to chance events like a busy or noisy network, and in these
  cases it is possible (and desirable) to try to recover gracefully, perhaps by trying
  the exact same operation again, or by tweaking some values first and then retrying.
* **Crash:**
  On the other hand, some errors cannot be recovered from at runtime. If the user
  requests access to some resource that is already being used by another process,
  there is no straightforward way to recover, and crashing (with an informative
  message) would be an acceptable response.
* **Print error message:**
  Especially for fatal errors that cause your program to crash, it is good style
  to print out a message indicating what has gone wrong. For non-fatal errors, you
  may find it useful to print messages for debugging purposes. For your final
  submit, you should make an effort to minimize output by commenting out any
  debugging statements.
* **Do nothing:**
  It is generally poor style to leave potential errors unhandled, as your program
  might continue executing, believing everything to be OK, only to crash later on
  in a way that will be much harder to debug. Your program should make an effort
  to handle all reasonable errors that may arise.

#### Error Handling in Go

Go has several error handling functions that may be of use to you:

* `log.Fatal(message string)`<sup>[(docs)](http://golang.org/pkg/log/#Fatal)</sup>:
  Print message to `os.Stderr` and terminate the program with a return code of 1.

* `log.Print(message string)`<sup>[(docs)](https://golang.org/pkg/log/#Print)</sup>:
  Prints message to stderr.

* `log.Panic(message string)`<sup>[(docs)](https://golang.org/pkg/log/#Panic)</sup>:
  Prints the error message, and then calls `panic()`, which propagates the error, and
  prints a stack trace if unhandled. `panic()` is similar to Java's `throw`.
  It differs from `log.Fatal()` in that deferred functions are executed before the
  program exits (perhaps freeing resources or flushing buffers).
  An interesting note is that callers can recover from a panic using `recover()`
  (analagous to `catch` in Java), but you will *not* need to make use of this
  functionality for this assignment.
  See [this blog post](https://blog.golang.org/defer-panic-and-recover)
  for more on `defer()`, `panic()`, and `recover()`.

### Getting started

You can work directly on your laptop if you installed Go and Git on your computer.  Alternatively,
you may do building and testing on the Vagrant VM. You may either write your code on
the Vagrant VM (both Emacs and Vim text editors are pre-installed) or directly on
your OS (allowing you to use any editor you have installed).

Many modern editors (Emacs, Atom, VS Code, Sublime, ...) provide convenient extensions
specifically designed for working with Golang. These extensions provide many
useful features, including:
* Automatically adding required import statements
* Automatically vetting your code for compilation errors
* Automatically formatting your code in idiomatic Go style
* Automatically running any unit tests you've written
It is strongly recommended that you install an extension appropriate for your
editor of choice to streamline your Go programming experience.

Before jumping into writing code, you will need to clone your code repository
from GitHub to your local machine.
The basic command to do this is `git clone https://github.com/cos316/<repository name>`.
You may also clone your repository using the GitHub Desktop client.
For more detailed instructions on cloning a repository, consult [this readme](https://github.com/cos316/COS316-Public/blob/master/assignments/GITHUB.md).

You can work with your code repository either using the VM or using your physical machine directly.

We have provided scaffolding code in the `assignment1` directory.
*You should read and understand this code before starting to program.*

You should program only in the locations of the provided files marked with `TODO`
comments. There is one `TODO` section for the client and one for the server.

You can add functions if you wish, but do not change file names, as they will be
used for automated testing.

The following section provides details for implementing the client and server programs in Go.

### Go

The documentation for Go socket programming is located [here](https://golang.org/pkg/net/).  

The overview at the top and the section on the [Conn type](https://golang.org/pkg/net/#Conn) will be most relevant.
You may also find the buffered [Reader](https://golang.org/pkg/bufio/#Reader)
and [Writer](https://golang.org/pkg/bufio/#Writer) types to be useful, but you
aren't required to use them, and you can construct a working solution without them.

The Go language (Golang) documentation can be cryptic, so be sure to familiarize
yourself with the language a bit first, especially if you are new to Go. You may
find the [Tour of Go](https://tour.golang.org/list) documentation useful if you
have never used Go before.

The files `client.go` and `server.go` contain the scaffolding code. You will need
to add socket programming code in the locations marked `TODO`. The reference
solutions have roughly 40  (well commented and spaced) lines of code in the
`TODO` sections of each file. Your implementations may be shorter or longer.

You should build your solution by running `make` in the `assignment1` directory.
Your code *must* build using the provided Makefile.
The server should be run as `./server [port] > [output file]`.
The client should be run as `./client [server IP] [server port] < [message file]`.
See "Testing" for more details.

### Testing

You should test your implementation by attempting to send messages from your
client to your server.
The server can be run in the background (append a `&` to the command) or in a separate
SSH window.
You should use `127.0.0.1` (i.e. the "localhost", or "loopback", address) as the
server IP and a high server port number between 10000 and 60000.

You can kill a background server with the command `fg` to bring it to the foreground
then `ctrl-c`. Conversely, you can send a foreground process to the background by
hitting `ctrl-z` to suspend the process, and typing the command `bg` to resume the
process in the background.

You should test your implementation by attempting
to send several different messages between your client and server.
For example:

0.  The short message "Go Tigers!\n"
0.  A long, randomly generated alphanumeric message
0.  A long, randomly generated binary message
0.  Several short messages sent sequentially from separate clients to one server
0.  Several long, random alphaumeric messages sent concurrently from separate
    clients to one server


### Debugging hints

Here are some debugging tips. If you are still having trouble, ask a question on
Piazza or see an instructor during office hours.

* There are defined buffer size constants in the scaffolding code. Use them.
  If you are not using one of them, either you have hard-coded a value, which is
  bad style, or you are very likely doing something wrong.
* There are multiple ways to read and write from stdin/stdout in Go. Any method
  is acceptable as long as it does not read an unbounded amount into memory at
  once and does not modify the message.
* If you are using buffered I/O to write to stdout, make sure to call `flush` or
  the end of a long message may not write.
* Remember to close the socket at the end of the client program.
* When testing, make sure you are using `127.0.0.1` as the server IP argument to
  the client and the same server port for both client and server programs.
* If you get "address already in use" errors, make sure you don't already have a
  server running. Otherwise, restart your ssh session with the command `logout`
  followed by `vagrant ssh`.
* If you are getting other connection errors, try a different port between
  10000 and 60000.

### Q&A

* **Do I need to handle signals such as SIGINT to clean up the server process when
  the user presses `ctrl-c`?**

  No, it is not necessary in this assignment. The default response to signals is good enough.
  Keep in mind it would be good practice to do so in general, however.

* **Should I use stream (TCP) or datagram (UDP) sockets?**

  Please use stream (TCP) sockets, to ensure that the exact message is delivered.
  Streams ensure reliable, in-order packet transmission, whereas datagram packets
  are not guaranteed to be delivered.

* **Should I support IPv6?**

  IPv6 support is not necessary. Your code will be tested using IPv4 only.

* **Should the client wait to receive a reply from the server?**

  No, in this assignment it should exit immediately after sending all the data.

* **Should the server handle client connections concurrently (in separate processes)?**

  No, as stated in the client specification, this is not required in this assignment.
  So no need to use `syscall.ForkExec()` or any other system-level `fork()`!

### Submission and grading

To submit your client and server, `git commit` the changes to your code and
`git push` them to your GitHub classroom repository.

See the [GitHub classroom README](https://github.com/cos316/COS316-Public/blob/master/assignments/GITHUB.md)
for more detailed instructions.

We will grade your submissions by compiling your client and server, and then
sending messages back and forth between each of your submitted programs and a
reference server or client, as appropriate. Your code will be scored based on
how many different kinds of messages are transmitted correctly, and how well
your implementations adhere to various aspects of the specification. Within a
couple minutes of submitting your assignment, the GitHub autograder will add a
comment to your most recent commit on GitHub, indicating your test results.

You may submit in this way and receive feedback as many times as you like,
whenever you like, but a lateness penalty will be applied to submissions
received after the deadline.

Code that does not compile is graded harshly; if you want partial credit on code
that doesn't compile, comment it out and make sure your file compiles!

Remember that you should submit two files: a server file and a client file, both
written in Go!
