package utils

import (
	"github.com/go-git/go-git/v5"
	"os"
)

func InitRepo(path, url string) error {
	var r *git.Repository
	if _, err := os.Stat(path); os.IsNotExist(err) {
		r, err = git.PlainClone(path, false, &git.CloneOptions{
			URL: url,
		})
		if err != nil {
			return err
		}
	} else {
		r, err = git.PlainOpen(path)
		if err != nil {
			return err
		}
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	return w.Pull(&git.PullOptions{})
}
