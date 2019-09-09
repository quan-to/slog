package main

import (
    "github.com/quan-to/slog"
    "time"
)

var log = slog.Scope("MAIN")


func call0(i slog.Instance, arg0 string) {
    l := i.SubScope("Call0").WithFields(map[string]interface{}{
        "arg0": arg0,
    })
    l.Await("Doing some work")
    time.Sleep(time.Second)
    l.Done("Finished some work")
    l.Note("Not sure what I'm doing...")
    l.Info("Calling Call1")
    call1(l, "call1arg")
    l.Done("Exiting")
}

func call1(i slog.Instance, huebr string) {
    l := i.SubScope("Call1").WithFields(map[string]interface{}{
        "huebr": huebr,
    })
    l.Info("Calling Call2")
    call2(l, "abcde")
    l.Warn("Call 1 finished")
}

func call2(i slog.Instance, pop string) {
    l := i.SubScope("Call2").WithFields(map[string]interface{}{
        "pop": pop,
    })

    l.IO("Doing some IO")
    l.Error("I'm useless. Please fix-me")
}

func main() {
    slog.SetShowLines(true)
    slog.SetScopeLength(40) // Expand Scope pad length

    log = log.Tag("REQ001") // Tag current as REQ001

    log.Info("Starting program")
    t := log.Operation(slog.AWAIT)

    t.Operation(slog.AWAIT).Info("HUE TESTE STACK O")
    t.Info("HUE TESTE STACK")
    
    call0(log, "MyArg0")

    call1(log, "Call1Arg")
    call2(log, "Call2Arg")
}