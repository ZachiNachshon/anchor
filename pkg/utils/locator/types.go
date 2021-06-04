package locator

type Locator interface {
	Scan() error
	Applications() []string
	Application(name string) *AppContent
	//Print()
}
