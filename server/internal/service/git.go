package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

// CommitInfo represents a single git commit.
type CommitInfo struct {
<<<<<<< HEAD
	Hash       string    `json:"hash"`
	ShortHash  string    `json:"short_hash"`
	Message    string    `json:"message"`
	Author     string    `json:"author"`
	AuthorEmail string   `json:"author_email"`
	Date       time.Time `json:"date"`
	FilesChanged int     `json:"files_changed"`
=======
	Hash         string    `json:"hash"`
	ShortHash    string    `json:"short_hash"`
	Message      string    `json:"message"`
	Author       string    `json:"author"`
	AuthorEmail  string    `json:"author_email"`
	Date         time.Time `json:"date"`
	FilesChanged int       `json:"files_changed"`
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
}

// BranchInfo represents a git branch.
type BranchInfo struct {
	Name     string `json:"name"`
	IsHead   bool   `json:"is_head"`
	IsRemote bool   `json:"is_remote"`
	Hash     string `json:"hash"`
}

// FileStatusEntry represents the status of a single file.
type FileStatusEntry struct {
	Path     string `json:"path"`
	Staging  string `json:"staging"`
	Worktree string `json:"worktree"`
}

// GitStatus represents the working tree status.
type GitStatus struct {
	Modified  []FileStatusEntry `json:"modified"`
	Added     []FileStatusEntry `json:"added"`
	Deleted   []FileStatusEntry `json:"deleted"`
	Untracked []FileStatusEntry `json:"untracked"`
}

// DiffEntry represents a file-level diff.
type DiffEntry struct {
	Path    string `json:"path"`
	OldPath string `json:"old_path,omitempty"`
	Change  string `json:"change"` // "add", "modify", "delete", "rename"
	Patch   string `json:"patch,omitempty"`
}

// CommitDetail has full commit details including diffs.
type CommitDetail struct {
	CommitInfo
	Diffs []DiffEntry `json:"diffs"`
}

// GitService wraps go-git operations for local project repositories.
type GitService struct{}

// NewGitService creates a new GitService.
func NewGitService() *GitService {
	return &GitService{}
}

<<<<<<< HEAD
// openRepo opens an existing git repository at the given path.
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (s *GitService) openRepo(localPath string) (*git.Repository, error) {
	repo, err := git.PlainOpen(localPath)
	if err != nil {
		return nil, fmt.Errorf("open repo at %s: %w", localPath, err)
	}
	return repo, nil
}

// Init initializes a git repository at the given path if not already present.
<<<<<<< HEAD
// Returns true if a new repo was initialized, false if it already existed.
func (s *GitService) Init(localPath string) (bool, error) {
	_, err := git.PlainOpen(localPath)
	if err == nil {
		return false, nil // already a repo
	}

=======
func (s *GitService) Init(localPath string) (bool, error) {
	_, err := git.PlainOpen(localPath)
	if err == nil {
		return false, nil
	}
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	_, err = git.PlainInit(localPath, false)
	if err != nil {
		return false, fmt.Errorf("init repo: %w", err)
	}
<<<<<<< HEAD

	// Add .multica-local to .gitignore
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	s.ensureGitignoreEntry(localPath, ".multica-local")
	return true, nil
}

// IsRepo checks whether the given path is a git repository.
func (s *GitService) IsRepo(localPath string) bool {
	_, err := git.PlainOpen(localPath)
	return err == nil
}

// Status returns the working tree status.
func (s *GitService) Status(localPath string) (*GitStatus, error) {
	repo, err := s.openRepo(localPath)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	wt, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("get worktree: %w", err)
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	status, err := wt.Status()
	if err != nil {
		return nil, fmt.Errorf("get status: %w", err)
	}

	gs := &GitStatus{
		Modified:  []FileStatusEntry{},
		Added:     []FileStatusEntry{},
		Deleted:   []FileStatusEntry{},
		Untracked: []FileStatusEntry{},
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	for path, s := range status {
		entry := FileStatusEntry{
			Path:     path,
			Staging:  statusCodeToString(s.Staging),
			Worktree: statusCodeToString(s.Worktree),
		}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
		switch {
		case s.Worktree == git.Untracked:
			gs.Untracked = append(gs.Untracked, entry)
		case s.Worktree == git.Deleted || s.Staging == git.Deleted:
			gs.Deleted = append(gs.Deleted, entry)
		case s.Staging == git.Added:
			gs.Added = append(gs.Added, entry)
		case s.Worktree == git.Modified || s.Staging == git.Modified:
			gs.Modified = append(gs.Modified, entry)
		default:
<<<<<<< HEAD
			// Catch-all for renamed, copied, etc.
			gs.Modified = append(gs.Modified, entry)
		}
	}

=======
			gs.Modified = append(gs.Modified, entry)
		}
	}
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	return gs, nil
}

