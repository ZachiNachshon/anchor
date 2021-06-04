package clipboard

type Clipboard interface {
	Load(content string) error
}
