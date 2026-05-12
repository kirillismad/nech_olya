package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("hello.txt")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("errors: file does not exist")
		}
		fmt.Printf("file %s reading error", err)
		return
	}
	fmt.Println(string(data))

	file, err := os.Open("hello.txt")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("errors: file does not exist")
		}
		fmt.Printf("file %s reading error", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	num := 1
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("%s номер строки: %d\n", line, num)
		num++
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %v\n", err)
	}

	data2, err2 := os.ReadFile("congig.yaml")
	if err2 != nil {
		if errors.Is(err2, os.ErrNotExist) {
			fmt.Printf("errors: file does not exist")
		}
		fmt.Printf("file %s reading error", err2)
		return
	}
	fmt.Println(string(data2))

}
