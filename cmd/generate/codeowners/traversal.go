package codeowners

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/jpmcb/gopherlogs"
	"github.com/jpmcb/gopherlogs/pkg/colors"

	"github.com/open-sauced/pizza-cli/pkg/logging"
)

// ProcessOptions are the options for iterating a git reflog and deriving the codeowners
type ProcessOptions struct {
	repo         *git.Repository
	previousDays int
	dirPath      string

	logger gopherlogs.Logger
}

func (po *ProcessOptions) process() (FileStats, error) {
	fs := make(FileStats)

	// Get the HEAD reference
	head, err := po.repo.Head()
	if err != nil {
		return nil, fmt.Errorf("could not get repo head: %w", err)
	}

	now := time.Now()
	previousTime := now.AddDate(0, 0, -po.previousDays)

	// Get the commit history for all files
	commitIter, err := po.repo.Log(&git.LogOptions{
		From:  head.Hash(),
		Since: &previousTime,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get repo log iterator: %w", err)
	}

	defer commitIter.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		po.logger.Style(0, colors.Reset).AnimateProgressWithOptions(
			gopherlogs.AnimatorWithContext(ctx),
			gopherlogs.AnimatorWithMaxLen(80),
			gopherlogs.AnimatorWithMessagef("Iterating commits for repo: %s ", po.dirPath),
		)
	}(ctx)

	err = commitIter.ForEach(func(commit *object.Commit) error {
		// Get the patch for this commit between the head and the parent commit
		patch, err := po.getPatchForCommit(commit)
		if err != nil {
			return fmt.Errorf("could not get patch for commit %s: %w", commit.Hash, err)
		}

		for _, fileStat := range patch.Stats() {
			if !po.isSubPath(po.dirPath, fileStat.Name) {
				// Explicitly ignore paths that do not exist in the repo.
				// This is relevant for old changes and filename changes.
				// Example: this will ignore some/file/path => new/name/path
				// changes that ONLY change the name / path of a file.
				//
				// These are edge cases to revisit in the future.
				return nil
			}

			fs.addStat(&fileStat, commit)
		}

		return nil
	})

	if err != nil {
		cancel()
		return nil, fmt.Errorf("could not process commit iterator: %w", err)
	}

	cancel()
	po.logger.V(logging.LogInfo).Style(0, colors.FgGreen).ReplaceLinef("Finished processing commits for: %s", po.dirPath)
	return fs, nil
}

func (po *ProcessOptions) isSubPath(basePath, relativePath string) bool {
	// Clean the paths to remove any '..' or '.' components
	basePath = filepath.Clean(basePath)
	fullPath := filepath.Join(basePath, relativePath)
	fullPath = filepath.Clean(fullPath)

	// Check if the full path starts with the base path
	return strings.HasPrefix(fullPath, basePath)
}

func (po *ProcessOptions) getPatchForCommit(commit *object.Commit) (*object.Patch, error) {
	// No parents (the initial, first commit). Use a stub of an object tree
	// to simulate "no" parent present in the diff
	if commit.NumParents() == 0 {
		parentTree := &object.Tree{}
		commitTree, err := commit.Tree()
		if err != nil {
			return nil, fmt.Errorf("could not get commit tree for commit %s: %w", commit.Hash, err)
		}

		return parentTree.Patch(commitTree)
	}

	commitTree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("could not get commit tree for commit %s: %w", commit.Hash, err)
	}

	parentCommit, err := commit.Parents().Next()
	if err != nil && err != storer.ErrStop {
		return nil, fmt.Errorf("could not get parent commit to commit %s: %w", commit.Hash, err)
	}

	parentTree, err := parentCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("could not get parent commit tree for parent commit %s: %w", parentCommit.Hash, err)
	}

	return parentTree.Patch(commitTree)
}