// Add stages files for commit.
func (s *GitService) Add(localPath string, files []string) error {
	repo, err := s.openRepo(localPath)
	if err != nil {
		return err
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("get worktree: %w", err)
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	for _, f := range files {
		if _, err := wt.Add(f); err != nil {
			return fmt.Errorf("stage %s: %w", f, err)
		}
	}
	return nil
}

// Commit creates a commit with the given message.
func (s *GitService) Commit(localPath, message, authorName, authorEmail string) (string, error) {
	repo, err := s.openRepo(localPath)
	if err != nil {
		return "", err
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	wt, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("get worktree: %w", err)
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	hash, err := wt.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", fmt.Errorf("commit: %w", err)
	}
	return hash.String(), nil
}

// Log returns paginated commit history.
func (s *GitService) Log(localPath string, limit, offset int) ([]CommitInfo, error) {
	repo, err := s.openRepo(localPath)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD

	ref, err := repo.Head()
	if err != nil {
		// Empty repo with no commits
=======
	ref, err := repo.Head()
	if err != nil {
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
		if err == plumbing.ErrReferenceNotFound {
			return []CommitInfo{}, nil
		}
		return nil, fmt.Errorf("get HEAD: %w", err)
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	iter, err := repo.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, fmt.Errorf("log: %w", err)
	}
	defer iter.Close()

	var commits []CommitInfo
	i := 0
	err = iter.ForEach(func(c *object.Commit) error {
		if i < offset {
			i++
			return nil
		}
		if len(commits) >= limit {
			return io.EOF
		}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
		stats, _ := c.Stats()
		commits = append(commits, CommitInfo{
			Hash:         c.Hash.String(),
			ShortHash:    c.Hash.String()[:7],
			Message:      strings.TrimSpace(c.Message),
			Author:       c.Author.Name,
			AuthorEmail:  c.Author.Email,
			Date:         c.Author.When,
			FilesChanged: len(stats),
		})
		i++
		return nil
	})
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("iterate log: %w", err)
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	return commits, nil
}

// Show returns full commit details including diffs.
func (s *GitService) Show(localPath, commitHash string) (*CommitDetail, error) {
	repo, err := s.openRepo(localPath)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	hash := plumbing.NewHash(commitHash)
	commit, err := repo.CommitObject(hash)
	if err != nil {
		return nil, fmt.Errorf("get commit %s: %w", commitHash, err)
	}

	detail := &CommitDetail{
		CommitInfo: CommitInfo{
			Hash:        commit.Hash.String(),
			ShortHash:   commit.Hash.String()[:7],
			Message:     strings.TrimSpace(commit.Message),
			Author:      commit.Author.Name,
			AuthorEmail: commit.Author.Email,
			Date:        commit.Author.When,
		},
		Diffs: []DiffEntry{},
	}

<<<<<<< HEAD
	// Get parent for diff (nil parent = initial commit)
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	var parentTree *object.Tree
	if commit.NumParents() > 0 {
		parent, err := commit.Parent(0)
		if err == nil {
			parentTree, _ = parent.Tree()
		}
	}

	commitTree, err := commit.Tree()
	if err != nil {
<<<<<<< HEAD
		return detail, nil // return without diffs
=======
		return detail, nil
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	}

	changes, err := object.DiffTree(parentTree, commitTree)
	if err != nil {
		return detail, nil
	}

	detail.CommitInfo.FilesChanged = len(changes)
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	for _, change := range changes {
		de := DiffEntry{}
		action, err := change.Action()
		if err != nil {
			continue
		}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
		switch action {
		case merkletrie.Insert:
			de.Change = "add"
			de.Path = change.To.Name
		case merkletrie.Delete:
			de.Change = "delete"
			de.Path = change.From.Name
		case merkletrie.Modify:
			de.Change = "modify"
			de.Path = change.To.Name
			if change.From.Name != change.To.Name {
				de.Change = "rename"
				de.OldPath = change.From.Name
			}
		}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
		patch, err := change.Patch()
		if err == nil && patch != nil {
			de.Patch = patch.String()
		}
<<<<<<< HEAD

		detail.Diffs = append(detail.Diffs, de)
	}

	return detail, nil
}

// Diff returns working tree diff (unstaged changes).
=======
		detail.Diffs = append(detail.Diffs, de)
	}
	return detail, nil
}

