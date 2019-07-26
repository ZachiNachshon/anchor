package locator

type Locator interface {
	Scan() error
	Print(opts *ListOpts)
	Name(identifier string) (string, error)
	Dockerfile(identifier string) (string, error)
	DockerContext(identifier string) (string, error)
	Manifest(identifier string) (string, error)
}
