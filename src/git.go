package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	_ "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

type GitRepo struct {
	CloneID       string
	RepoName      string
	BranchName    string
	BranchRef     string
	RepoLocalPath string
}

func NewGitRepo(gitConfig GitRepoConfig, localPath string) GitRepo {
	var g = GitRepo{
		CloneID:       gitConfig.CloneId,
		RepoName:      gitConfig.RepoName,
		BranchName:    gitConfig.BranchName,
		RepoLocalPath: localPath,
	}
	g.BranchRef = fmt.Sprintf("refs/heads/%s", g.BranchName)
	return g
}

func (g *GitRepo) GetHeadCommitID() string {
	var err error
	var headCommitID string
	repo, err := git.PlainOpen(g.RepoLocalPath)
	if err != nil {
		logger.Error(err.Error())
		return headCommitID
	}
	ref, err := repo.Head()
	if err != nil {
		logger.Error(err.Error())
		return headCommitID
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		logger.Error(err.Error())
		return headCommitID
	}
	headCommitID = commit.ID().String()
	return headCommitID
}

func (g *GitRepo) Clone() error {
	var err error
	publicKey, err := getSshPublicKey()
	if err != nil {
		return err
	}
	_, err = git.PlainClone(g.RepoLocalPath, false, &git.CloneOptions{
		URL:           g.CloneID,
		Auth:          publicKey,
		ReferenceName: plumbing.ReferenceName(g.BranchRef),
		SingleBranch:  true,
		Progress:      nil,
	})
	return err
}

func (g *GitRepo) Pull() error {
	var err error
	repo, err := git.PlainOpen(g.RepoLocalPath)
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
		ReferenceName: plumbing.ReferenceName(g.BranchRef),
		Auth:          publicKey,
		Progress:      nil,
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

func (g *GitRepo) CommitAndPush(message string) error {
	var err error
	repo, err := git.PlainOpenWithOptions(g.RepoLocalPath, &git.PlainOpenOptions{DetectDotGit: true, EnableDotGitCommonDir: true})
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
	logger.Infof("Committed changes to '%s' with commit hash: '%s'", g.RepoLocalPath, commitHash.String())
	publicKey, err := getSshPublicKey()
	if err != nil {
		logger.Error("SSH Key not returned")
		return err
	}
	err = repo.Push(&git.PushOptions{RemoteName: "origin", Auth: publicKey, Progress: nil})
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
