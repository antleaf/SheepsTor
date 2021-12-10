package main

import (
	"github.com/go-git/go-git/v5"
	_ "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

func Clone(cloneID, branchRef, repoLocalPath string) error {
	var err error
	publicKey, err := getSshPublicKey()
	if err != nil {
		return err
	}
	_, err = git.PlainClone(repoLocalPath, false, &git.CloneOptions{
		URL:           cloneID,
		Auth:          publicKey,
		ReferenceName: plumbing.ReferenceName(branchRef),
		SingleBranch:  true,
		Progress:      os.Stdout,
	})
	return err
}

func Pull(repoLocalPath, branchRef string) error {
	var err error
	repo, err := git.PlainOpen(repoLocalPath)
	if err != nil {
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	publicKey, err := getSshPublicKey()
	if err != nil {
		return err
	}
	err = w.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.ReferenceName(branchRef),
		Auth:          publicKey,
		Progress:      os.Stdout,
	})
	if err != nil {
		switch err.Error() {
		//TODO find better way to do this checking type of error rather than  checking error string
		case "already up-to-date":
			err = nil
		//case "non-fast-forward update":
		//	zapLogger.Info(fmt.Sprintf("Non-fast-forward update for '%s'", website.Id))
		//	err = nil
		default:
			//zapLogger.Error(err.Error())
			return err
		}
	}
	return err
}

func CommitAndPush(repoLocalPath, message string) error {
	var err error
	repo, err := git.PlainOpenWithOptions(repoLocalPath, &git.PlainOpenOptions{DetectDotGit: true, EnableDotGitCommonDir: true})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	err = w.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	commitHash, err := w.Commit(message, &git.CommitOptions{All: true})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.Infof("Committed changes to '%s' with commit hash: '%s'", repoLocalPath, commitHash.String())
	publicKey, err := getSshPublicKey()
	if err != nil {
		logger.Error("SSH Key not returned")
		return err
	}
	err = repo.Push(&git.PushOptions{RemoteName: "origin", Auth: publicKey, Progress: os.Stdout})
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.Info("Pushed changes")
	return err
}

func getSshPublicKey() (*ssh.PublicKeys, error) {
	var publicKey *ssh.PublicKeys
	usr, err := user.Current()
	if err != nil {
		return publicKey, err
	}
	privateSSHKeyPath := filepath.Join(usr.HomeDir, ".ssh", "id_rsa")
	sshKey, err := ioutil.ReadFile(privateSSHKeyPath)
	if err != nil {
		return publicKey, err
	}
	publicKey, err = ssh.NewPublicKeys("git", []byte(sshKey), "")
	if err != nil {
		return publicKey, err
	}
	return publicKey, err
}
