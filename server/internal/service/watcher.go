package service

import (
	"log/slog"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/multica-ai/multica/server/internal/events"
	"github.com/multica-ai/multica/server/pkg/protocol"
)

<<<<<<< HEAD
// WatcherService watches project directories for file changes and
// broadcasts events via the event bus.
type WatcherService struct {
	bus     *events.Bus
	mu      sync.Mutex
	watches map[string]*projectWatch // projectID -> watch
=======
// WatcherService watches project directories for file changes.
type WatcherService struct {
	bus     *events.Bus
	mu      sync.Mutex
	watches map[string]*projectWatch
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
}

type projectWatch struct {
	watcher     *fsnotify.Watcher
	projectID   string
	workspaceID string
	localPath   string
	stopCh      chan struct{}
}

// DefaultIgnorePatterns are directories that should not trigger change events.
var DefaultIgnorePatterns = []string{
<<<<<<< HEAD
	".git",
	"node_modules",
	".next",
	"__pycache__",
	".multica-local",
	"vendor",
	".idea",
	".vscode",
=======
	".git", "node_modules", ".next", "__pycache__",
	".multica-local", "vendor", ".idea", ".vscode",
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
}

// NewWatcherService creates a new WatcherService.
func NewWatcherService(bus *events.Bus) *WatcherService {
	return &WatcherService{
		bus:     bus,
		watches: make(map[string]*projectWatch),
	}
}

// Watch starts watching a project directory for changes.
func (ws *WatcherService) Watch(projectID, workspaceID, localPath string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
<<<<<<< HEAD

	// Already watching
	if _, ok := ws.watches[projectID]; ok {
		return nil
	}

=======
	if _, ok := ws.watches[projectID]; ok {
		return nil
	}
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	pw := &projectWatch{
		watcher:     watcher,
		projectID:   projectID,
		workspaceID: workspaceID,
		localPath:   localPath,
		stopCh:      make(chan struct{}),
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if err := watcher.Add(localPath); err != nil {
		watcher.Close()
		return err
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	ws.watches[projectID] = pw
	go ws.runWatch(pw)
	slog.Info("started watching project", "project_id", projectID, "path", localPath)
	return nil
}

// Unwatch stops watching a project directory.
func (ws *WatcherService) Unwatch(projectID string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if pw, ok := ws.watches[projectID]; ok {
		close(pw.stopCh)
		pw.watcher.Close()
		delete(ws.watches, projectID)
<<<<<<< HEAD
		slog.Info("stopped watching project", "project_id", projectID)
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	}
}

// Close stops all watches.
func (ws *WatcherService) Close() {
	ws.mu.Lock()
	defer ws.mu.Unlock()
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	for id, pw := range ws.watches {
		close(pw.stopCh)
		pw.watcher.Close()
		delete(ws.watches, id)
	}
}

<<<<<<< HEAD
// runWatch processes filesystem events with debouncing.
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (ws *WatcherService) runWatch(pw *projectWatch) {
	const debounceWindow = 100 * time.Millisecond
	var timer *time.Timer

	for {
		select {
		case <-pw.stopCh:
			if timer != nil {
				timer.Stop()
			}
			return
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
		case event, ok := <-pw.watcher.Events:
			if !ok {
				return
			}
<<<<<<< HEAD

			// Ignore changes in excluded directories
			if shouldIgnore(pw.localPath, event.Name) {
				continue
			}

			// Debounce: reset timer on each event
=======
			if shouldIgnore(pw.localPath, event.Name) {
				continue
			}
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(debounceWindow, func() {
				ws.bus.Publish(events.Event{
					Type:        protocol.EventProjectFilesChanged,
					WorkspaceID: pw.workspaceID,
					Payload: map[string]any{
						"project_id": pw.projectID,
						"path":       pw.localPath,
					},
				})
			})
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
		case err, ok := <-pw.watcher.Errors:
			if !ok {
				return
			}
			slog.Error("watcher error", "project_id", pw.projectID, "error", err)
		}
	}
}

<<<<<<< HEAD
// shouldIgnore returns true if the changed file is in an ignored directory.
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func shouldIgnore(basePath, changedPath string) bool {
	rel, err := filepath.Rel(basePath, changedPath)
	if err != nil {
		return false
	}
<<<<<<< HEAD

	parts := strings.Split(rel, string(filepath.Separator))
	for _, part := range parts {
=======
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
		for _, pattern := range DefaultIgnorePatterns {
			if part == pattern {
				return true
			}
		}
	}
	return false
}
