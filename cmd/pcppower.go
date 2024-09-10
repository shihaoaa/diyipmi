package main

import (
        "flag"
        "fmt"
        "os"
        "sync"
        "time"

        "github.com/stianeikeland/go-rpio/v4"
)

var lock sync.Mutex

// 封装操作树莓派rpio功能
func switchPin(pin *rpio.Pin, pTime int64, dTime time.Duration) {
        pin.Output()
        pin.High()
        time.Sleep(time.Duration(pTime) * dTime)
        pin.Low()
}

func useage() {
        fmt.Printf("Usage: pcpower <action>\n\nactions:\n    start:     power on the pi-bound pc.\n    stop:    power off the pi-bound pc.\n    restart:        reboot the the pi-bound pc.\n    status:       query the the pi-bound pc power status.\n")
        os.Exit(125)
}

func main() {
        flag.Usage = useage
        flag.Parse()
        actions := flag.Args()
        if len(actions) != 1 {
                useage()
        }
        action := actions[0]
        var rightUse bool = false
        actionsStyle := []string{"start", "stop", "restart", "status"}
        for _, actionStle := range actionsStyle {
                if action == actionStle {
                        rightUse = true
                }
        }
        if !rightUse {
                fmt.Printf("#: Unknown action arg:--> %s\n\n", action)
                useage()
        }
        err := rpio.Open()
        if err != nil {
                fmt.Println(err)
        }
        defer rpio.Close()
        if action == "status" {
                lock.Lock()
                defer lock.Unlock()
                pin := rpio.Pin(23)
                pin.Input()
                powerStatus := 1
                for i := 0; i < 2; i++ {
                        if pin.Read() == 0 {
                                powerStatus = 0
                                break
                        }
                        time.Sleep(time.Millisecond * 400)
                }
                fmt.Println(powerStatus)
                return
        }
        pin := rpio.Pin(18)
        switch action {
        case "start":
                switchPin(&pin, 400, time.Millisecond)
                fmt.Printf("action: %s complete\n", action)
        case "stop":
                switchPin(&pin, 6, time.Second)
                fmt.Printf("action: %s complete\n", action)
        case "restart":
                switchPin(&pin, 6, time.Second)
                time.Sleep(1 * time.Second)
                switchPin(&pin, 300, time.Millisecond)
                fmt.Printf("action: %s complete\n", action)
        }
}


