package input

var CreateFakeUserInput = func() *fakeUserInputImpl {
	return &fakeUserInputImpl{}
}

type fakeUserInputImpl struct {
	UserInput
	AskYesNoQuestionMock        func(question string) (bool, error)
	AskForNumberMock            func() (int, error)
	AskForNumberWithDefaultMock func() (int, error)
	PressAnyKeyToContinueMock   func() error
}

func (in *fakeUserInputImpl) AskYesNoQuestion(question string) (bool, error) {
	return in.AskYesNoQuestionMock(question)
}

func (in *fakeUserInputImpl) AskForNumber() (int, error) {
	return in.AskForNumberMock()
}

func (in *fakeUserInputImpl) AskForNumberWithDefault() (int, error) {
	return in.AskForNumberWithDefaultMock()
}

func (in *fakeUserInputImpl) PressAnyKeyToContinue() error {
	return in.PressAnyKeyToContinueMock()
}
