package models

import (
	"encoding/json"

	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"

	"github.com/opensourceways/sync-file-server/backend"
)

type File = backend.File
type Branch = backend.Branch
type RepoFile = backend.RepoFile

type SyncRepoFileOption struct {
	Branch
	BranchSHA string   `json:"branch_sha" required:"true"`
	FileNames []string `json:"file_names" required:"true"`
}

type msgTask struct {
	Org       string   `json:"org"`
	Repo      string   `json:"repo"`
	Branch    string   `json:"branch"`
	BranchSHA string   `json:"branchSHA"`
	Files     []string `json:"files"`
}

func SyncFromMQ(topic, component string) (mq.Subscriber, error) {
	return kafka.Subscribe(topic, handle, func(options *mq.SubscribeOptions) {
		options.Queue = component
	})
}

func handle(e mq.Event) error {
	body := e.Message().Body
	mt := new(msgTask)
	if err := json.Unmarshal(body, mt); err != nil {
		return err
	}

	opt := SyncRepoFileOption{
		Branch: Branch{
			Org:    mt.Org,
			Repo:   mt.Repo,
			Branch: mt.Branch,
		},
		BranchSHA: mt.BranchSHA,
		FileNames: mt.Files,
	}

	if err := opt.Create(); err != nil {
		return err
	}

	return nil
}

func (s SyncRepoFileOption) Create() error {
	c := backend.GetClient()
	log := logEntryForBranch(s.Branch, s.BranchSHA)

	allFiles, err := c.ListAllFilesOfRepo(s.Branch)
	if err != nil {
		log.WithError(err).Error()
		return err
	}

	files := make(map[string][]RepoFile)
	for _, f := range s.FileNames {
		files[f] = []RepoFile{}
	}

	for _, f := range allFiles {
		name := parseFileName(f.Path)
		if v, ok := files[name]; ok {
			files[name] = append(v, f)
		}
	}

	logError := func(f, msg string, err error) {
		log.WithField("file name", f).WithError(err).Error(msg)
	}

	for fileName, item := range files {
		if len(item) == 0 {
			logError(fileName, "there is not corresponding file", nil)
			continue
		}

		todo, err := s.filterFile(fileName, item)
		if err != nil {
			logError(fileName, "filter file", err)
			return err
		}

		if err := syncFile(s.Branch, s.BranchSHA, todo); err != nil {
			logError(fileName, "sync file", err)
		}
		log.Info("sync file success: ", fileName)
	}

	return nil
}

func (s SyncRepoFileOption) filterFile(fileName string, files []RepoFile) ([]string, error) {
	if len(files) == 0 {
		return nil, nil
	}

	c := backend.GetClient()

	summary, err := c.GetFileSummary(s.Branch, fileName)
	if err != nil {
		return nil, err
	}

	m := map[string]string{}
	for _, item := range summary {
		m[item.Path] = item.SHA
	}

	todo := make([]string, 0, len(files))
	for _, item := range files {
		if item.SHA != m[item.Path] {
			todo = append(todo, item.Path)
		}
	}
	return todo, nil
}