// Diff returns working tree diff.
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (s *GitService) Diff(localPath string) ([]DiffEntry, error) {
	repo, err := s.openRepo(localPath)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	wt, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("get worktree: %w", err)
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	status, err := wt.Status()
	if err != nil {
		return nil, fmt.Errorf("get status: %w", err)
	}

	var diffs []DiffEntry
	for path, s := range status {
		if s.Worktree == git.Unmodified && s.Staging == git.Unmodified {
			continue
		}
		change := "modify"
		switch {
		case s.Worktree == git.Untracked:
			change = "add"
		case s.Worktree == git.Deleted || s.Staging == git.Deleted:
			change = "delete"
		case s.Staging == git.Added:
			change = "add"
		}
<<<<<<< HEAD
		diffs = append(diffs, DiffEntry{
			Path:   path,
			Change: change,
		})
	}

	return diffs, nil
}

// Branches returns all branches (local and remote).
=======
		diffs = append(diffs, DiffEntry{Path: path, Change: change})
	}
	return diffs, nil
}

// Branches returns all branches.
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
func (s *GitService) Branches(localPath string) ([]BranchInfo, error) {
	repo, err := s.openRepo(localPath)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	head, _ := repo.Head()
	headHash := ""
	if head != nil {
		headHash = head.Hash().String()
	}

	var branches []BranchInfo
<<<<<<< HEAD

	// Local branches
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	refs, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("list branches: %w", err)
	}
	refs.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, BranchInfo{
			Name:   ref.Name().Short(),
			IsHead: ref.Hash().String() == headHash,
			Hash:   ref.Hash().String(),
		})
		return nil
	})

<<<<<<< HEAD
	// Remote branches
=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	remoteRefs, err := repo.References()
	if err == nil {
		remoteRefs.ForEach(func(ref *plumbing.Reference) error {
			if ref.Name().IsRemote() {
				branches = append(branches, BranchInfo{
					Name:     ref.Name().Short(),
					IsRemote: true,
					Hash:     ref.Hash().String(),
				})
			}
			return nil
		})
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	return branches, nil
}

// Checkout switches to the given branch.
func (s *GitService) Checkout(localPath, branch string) error {
	repo, err := s.openRepo(localPath)
	if err != nil {
		return err
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	wt, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("get worktree: %w", err)
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	return wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
	})
}

// CreateBranch creates a new branch from HEAD.
func (s *GitService) CreateBranch(localPath, name string) error {
	repo, err := s.openRepo(localPath)
	if err != nil {
		return err
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("get HEAD: %w", err)
	}
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	ref := plumbing.NewHashReference(plumbing.NewBranchReferenceName(name), head.Hash())
	return repo.Storer.SetReference(ref)
}

<<<<<<< HEAD
// ensureGitignoreEntry appends an entry to .gitignore if not already present.
func (s *GitService) ensureGitignoreEntry(localPath, entry string) {
	gitignorePath := filepath.Join(localPath, ".gitignore")

	content, err := os.ReadFile(gitignorePath)
	if err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == entry {
				return // already present
			}
		}
	}

=======
func (s *GitService) ensureGitignoreEntry(localPath, entry string) {
	gitignorePath := filepath.Join(localPath, ".gitignore")
	content, err := os.ReadFile(gitignorePath)
	if err == nil {
		for _, line := range strings.Split(string(content), "\n") {
			if strings.TrimSpace(line) == entry {
				return
			}
		}
	}
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
<<<<<<< HEAD

=======
>>>>>>> aef083616f315280ce283baf1ae5fd21992cd609
	if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
		f.WriteString("\n")
	}
	f.WriteString(entry + "\n")
}

func statusCodeToString(c git.StatusCode) string {
	switch c {
	case git.Unmodified:
		return "unmodified"
	case git.Modified:
		return "modified"
	case git.Added:
		return "added"
	case git.Deleted:
		return "deleted"
	case git.Renamed:
		return "renamed"
	case git.Copied:
		return "copied"
	case git.Untracked:
		return "untracked"
	default:
		return "unknown"
	}
}
