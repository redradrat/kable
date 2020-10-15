package helm

import (
	"github.com/go-git/go-git/v5"
)

// Checks out the given repository to the given directory and returns the path
func Checkout(url, path string) error {
	repo, err := git.PlainOpen(path)
	missing := err == git.ErrRepositoryNotExists
	if err != nil && !missing {
		return err
	}

	if !missing {
		wt, err := repo.Worktree()
		if err != nil {
			return err
		}
		if err := wt.Pull(&git.PullOptions{}); err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}
	} else {
		_, err := git.PlainClone(path, false, &git.CloneOptions{
			URL:          url,
			SingleBranch: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
