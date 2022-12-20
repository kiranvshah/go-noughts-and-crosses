/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// boardType represents a board - a slice containing 3 slices, each containing 3 strings
type boardType [3][3]string

// printBoard displays a boardType by printing it to the command line
func printBoard(board boardType) {
	// fmt.Println("\033[H\033[2J    A B C")
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

// getCellCoords prompts the user to enter coordiantes of a particular cell, and returns the cell's coordinates (including human friendly ones, e.g. B3), or an error.
// It takes an input of the current board, so it can throw an error if the selected coords have already been played.
func getCellCoords(board boardType) (x, y int, humanFriendlyCoords string, err error) {

	promptCol := promptui.Select{
		Label: "Column",
		Items: []string{
			"A", "B", "C",
		},
		HideSelected: true,
	}
	x, xHumanFriendly, errCol := promptCol.Run()
	if errCol != nil {
		err = errCol
		return
	}

	promptRow := promptui.Select{
		Label: "Row",
		Items: []string{
			"1", "2", "3",
		},
		HideSelected: true,
	}
	y, yHumanFriendly, errRow := promptRow.Run()
	if errRow != nil {
		err = errRow
		return
	}

	fmt.Println("") // https://github.com/manifoldco/promptui/issues/180#issuecomment-882836869

	if board[y][x] != "-" {
		err = errors.New("selected coords have already been played")
	}

	humanFriendlyCoords = xHumanFriendly + yHumanFriendly

	return // naked return
}

// computerGivesCellCoords simulates a computer's turn by providing the cell coords that it is going to play.
func computerGivesCellCoords(board boardType) (x, y int, humanFriendlyCoords string) {
	availableCoords := make([][2]int, 0, 9)
	for rowIndex, row := range board {
		for colIndex, cellValue := range row {
			if cellValue == "-" {
				availableCoords = append(availableCoords, [2]int{rowIndex, colIndex})
			}
		}
	}

	selectedCoords := availableCoords[rand.Intn(len(availableCoords))]
	x, y = selectedCoords[0], selectedCoords[1]

	if x == 0 {
		humanFriendlyCoords += "A"
	} else if x == 1 {
		humanFriendlyCoords += "B"
	} else if x == 2 {
		humanFriendlyCoords += "C"
	}
	humanFriendlyCoords += fmt.Sprintf("%d", y+1)

	return
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
		againstComputer, errWithAgainstComputerFlag := cmd.Flags().GetBool("against-computer")
		if errWithAgainstComputerFlag != nil {
			fmt.Printf("Error with against-computer flag: %s\n", errWithAgainstComputerFlag)
		}
		if againstComputer {
			fmt.Print("Starting singleplayer game against computer. You are o.\n\n")
		} else {
			fmt.Print("Starting two-player game.\n\n")
		}

		board := boardType{{"-", "-", "-"}, {"-", "-", "-"}, {"-", "-", "-"}}
		nextToMove := "x"

		for {

			// check for winner
			isWinner, winner := boardHasThreeInARow(board)
			if isWinner {
				printBoard(board)
				fmt.Printf("%s wins!\n", winner)
				return
			}

			// check for draw
			if boardIsDraw(board) {
				printBoard(board)
				fmt.Println("Draw! Game over")
				return
			}

			if nextToMove == "o" && againstComputer {
				// computer has turn
				x, y, humanFriendlyCoords := computerGivesCellCoords(board)
				board[y][x] = nextToMove
				fmt.Printf("Computer placed %s at %s.\n\n\n\n\n", nextToMove, humanFriendlyCoords)
			} else {
				// player has turn
				printBoard(board)
				fmt.Printf("%s to move.\n\n", nextToMove)
				carryOnTryingToGetCoords := true
				for carryOnTryingToGetCoords {
					x, y, humanFriendlyCoords, err := getCellCoords(board)
					if err != nil {
						if err.Error() == "selected coords have already been played" {
							fmt.Print("Selected coords have already been played. Try again.\n")
						} else {
							fmt.Printf("Error getting coordiantes: %s\n", err)
							return
						}
					} else {
						// place go
						board[y][x] = nextToMove
						fmt.Printf("Placed %s at %s.\n\n\n\n\n", nextToMove, humanFriendlyCoords)
						carryOnTryingToGetCoords = false
					}
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
	rand.Seed(time.Now().UnixNano())

	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().BoolP("against-computer", "C", false, "Whether one player (o) should be controlled by the computer.")
}
