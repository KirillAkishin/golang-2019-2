package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func GetValueCalculatedByRPN(input string) (result float64, err error) {
	err = nil
	strs := strings.Split(input, " ")
	stack := make([]float64, len(strs))
	sp := 0
	strs[len(strs)-1] = strs[len(strs)-1][:len(strs[len(strs)-1])-1] //удаление символа переноса строки
	for idx := range strs {
		symbols := utf8.RuneCountInString(strs[idx])
		if symbols == 1 {
			switch strs[idx] {
			case "=":
				{
					break
				}
			case "+":
				if sp < 2 {
					err = errors.New("1")
					fmt.Print("Невозможно выполнить POP для пустого стека.\n")
					return result, err
				}
				stack[sp-2] = stack[sp-2] + stack[sp-1]
				sp--
			case "-":
				if sp < 2 {
					err = errors.New("1")
					fmt.Print("Невозможно выполнить POP для пустого стека.\n")
					return result, err
				}
				stack[sp-2] = stack[sp-2] - stack[sp-1]
				sp--
			case "*":
				if sp < 2 {
					err = errors.New("1")
					fmt.Print("Невозможно выполнить POP для пустого стека.\n")
					return result, err
				}
				stack[sp-2] = stack[sp-1] * stack[sp-2]
				sp--
			case "/":
				if sp < 2 {
					err = errors.New("1")
					fmt.Print("Невозможно выполнить POP для пустого стека.\n")
					return result, err
				}
				if stack[sp-1] == 0. {
					err = errors.New("3")
					fmt.Print("Деление на ноль.\n")
					return result, err
				}
				stack[sp-2] = stack[sp-2] / stack[sp-1]
				sp--
			default:
				x, err := strconv.ParseFloat(strs[idx], 64)
				if err == nil {
					stack[sp] = x
					sp++
				} else {
					err = errors.New("2")
					fmt.Print("Can't read number\n")
					return result, err
				}
			}
		} else {
			x, err := strconv.ParseFloat(strs[idx], 64)
			if err == nil {
				stack[sp] = x
				sp++
			} else {
				err = errors.New("2")
				fmt.Print("Can't read number\n")
				return result, err
			}
		}
	}
	result = stack[sp-1]
	return result, err
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	result, err := GetValueCalculatedByRPN(input)
	if err == nil {
		fmt.Print("Result = ", result, "\n")
	}
	return
}
