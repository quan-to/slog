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
