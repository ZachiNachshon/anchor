package input

import (
	"bufio"
	"fmt"
	"os"
)

type yesNoInput struct {
}

func NewYesNoInput() Interactive {
	return &yesNoInput{}
}

func (input *yesNoInput) WaitForInput(question string) (bool, error) {
	fmt.Print(question + " (y/n): ")

	reader := bufio.NewReader(os.Stdin)

	var selection bool
	var keepAsking = true
	for keepAsking {
		char, _, err := reader.ReadRune()
		if err != nil {
			return false, err
		}

		switch char {
		case 'y':
			selection = true
			keepAsking = false
			break
		case 'n':
			selection = false
			keepAsking = false
			break
		case '\n':
			// Do nothing
			break
		default:
			fmt.Print(question + " (y/n): ")
			break
		}
	}

	return selection, nil
}
