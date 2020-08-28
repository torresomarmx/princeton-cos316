# COS316, Assignment 3: In-Memory Cache

## Due: October 22 at 23:00

# In-Memory Cache


In this project, you will build an in-memory, look-aside, write-allocate cache.
You will implement two eviction algorithms: first-in first-out (FIFO) and
least-recently-used (LRU). Each of these implementations adhere to an abstract
interface provided to you, called `Cache`.

## API

A cache as a fixed-size store of key-value bindings. In your cache, keys are
strings (e.g., the name of a file or variable), and values are arbitrary-length
byte slices.

Because the cache is fixed-size, it *cannot* grow to accept arbitrarily many
keys and values. We initialize a cache with a `limit` parameter that defines
how precisely how many bytes worth of keys and values it can accommodate.

Instead of growing, once the cache becomes too full, we must **evict**
some item that is already in the cache before we are able to **admit** new
values to the cache.

The mechanism by which we decide which items to remove from the cache is
known as an **eviction algorithm**, which we will cover later in this document.

A general-purpose `Cache` interface is defined by the following functions.
Any type that implements these functions is automatically able to be used as
a `Cache` type in Golang.

This interface is defined in `cache.go`. You should examine this file and
understand its purpose.

```go
type Cache interface {
	// Get returns the value associated with the given key, if it exists.
	// This operation counts as a "use" for that key-value pair
	// ok is true if a value was found and false otherwise.
	Get(key string) (value []byte, ok bool)

	// Remove removes and returns the value associated with the given key, if it exists.
	// ok is true if a value was found and false otherwise
	Remove(key string) (value []byte, ok bool)

	// Set associates the given value with the given key, possibly evicting values
	// to make room. Returns true if the items was added successfully, else false.
	Set(key string, value []byte) bool

	// Len returns the number of items in the cache.
	Len() int

	// MaxStorage returns the maximum number of bytes this cache can store
	MaxStorage() int

	// RemainingStorage returns the number of unused bytes available in this cache
	RemainingStorage() int


	// Stats returns a pointer to a Stats object that indicates how many hits
	// and misses this cache has resolved over its lifetime.
	Stats() *Stats
}
```

### First-in First-out (FIFO) Eviction

For the first part of the assignment, you will be implementing a cache using a
first-in first-out eviction algorithm.

This will require you to modify the file `fifo.go`.  Unit testing should be
done in `fifo_test.go`.

FIFO eviction means that if there is not enough space in the cache for a new
item, the cache evicts items one at a time until there is enough
room, starting with the item that was added to the cache first (i.e. the
oldest item presently in the cache), and proceeding in that same order.

### Least Recently Used (LRU) Eviction

For the second part of the assignment, you will be implementing a cache using a
least recently used eviction algorithm.

This will require you to modify the file `lru.go`.
Unit testing should be done in `lru_test.go`.

LRU eviction means that if there is not enough space in the cache for a new
item, the cache evicts items one at a time until there is enough room,
starting with the item that was *used* by a client least recently. A item
is considered *user* any time is the object of a `Get()` or `Set()` call.

## Additional Specifications

* *What sorts of keys and values are acceptable?*
  Keys can be any valid Go string and items
  Note that the empty string is an acceptable key for an item. Likewise,
  the empty byte slice is an acceptable value. The `nil` byte slice is also an
  acceptable value for an item.

* *What does the `ok` return value signify?*
  Your cache should use `ok=true` to indicate that the requested operation
  executed successfully, or `ok=false` to indicate some issue.
  For example, if `Get()` or `Remove()` returns `ok=false`, it may mean that
  no item exists for the requested key.

* *If `ok` is false, what should the other returned value be?*
  If `Get()` or `Remove()` returns `false` as its second return value, clients
  should assume that the first return value is invalid, and its specific value
  is therefore not relevant. In this case, it would be reasonable to return `nil`,
  but you are not required to do so.

* *How much memory does an item consume?*
  For this assignment, assume that the memory consumed by an item is precisely
  `len(key) + len(value)`. In practice, simply counting the number of bytes in a
  string and the number of entries in a byte array is not sufficient since the
  data structures you use to store the items almost certainly incur overhead,
  such as the size of the pointers to a key or value. However, we ignore these
  factors for the assignment to make testing simpler.

* *What if an item could never fit into a cache?* You may encounter situations
  where a client requests adding an item that is larger than the maximum
  capacity of the cache. In these cases, `Set()` should return `false` to
  indicate the binding was not admitted to the cache, and the contents of the
  cache should be left unmodified.


## Performance

* Your implementation should be memory-efficient in the sense that it evicts
  values from the cache only as a last resort. If it is possible to store a
  binding in the cache without evicting another, your implementation must do so.

* Your implementation should be time-efficient. Specifically, `Get`, `Set`,
  `Len()` and `Stats()` should be constant-time in the number of items in the
  cache. Carefully consider the data structures you use to implement FIFO and
  LRU.  For example, iterating over all items in the cache to find the
  oldest/least-recently used will *not* be acceptably fast when the cache is
  large.

## A Note on External Libraries

* Your code may use any data structures that have been implemented in the Go
  standard libraries, or any data structures that you implement yourself from
  scratch, but you may not use data structures defined in third party
  libraries. Your code must not rely on any existing LRU or FIFO implementation,
  regardless of where it came from.

## Getting started

As in previous assignments, you will need to clone your GitHub classroom
repository, and add the downloaded repo as a synced folder in your Vagrant VM
before you start programming.
Refer to the [GitHub classroom README](https://github.com/cos316/COS316-Public/blob/master/assignments/GITHUB.md)
for more detailed instructions.

## Unit Testing

Recall Go uses the [testing package](https://golang.org/pkg/testing/) to create
unit tests for Go packages.

For this assignment, you are provided with several files:
* `helpers_test.go` contains useful helper functions you may use (or modify)
  for debugging purposes, or to help creating your own tests.
* `fifo_test.go` contains a very basic unit test for your first-in first-out
  cache. You are encouraged to extend this file to create your own unit tests.
* `lru_test.go` contains a very basic unit test for your least recently used
  cache. You are encouraged to extend this file to create your own unit tests.

Read through all three files and try to understand how and why they work the
way that they do.  Hopefully, it will give you some ideas you can build off of
to create more comprehensive tests.

You can run your unit tests with the command `go test`, which simply reports the
result of the test, and the reason for failure, if any, or you may add the `-v`
flag to see the verbose output of the unit tests.

## Submission & Grading

Your assignment will be automatically submitted every time you push your changes
to your GitHub Classroom repo. Within a couple minutes of your submission, the
autograder will make a comment on your commit listing the output of our testing
suite when run against your code. **Note that you will be graded only on your
changes to the `cache` package**, and not on your changes to any other files,
though you may modify any files you wish.

You may submit and receive feedback in this way as many times as you like,
whenever you like, but a substantial lateness penalty will be applied to
submissions past the deadline.
