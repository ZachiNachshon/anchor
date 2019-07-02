package locator

type Locator interface {
	Dockerfile(name string) (string, error)
	DockerfileDir(name string) (string, error)
	Manifest(name string) (string, error)
	ManifestDir(name string) (string, error)

	GetRootFromManifestFile(path string) string
	GetRootFromDockerfile(path string) string
}
