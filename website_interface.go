package sheepstor

type WebsiteInterface interface {
	Build() error
	ProvisionSources() error
	CommitAndPush(message string) error
	HasID(id string) bool
	HasRepoNameAndBranchRef(repoName, branchRef string) bool
	GetGitRepo() GitRepo
	GetID() string
}
