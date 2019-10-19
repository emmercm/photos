package monitor

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	scannerOnce sync.Once
	scanner     *Scanner
)

// Scanner represents a brute-force file scanner
type Scanner struct {
	eventChannel EventChannel
	directories  map[string]context.CancelFunc
}

// NewScanner returns a new Scanner
func NewScanner(eventChannel EventChannel) *Scanner {
	scannerOnce.Do(func() {
		scanner = &Scanner{
			eventChannel: eventChannel,
			directories:  make(map[string]context.CancelFunc),
		}
	})

	return scanner
}

func (s *Scanner) checkDirectory(ctx context.Context, directory string, sleep time.Duration) error {
	directoryAbs, err := filepath.Abs(directory)
	if err != nil {
		return err
	}

	err = filepath.Walk(directoryAbs, func(path string, info os.FileInfo, err error) error {
		select {
		case <-ctx.Done():
			return io.EOF
		default:
		}

		if info != nil && info.IsDir() {
			return nil
		}

		if !matchesGlob(path) {
			return nil
		}

		event := Event{
			Filename: path,
			Type:     Check,
		}

		// log.Debug(event)

		s.eventChannel <- event

		time.Sleep(sleep)

		return nil
	})
	if err != nil && err != io.EOF {
		return err
	}

	return nil
}

// WatchDirectory performs an initial directory scan and then monitors for changes
func (s *Scanner) WatchDirectory(directory string) error {
	directoryAbs, err := filepath.Abs(directory)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.directories[directoryAbs] = cancel

	go func() {
		if err := s.checkDirectory(ctx, directoryAbs, 0); err != nil {
			log.Error(err)
		}
	}()

	// go func() {
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			return
	// 		default:
	// 		}
	//
	// 		if err := s.checkDirectory(ctx, directoryAbs, time.Second); err != nil {
	// 			log.Error(err)
	// 		}
	// 	}
	// }()

	return nil
}

// UnwatchDirectory stops watching a directory for changes
func (s *Scanner) UnwatchDirectory(directory string) error {
	directoryAbs, err := filepath.Abs(directory)
	if err != nil {
		return err
	}

	cancel, ok := s.directories[directoryAbs]
	if !ok {
		return nil
	}

	cancel()

	delete(s.directories, directoryAbs)

	return nil
}
