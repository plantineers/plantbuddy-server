package utils

import (
	"fmt"
	"io"
	"os"
)

// Copy database file from current directory to the data directory
// if it does not exist
func CopyDatabase() {
	if _, err := os.Stat("buddy_data/buddy.sqlite"); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Copying database file...")
			err := os.MkdirAll("buddy_data", 0755)
			if err != nil {
				panic(err)
			}
			dest, err := os.Create("buddy_data/buddy.sqlite")
			if err != nil {
				panic(err)
			}
			src, err := os.Open("buddy.sqlite")
			if err != nil {
				panic(err)
			}
			_, err = io.Copy(dest, src)
			if err != nil {
				panic(err)
			}
			dest.Close()
			src.Close()
			fmt.Println("Done copying database file.")
		}
	}
}
