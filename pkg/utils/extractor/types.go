package extractor

type DockerCommand string

const (
	DockerCommandRun   DockerCommand = "docker run"
	DockerCommandBuild DockerCommand = "docker build"
	DockerCommandPush  DockerCommand = "docker push"
	DockerCommandTag   DockerCommand = "docker tag"
)

type ManifestCommand string

const (
	ManifestCommandPortForward ManifestCommand = "kubectl port-forward"
	ManifestCommandWait        ManifestCommand = "kubectl wait"
	ManifestCommandStateful    ManifestCommand = "anchor-stateful"
)

type Extractor interface {
	DockerCmd(identifier string, dockerCommand DockerCommand) (string, error)
	ManifestCmd(identifier string, manifestCommand ManifestCommand) (string, error)
	ManifestContent(identifier string, manifestCommand ManifestCommand) (bool, error)
}
