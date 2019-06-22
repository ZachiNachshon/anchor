package input

type Interactive interface {
	WaitForInput(question string) (bool, error)
}
