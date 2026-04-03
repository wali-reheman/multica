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

// WatcherService watches project directories for file changes.
type WatcherService struct {
	bus     *events.Bus
	mu      sync.Mutex
	watches map[string]*projectWatch
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
	".git", "node_modules", ".next", "__pycache__",
	".multica-local", "vendor", ".idea", ".vscode",
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
	if _, ok := ws.watches[projectID]; ok {
		return nil
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	pw := &projectWatch{
		watcher:     watcher,
		projectID:   projectID,
		workspaceID: workspaceID,
		localPath:   localPath,
		stopCh:      make(chan struct{}),
	}
	if err := watcher.Add(localPath); err != nil {
		watcher.Close()
		return err
	}
	ws.watches[projectID] = pw
	go ws.runWatch(pw)
	slog.Info("started watching project", "project_id", projectID, "path", localPath)
	return nil
}

// Unwatch stops watching a project directory.
func (ws *WatcherService) Unwatch(projectID string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	if pw, ok := ws.watches[projectID]; ok {
		close(pw.stopCh)
		pw.watcher.Close()
		delete(ws.watches, projectID)
	}
}

// Close stops all watches.
func (ws *WatcherService) Close() {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	for id, pw := range ws.watches {
		close(pw.stopCh)
		pw.watcher.Close()
		delete(ws.watches, id)
	}
}

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
		case event, ok := <-pw.watcher.Events:
			if !ok {
				return
			}
			if shouldIgnore(pw.localPath, event.Name) {
				continue
			}
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
		case err, ok := <-pw.watcher.Errors:
			if !ok {
				return
			}
			slog.Error("watcher error", "project_id", pw.projectID, "error", err)
		}
	}
}

func shouldIgnore(basePath, changedPath string) bool {
	rel, err := filepath.Rel(basePath, changedPath)
	if err != nil {
		return false
	}
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
		for _, pattern := range DefaultIgnorePatterns {
			if part == pattern {
				return true
			}
		}
	}
	return false
}
