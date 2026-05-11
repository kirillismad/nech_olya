package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	data := []byte("Hello, файл!\n")
	err := os.WriteFile("hello_file.txt", data, 0644)
	if err != nil {
		fmt.Println("failed to write to file")
		return
	}
	file, err := os.OpenFile("hello_file.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open file")
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	w.WriteString("Good morning!\n")
	w.Flush()
	w.WriteString("Good night!\n")
	w.Flush()

	data, err = os.ReadFile("hello_file.txt")
	if err != nil {
		fmt.Println("file not exist")
		return
	}
	fmt.Println(string(data))
}
