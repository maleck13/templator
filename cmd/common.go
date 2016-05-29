package cmd

import (
	"bufio"
	"fmt"
	"os"
)

func QuestionAndAnswer(q string, answer func(string)) error {
	fmt.Print(q)
	lineScanner := bufio.NewScanner(os.Stdin)
	lineScanner.Scan()
	answer(lineScanner.Text())

	if err := lineScanner.Err(); err != nil {
		return err
	}
	return nil

}
