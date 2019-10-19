package monitor

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/rjeczalik/notify"
	log "github.com/sirupsen/logrus"
)

// ChannelBufferSize is the buffer size for the OS file event channel
const ChannelBufferSize = 1024

var (
	watcherOnce sync.Once
	watcher     *Watcher
)

// Watcher represents an OS file event watcher
type Watcher struct {
	eventInfo    chan notify.EventInfo
	eventChannel EventChannel
	directories  []string
}

// NewWatcher returns a new Watcher
func NewWatcher(eventChannel EventChannel) *Watcher {
	watcherOnce.Do(func() {
		watcher = &Watcher{
			eventInfo:    make(chan notify.EventInfo, ChannelBufferSize),
			eventChannel: eventChannel,
		}

		go watcher.watchLoop()
	})

	return watcher
}

func (w *Watcher) watchLoop() {
	for {
		eventInfo := <-w.eventInfo

		if !matchesGlob(eventInfo.Path()) {
			continue
		}

		event := Event{
			Filename: eventInfo.Path(),
			Type:     eventType(eventInfo),
		}

		// log.Debug(event)

		w.eventChannel <- event
	}
}

func eventType(eventInfo notify.EventInfo) EventType {
	switch eventInfo.Event() {
	case notify.Create:
		return Create
	case notify.Remove:
		return Delete
	case notify.Rename:
		if _, err := os.Stat(eventInfo.Path()); os.IsNotExist(err) {
			return Delete
		}
		return Create
	default:
		return Update
	}
}

// Close stops the Watcher entirely
func (w *Watcher) Close() {
	notify.Stop(w.eventInfo)
	w.directories = []string{}
}

// WatchDirectory listens for OS file events
func (w *Watcher) WatchDirectory(directory string) error {
	directoryAbs, err := filepath.Abs(directory)
	if err != nil {
		return err
	}

	entry := log.WithFields(log.Fields{
		"event":     "watch directory",
		"directory": directoryAbs,
	})

	path := fmt.Sprintf("%v/...", directoryAbs)
	if err := notify.Watch(path, w.eventInfo, notify.All); err != nil {
		// notify.Stop(w.eventInfo)
		return err
	}

	w.directories = append(w.directories, directoryAbs)
	entry = entry.WithFields(log.Fields{
		"directories": w.directories,
	})

	entry.Info()
	return nil
}

// UnwatchDirectory stops watching a directory for events
func (w *Watcher) UnwatchDirectory(directory string) error {
	directoryAbs, err := filepath.Abs(directory)
	if err != nil {
		return err
	}

	entry := log.WithFields(log.Fields{
		"event":     "unwatch directory",
		"directory": directoryAbs,
	})

	for i, v := range w.directories {
		if v == directoryAbs {
			directories := append(w.directories[:i], w.directories[i+1:]...)

			w.Close()

			for _, q := range directories {
				if err := w.WatchDirectory(q); err != nil {
					return err
				}
			}

			entry = entry.WithFields(log.Fields{
				"directories": w.directories,
			})

			break
		}
	}

	entry.Info()
	return nil
}
