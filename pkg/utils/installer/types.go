package installer

type Installer interface {
	install() error
	verify() error
	Check() error
}
