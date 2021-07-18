package git

var CreateFakeGit = func() *fakeGitImpl {
	return &fakeGitImpl{}
}

type fakeGitImpl struct {
	Git
	CloneMock                    func(url string, branch string, clonePath string) error
	InitMock                     func(path string) error
	AddOriginMock                func(path string, url string) error
	FetchShallowMock             func(path string, branch string) error
	ResetMock                    func(path string, revision string) error
	CheckoutMock                 func(path string, branch string) error
	CleanMock                    func(path string) error
	GetRemoteHeadCommitHashMock  func(path string, repoUrl string, branch string) (string, error)
	GetLocalOriginCommitHashMock func(path string, branch string) (string, error)
	LogRevisionsDiffPrettyMock   func(path string, prevRevision string, newRevision string) error
}

func (g *fakeGitImpl) Clone(url string, branch string, clonePath string) error {
	return g.CloneMock(url, branch, clonePath)
}

func (g *fakeGitImpl) Init(path string) error {
	return g.InitMock(path)
}

func (g *fakeGitImpl) AddOrigin(path string, url string) error {
	return g.AddOriginMock(path, url)
}

func (g *fakeGitImpl) FetchShallow(path string, branch string) error {
	return g.FetchShallowMock(path, branch)
}

func (g *fakeGitImpl) Reset(path string, revision string) error {
	return g.ResetMock(path, revision)
}

func (g *fakeGitImpl) Checkout(path string, branch string) error {
	return g.CheckoutMock(path, branch)
}

func (g *fakeGitImpl) Clean(path string) error {
	return g.CleanMock(path)
}

func (g *fakeGitImpl) GetRemoteHeadCommitHash(path string, repoUrl string, branch string) (string, error) {
	return g.GetRemoteHeadCommitHashMock(path, repoUrl, branch)
}

func (g *fakeGitImpl) GetLocalOriginCommitHash(path string, branch string) (string, error) {
	return g.GetLocalOriginCommitHashMock(path, branch)
}

func (g *fakeGitImpl) LogRevisionsDiffPretty(path string, prevRevision string, newRevision string) error {
	return g.LogRevisionsDiffPrettyMock(path, prevRevision, newRevision)
}
