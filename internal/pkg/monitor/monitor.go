package monitor

import (
	"fmt"

	"github.com/emmercm/photos/internal/pkg/model"
	log "github.com/sirupsen/logrus"
)

// Monitor represents different processes that can monitor for file changes
type Monitor interface {
	WatchDirectory(string) error
	UnwatchDirectory(string) error
}

// Collection represents a set of Monitors
type Collection struct {
	eventChannel EventChannel
	monitors     []Monitor
}

// NewCollection creates a new Collection of Monitors
func NewCollection() *Collection {
	eventChannel := make(EventChannel)

	return &Collection{
		eventChannel: eventChannel,
		monitors: []Monitor{
			Monitor(NewScanner(eventChannel)),
			Monitor(NewWatcher(eventChannel)),
		},
	}
}

// WatchDirectory calls WatchDirectory() on all Monitors
func (c *Collection) WatchDirectory(directory string) error {
	for _, m := range c.monitors {
		if err := m.WatchDirectory(directory); err != nil {
			return err
		}
	}

	return nil
}

// UnwatchDirectory calls UnwatchDirectory() on all Monitors
func (c *Collection) UnwatchDirectory(directory string) error {
	for _, m := range c.monitors {
		if err := m.UnwatchDirectory(directory); err != nil {
			return err
		}
	}

	return nil
}

// Start listens for events from Monitors
func (c *Collection) Start() {
	for {
		event := <-c.eventChannel
		if err := processEvent(event); err != nil {
			log.Error(err)
		}
	}
}

func processEvent(event Event) error {
	fmt.Println(event)

	f, err := model.FileFromPath(event.Filename)
	if err != nil {
		return err
	}

	if event.Type == Delete {
		return f.Delete()
	}

	hasChanged, err := f.HasChanged()
	if err != nil {
		return err
	}

	if hasChanged {
		if err := f.Update(); err != nil {
			return err
		}

		if err := f.Save(); err != nil {
			return err
		}
	}

	a, err := f.GetDirnameAlbum()
	if err != nil {
		return err
	}

	if err := a.AddFile(f); err != nil {
		return err
	}

	return nil
}
