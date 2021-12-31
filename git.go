package sheepstor

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

func NewGitRepo(cloneID, repoName, branchName, localPath string) GitRepo {
	var g = GitRepo{
		CloneID:       cloneID,
		RepoName:      repoName,
		BranchName:    branchName,
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
		return headCommitID
	}
	ref, err := repo.Head()
	if err != nil {
		return headCommitID
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
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
			logger.Debugf("Already up-to-date for '%s'", g.RepoLocalPath)
			err = nil
		case "non-fast-forward update":
			logger.Debugf("Non-fast-forward update for '%s'", g.RepoLocalPath)
			err = nil
		default:
			return err
		}
	}
	return err
}

func (g *GitRepo) CommitAndPush(message string) error {
	var err error
	repo, err := git.PlainOpenWithOptions(g.RepoLocalPath, &git.PlainOpenOptions{DetectDotGit: true, EnableDotGitCommonDir: true})
	if err != nil {
		return err
	}
	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	err = w.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return err
	}
	_, err = w.Commit(message, &git.CommitOptions{All: true})
	if err != nil {
		return err
	}
	publicKey, err := getSshPublicKey()
	if err != nil {
		return err
	}
	err = repo.Push(&git.PushOptions{RemoteName: "origin", Auth: publicKey, Progress: nil})
	if err != nil {
		return err
	}
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
