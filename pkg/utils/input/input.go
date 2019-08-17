package input

import (
	"bufio"
	"fmt"
	"github.com/ZachiNachshon/anchor/pkg/logger"
	"os"
	"strconv"
)

type YesNoInput struct{}

func NewYesNoInput() *YesNoInput {
	return &YesNoInput{}
}

func (input *YesNoInput) WaitForInput(question string) (bool, error) {
	fmt.Print("\n  " + question + " (y/n): ")

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

	fmt.Print("\n")

	return selection, nil
}

type NumericInput struct{}

func NewNumericInput() *NumericInput {
	return &NumericInput{}
}

func (input *NumericInput) WaitForInput() (int, error) {
	return input.waitForInputInner(false)
}

func (input *NumericInput) WaitForInputAllowDefault() (int, error) {
	return input.waitForInputInner(true)
}

func (input *NumericInput) waitForInputInner(allowDefault bool) (int, error) {
	fmt.Print("  Enter a value: ")

	//_ = exec.Command(string(shell.BASH), "-c", "stty sane")

	// Check is shell is zsh and execute 'stty sane' to fix the ^M char for enter key press
	//_, _ = common.ShellExec.ExecuteWithOutput("stty sane")

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
			if allowDefault {
				return -1, nil
			}
			fmt.Print("  Enter a value: ")
			break
		default:
			if number, err := strconv.Atoi(value); err != nil {
				if allowDefault {
					return -1, nil
				}
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
