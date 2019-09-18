SLog - Golang Super Simple Logging
===================================

[![Apache License](https://img.shields.io/badge/License-Apache-brightgreen.svg)](https://tldrlegal.com/license/apache-license-2.0-\(apache-2.0\)) [![Coverage Status](https://coveralls.io/repos/github/quan-to/slog/badge.svg?branch=master)](https://coveralls.io/github/quan-to/slog?branch=master) [![Build Status](https://travis-ci.org/quan-to/slog.svg?branch=master)](https://travis-ci.org/quan-to/slog)

Usage:

```go
package main

import (
    "github.com/quan-to/slog"
    "time"
)

var log = slog.Scope("MAIN")


func Call0(i slog.Instance, arg0 string) {
    l := i.SubScope("Call0").WithFields(map[string]interface{}{
        "arg0": arg0,
    })
    l.Await("Doing some work")
    time.Sleep(time.Second)
    l.Done("Finished some work")
    l.Note("Not sure what I'm doing...")
    l.Info("Calling Call1")
    Call1(l, "call1arg")
    l.Done("Exiting")
}

func Call1(i slog.Instance, huebr string) {
    l := i.SubScope("Call1").WithFields(map[string]interface{}{
        "huebr": huebr,
    })
    l.Info("Calling Call2")
    Call2(l, "abcde")
    l.Warn("Call 1 finished")
}

func Call2(i slog.Instance, pop string) {
    l := i.SubScope("Call2").WithFields(map[string]interface{}{
        "pop": pop,
    })

    l.IO("Doing some IO")
    l.Error("I'm useless. Please fix-me")
}

func main() {
    slog.SetScopeLength(40) // Expand Scope pad length

    log = log.Tag("REQ001") // Tag current as REQ001

    log.Info("Starting program")

    Call0(log, "MyArg0")

    Call1(log, "Call1Arg")
    Call2(log, "Call2Arg")
}
```

Output:

![Sample Output](https://user-images.githubusercontent.com/578310/64198701-289b6b80-ce5f-11e9-8771-88ae4e07a213.png)


### Log Pattern

The slog output is expected to be like this:

```
2019-09-16T15:35:52-03:00 | I | IO    | REQ001 | MAIN > Call0 > Call1 > Call2             | test.go:38 | Doing some IO  | {"arg0":"MyArg0","huebr":"call1arg","pop":"abcde"}
```

There are some fields that are optional (like filename/line number) but it should always follow the same pattern:
```
DATETIME | LEVEL | OPERATION | TAG | SCOPE | [FILENAME:LINE NUMBER] | MESSAGE | LOG FIELDS
```

*   `DATETIME` => An ISO Datetime when the log is displayed
*   `LEVEL` => The level of the log line
    *   `I` => INFO - Shows an information usually to track what's happening inside an application
    *   `W` => WARN - Shows an warning regarding something that went in a way that might require some attention
    *   `E` => ERROR - Shows an application error that can be expected or not
    *   `D` => DEBUG - Shows some debug information to help tracking issues
    *   `F` => FATAL - Shows an error that will quit the application in that point
*   `TAG` => Line log tag. Use this for tracking related log lines. For example with a HTTP Request ID
*   `SCOPE` => The scope of the current log. Use this to trace the chain of calls inside the application. For example in a context change
*   `FILENAME: LINE NUMBER` => *OPTIONAL* When ShowLines is enabled, it will show the filename and the line number of the caller of the slog library. Use this on debug mode to see which piece of code called the log library. Disabled by default
*   `MESSAGE` => The message
*   `LOG FIELDS` => When an instance is created using `WithFields` call, the fields will be serialized to either JSON or Key-Value depending on the configuration of the log instance. Defaults to JSON

#### Operations Usage

The library implements the concept of operation type. This describes which type of operation the log line represents.

*   `IO` => Anytime you do a I/O Operation such as Database or File Read/Write
*   `AWAIT` => Anytime you will do a asynchronous operation that will block the flow
*   `DONE` => After any `AWAIT` operation you should always use a `DONE` to log that the operation that did the `AWAIT` log line has finished.
*   `NOTE` => Also know as Verbose, that operation is used when you're just letting the log reader some note about the operation
*   `MSG` => Common messages such like normal information (which does not fit in other operations)

There are some syntax sugars to make easy to use Log Levels with Log Operations together:

* Warning Messages
    *   `WarnDone` => Same as `log.Operation(DONE).Warn(message)`
    *   `WarnNote` => Same as `log.Operation(NOTE).Warn(message)`
    *   `WarnAwait` => Same as `log.Operation(AWAIT).Warn(message)`
    *   `WarnSuccess` => Same as `log.Operation(DONE).Warn(message)`
    *   `WarnIO` => Same as `log.Operation(IO).Warn(message)`
* Error Messages
    *   `ErrorDone` => Same as `log.Operation(DONE).Error(message)`
    *   `ErrorNote` => Same as `log.Operation(NOTE).Error(message)`
    *   `ErrorAwait` => Same as `log.Operation(AWAIT).Error(message)`
    *   `ErrorSuccess` => Same as `log.Operation(DONE).Error(message)`
    *   `ErrorIO` => Same as `log.Operation(IO).Error(message)`
* Debug Messages
    *   `DebugDone` => Same as `log.Operation(DONE).Debug(message)`
    *   `DebugNote` => Same as `log.Operation(NOTE).Debug(message)`
    *   `DebugAwait` => Same as `log.Operation(AWAIT).Debug(message)`
    *   `DebugSuccess` => Same as `log.Operation(DONE).Debug(message)`
    *   `DebugIO` => Same as `log.Operation(IO).Debug(message)`

Use these syntax sugars whenever is possible, instead of calling `Operation(X).Level` directly.

### Multiline Logs

If a multiline log is displayed, the library will correctly ident all the messages:

```
2019-09-16T15:39:42-03:00 | I | MSG   | REQ001 | MAIN  | Multiline call
                                                         Thats the second line
                                                         Thats the third line  | {}
```

### Use Patterns

*   After calling a `Await`, you should always call a `Done` or `Success`
*   Don't add extensive fields to `WithFields` as it will pollute the log
*   Avoid multiline logs as this make parsing hard.
*   Avoid using pipes `|` in your log message or fields
*   Instead of `Operation(AWAIT).Warn` use `WarnAwait`
*   Use `Tag` to indentify calls in the same flow (for example on a HTTP Request)
*   All `LogInstance` calls returns a new `LogInstance` with the modified data. It will never change it's parent data making it completely immutable.
