package models

import (
	"fmt"
	"path"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/sync-file-server/backend"
)

type SyncFileOption struct {
	Branch
	BranchSHA string   `json:"branch_sha" required:"true"`
	Files     []string `json:"files" required:"true"`
}

func (s SyncFileOption) Create() error {
	files := map[string][]string{}
	for _, f := range s.Files {
		name := parseFileName(f)
		files[name] = append(files[name], f)
	}

	for _, item := range files {
		syncFile(s.Branch, s.BranchSHA, item)
	}
	return nil
}

type fetchFileResult struct {
	File
	err error
}

func fetchFileError(f string, err error) fetchFileResult {
	r := fetchFileResult{err: err}
	r.Path = f
	return r
}

func syncFile(branch Branch, branchSHA string, files []string) error {
	n := len(files)
	if n == 0 {
		return nil
	}

	c := backend.GetClient()

	getFileContent := func(f string) (fetchFileResult, error) {
		sha, content, err := c.GetFileConent(branch, f)
		if err != nil {
			return fetchFileResult{}, err
		}

		return fetchFileResult{
			File: File{
				RepoFile: RepoFile{
					Path: f,
					SHA:  sha,
				},
				Content: content,
			},
		}, nil
	}

	log := logEntryForBranch(branch, branchSHA)

	result := make([]File, 0, n)
	for _, f := range files {
		if r, err := getFileContent(f); err != nil {
			log.WithField("file", r.Path).WithError(r.err).Error("sync file")
		} else {
			result = append(result, r.File)
		}
	}

	if n := len(result); n > 0 {
		if err := c.SaveFiles(branch, branchSHA, result); err != nil {
			fs := make([]string, n)
			for i := range result {
				fs[i] = result[i].Path
			}

			return fmt.Errorf(
				"error to save files: %s, err: %v",
				strings.Join(fs, "; "), err,
			)
		}
	}

	return nil
}

func parseFileName(s string) string {
	return path.Base(s)
}

func logEntryForBranch(b Branch, sha string) *logrus.Entry {
	return logrus.WithFields(
		logrus.Fields{
			"org":        b.Org,
			"repo":       b.Repo,
			"branch":     b.Branch,
			"branch sha": sha,
		},
	)
}
