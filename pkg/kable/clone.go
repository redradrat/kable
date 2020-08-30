package kable

import (
	"github.com/go-git/go-git/v5"
)

func clone(cloner ClonerConfig) error {
	out, err := git.PlainClone(cloner.BaseDir, false, &cloner.CloneOptions)
	if err != nil {
		return err
	}
	out.Config()
	return nil

}
