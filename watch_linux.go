package usb

import (
	"fmt"
	"os"
	"time"
)

type Event struct {
	Type   EventType
	Device string
}

type EventType int

const (
	Connect EventType = iota
	Disconnect
)

func (e EventType) String() string {
	switch e {
	case Connect:
		return "Connect"
	case Disconnect:
		return "Disconnect"
	default:
		return fmt.Sprintf("EventType(%d)", int(e))
	}
}

func Watch() Watcher {
	w := Watcher{
		events:           make(chan Event, 16),
		stop:             make(chan bool, 1),
		devicesAvailable: make(map[string]bool),
	}
	w.Events = w.events
	for n := 'a'; n <= 'j'; n++ {
		w.devicesAvailable[fmt.Sprintf("/dev/sd%s1", string(n))] = false
	}
	go func() {
		for {
			select {
			case <-w.stop:
				return
			default:
				for path, avail := range w.devicesAvailable {
					_, err := os.Stat(path)
					if err == nil && !avail {
						w.devicesAvailable[path] = true
						w.events <- Event{
							Type:   Connect,
							Device: path,
						}
					} else if err != nil && avail {
						w.devicesAvailable[path] = false
						w.events <- Event{
							Type:   Disconnect,
							Device: path,
						}
					}
				}
				time.Sleep(time.Second)
			}
		}
	}()
	return w
}

type Watcher struct {
	Events           <-chan Event
	events           chan Event
	stop             chan bool
	devicesAvailable map[string]bool
}

func (w Watcher) Stop() {
	w.stop <- true
}
