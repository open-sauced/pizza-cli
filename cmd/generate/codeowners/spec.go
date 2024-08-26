package codeowners

import (
	"fmt"
	"sort"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// FileStats is a mapping of filenames to author stats.
// Example: { "path/to/file": { Author stats }}
type FileStats map[string]AuthorStats

func (fs FileStats) addStat(filestat *object.FileStat, commit *object.Commit) {
	author := fmt.Sprintf("%s <%s>", commit.Author.Name, commit.Author.Email)
	filename := filestat.Name

	if _, ok := fs[filename]; !ok {
		fs[filename] = make(AuthorStats)
	}

	if _, ok := fs[filename][author]; !ok {
		fs[filename][author] = &CodeownerStat{
			Name:  commit.Author.Name,
			Email: commit.Author.Email,
		}
	}

	fs[filename][author].Lines += filestat.Addition + filestat.Deletion
}

// AuthorStats is a mapping of author name email combinations to codeowner stats.
// Example: { "First Last name@domain.com": { Codeowner stat }}
type AuthorStats map[string]*CodeownerStat

// CodeownerStat is the base struct of name, email, lines changed,
// and the configured GitHub alias for a given codeowner. This is derived from
// git reflog commits and diffs.
type CodeownerStat struct {
	Name        string
	Email       string
	Lines       int
	GitHubAlias string
}

// AuthorStatSlice is a slice of codeowner stats. This is a utility type that makes
// turning a mapping of author stats to slices easy.
type AuthorStatSlice []*CodeownerStat

func (as AuthorStats) ToSortedSlice() AuthorStatSlice {
	slice := make(AuthorStatSlice, 0, len(as))

	for _, stat := range as {
		slice = append(slice, stat)
	}

	sort.Slice(slice, func(i, j int) bool {
		// sort the author stats by descending number of lines
		return slice[i].Lines > slice[j].Lines
	})

	return slice
}
