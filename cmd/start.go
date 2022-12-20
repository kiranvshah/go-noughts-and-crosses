/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// boardType represents a board - a slice containing 3 slices, each containing 3 strings
type boardType [3][3]string

// printBoard displays a boardType by printing it to the command line
func printBoard(board boardType) {
	fmt.Println("    A B C")
	fmt.Println()
	for rowIndex, row := range board {
		fmt.Printf("%d   ", rowIndex+1)
		for _, cellValue := range row {
			fmt.Print(cellValue + " ")
		}
		fmt.Println()
	}
}

// getCellCoords prompts the user to enter coordiantes of a particular cell, and returns the cell's coordinates, or an error.
// It takes an input of the current board, so it can throw an error if the selected coords have already been played.
func getCellCoords(board boardType) (x, y int, err error) {
	prompt := promptui.Prompt{
		Label: "Please enter the coordinates of the cell you'd like to play in (e.g. A1):",
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . | green }} ",
			Invalid: "{{ . | red }} ",
			Success: "{{ . | bold }} ",
		},
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.New("no cell coordinates provided")
			}
			if len(input) != 2 {
				return errors.New("cell coordiantes should be two characters long")
			}
			if !unicode.IsLetter(rune(input[0])) {
				return errors.New("first coordinate should be a letter")
			}
			if !unicode.IsNumber(rune(input[1])) {
				return errors.New("second coordinate should be a number")
			}
			firstChar := unicode.ToUpper(rune(input[0]))
			if firstChar != []rune("A")[0] && firstChar != []rune("B")[0] && firstChar != []rune("C")[0] {
				return errors.New("first coordiante should be A, B or C")
			}
			secondChar, errWhenConvertingSecondChar := strconv.ParseInt(string(input[1]), 10, 64)
			if errWhenConvertingSecondChar != nil {
				return errors.New("error when converting second coordinate to a number")
			}
			if secondChar != 1 && secondChar != 2 && secondChar != 3 {
				return errors.New("second coordinate should be 1, 2, or 3")
			}
			return nil
		},
	}

	coords, err := prompt.Run()
	firstChar := unicode.ToUpper(rune(coords[0]))
	secondChar, _ := strconv.ParseInt(string(coords[1]), 10, 64)

	if err != nil {
		// set x variable
		if firstChar == []rune("A")[0] {
			x = 0
		} else if firstChar == []rune("B")[0] {
			x = 1
		} else if firstChar == []rune("C")[0] {
			x = 2
		}

		// set y variable
		y = int(secondChar) - 1

		if board[y][x] != "-" {
			err = errors.New("selected coords have already been played")
		}
	}

	return // naked return
}

// boardHasThreeInARow, given an input of a boardType, returns a bool representing whether someone has won on the board,
// and a string representing the winner ("o", "x" or "")
func boardHasThreeInARow(board boardType) (hasThreeInRow bool, winner string) {

	// check for horizontal three in a row
	for _, row := range board {
		oCount := 0
		xCount := 0
		for _, cellValue := range row {
			if cellValue == "o" {
				oCount++
			} else if cellValue == "x" {
				xCount++
			}
		}
		if oCount == 3 {
			hasThreeInRow = true
			winner = "o"
			return
		}
		if xCount == 3 {
			hasThreeInRow = true
			winner = "x"
			return
		}
	}

	// check for vertical three in a row
	for colIndex := 0; colIndex < 3; colIndex++ {
		oCount := 0
		xCount := 0
		for _, row := range board {
			cellValue := row[colIndex]
			if cellValue == "o" {
				oCount++
			} else if cellValue == "x" {
				xCount++
			}
		}
		if oCount == 3 {
			hasThreeInRow = true
			winner = "o"
			return
		}
		if xCount == 3 {
			hasThreeInRow = true
			winner = "x"
			return
		}
	}

	{ // check for "\" diagonal
		oCount := 0
		xCount := 0
		for rowAndColIndex := 0; rowAndColIndex < 3; rowAndColIndex++ {
			cellValue := board[rowAndColIndex][rowAndColIndex]
			if cellValue == "o" {
				oCount++
			} else if cellValue == "x" {
				xCount++
			}
		}
		if oCount == 3 {
			hasThreeInRow = true
			winner = "o"
			return
		}
		if xCount == 3 {
			hasThreeInRow = true
			winner = "x"
			return
		}
	}

	{ // check for "/" diagonal
		oCount := 0
		xCount := 0
		for rowIndex := 0; rowIndex < 3; rowIndex++ {
			colIndex := 2 - rowIndex
			cellValue := board[rowIndex][colIndex]
			if cellValue == "o" {
				oCount++
			} else if cellValue == "x" {
				xCount++
			}
		}
		if oCount == 3 {
			hasThreeInRow = true
			winner = "o"
			return
		}
		if xCount == 3 {
			hasThreeInRow = true
			winner = "x"
			return
		}
	}

	return
}

// boardIsDraw, given an input of a boardType, returns a bool representing whether the game is a draw and all the spaces are full
func boardIsDraw(board boardType) bool {
	blankSpaceCount := 0
	for _, row := range board {
		for _, cellValue := range row {
			if cellValue == "-" {
				blankSpaceCount++
			}
		}
	}

	return blankSpaceCount == 0
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a game of noughts and crosses",
	Long:  `Starts a game of noughts and crosses to be played on the command line`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")

		board := boardType{{"-", "-", "-"}, {"-", "-", "-"}, {"-", "-", "-"}}
		printBoard(board)

		nextToMove := "x"

		for {

			// check for winner
			isWinner, winner := boardHasThreeInARow(board)
			if isWinner {
				fmt.Printf("%s wins!\n", winner)
				return
			}

			// check for draw
			if boardIsDraw(board) {
				fmt.Println("Draw! Game over")
				return
			}

			// player has turn
			fmt.Printf("%s to move.", nextToMove)
			successfullyGotCoords := false
			for !successfullyGotCoords {
				x, y, err := getCellCoords(board)
				if err != nil {
					fmt.Printf("Error getting coordiantes: %s", err)
				} else {
					// place go
					successfullyGotCoords = true
					board[y][x] = nextToMove

				}
			}

			// switch next player to have turn
			if nextToMove == "x" {
				nextToMove = "o"
			} else {
				nextToMove = "x"
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
