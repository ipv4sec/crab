package utils

import (
	"github.com/go-git/go-git/v5"
)

func InitRepo(path, url string) error {
	r, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	return w.Pull(&git.PullOptions{})
}
