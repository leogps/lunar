/*
 * Copyright (c) 2024, Paul Gundarapu.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 */

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// PromptAndValidate prompts the user for input and validates it based on the type.
func PromptAndValidate[T any](prompt string) (T, error) {
	var zero T // zero value for T, to return on error
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Determine the type of T and parse accordingly
		switch any(zero).(type) {
		case int:
			value, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid input. Expected an int value")
				continue
			}
			return any(value).(T), nil

		case float64:
			value, err := strconv.ParseFloat(input, 64)
			if err != nil {
				fmt.Println("Invalid input. Expected dollar amount")
				continue
			}
			return any(value).(T), nil

		case string:
			return any(input).(T), nil

		case bool:
			inputNormalized := strings.ToLower(input)
			// Check for acceptable boolean values
			if inputNormalized == "true" || inputNormalized == "yes" || inputNormalized == "y" {
				return any(true).(T), nil
			} else if inputNormalized == "false" || inputNormalized == "no" || inputNormalized == "n" {
				return any(false).(T), nil
			} else {
				fmt.Println("Invalid input. Expected one of: 'true', 'false', 'yes', 'no', 'y', 'n')")
				continue
			}

		default:
			fmt.Printf("Unsupported type: %s\n", reflect.TypeOf(zero).Kind())
			return zero, fmt.Errorf("unsupported type")
		}
	}
}
