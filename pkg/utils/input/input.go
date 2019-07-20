package input

import (
	"bufio"
	"fmt"
	"github.com/anchor/pkg/logger"
	"os"
	"strconv"
)

type YesNoInput struct{}

func NewYesNoInput() *YesNoInput {
	return &YesNoInput{}
}

func (input *YesNoInput) WaitForInput(question string) (bool, error) {
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

type NumericInput struct{}

func NewNumericInput() *NumericInput {
	return &NumericInput{}
}

func (input *NumericInput) WaitForInput() (int, error) {
	fmt.Print("  Enter a value: ")

	reader := bufio.NewReader(os.Stdin)

	var selection int
	var keepAsking = true
	for keepAsking {
		char, _, err := reader.ReadLine()
		value := string(char)
		if err != nil {
			return -1, err
		}

		switch value {
		case "\n":
			fmt.Print("  Enter a value: ")
			break
		default:
			if number, err := strconv.Atoi(value); err != nil {
				logger.Info("Selection must be a numeric value")
				fmt.Print("\n  Enter a value: ")
			} else {
				selection = number
				keepAsking = false
				break
			}
		}
	}

	return selection, nil
}
