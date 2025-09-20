package main

import "fmt"

func main() {
	var (
		num1, num2 int
		operation  string
	)

	_, err := fmt.Scan(&num1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err2 := fmt.Scan(&num2)
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err3 := fmt.Scan(&operation)
	if err3 != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operation {
	case "+":
		fmt.Println(num1 + num2)
	case "-":
		fmt.Println(num1 - num2)
	case "*":
		fmt.Println(num1 * num2)
	case "/":
		if num2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(num1 / num2)
	default:
		fmt.Println("Invalid operation")
	}
}
