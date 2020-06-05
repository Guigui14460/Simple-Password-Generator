package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

const (
	// LowerLetters is the list of lowercase letters.
	LowerLetters = "abcdefghijklmnopqrstuvwxyz"
	// UpperLetters is the list of uppercase letters.
	UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Digits is the list of permitted digits.
	Digits = "0123456789"
	// Symbols is the list of permitted symbols.
	Symbols = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
)

var (
	// ErrExceedsTotalLength is the error returned when the number of digits and
	// symbols is greater than the total length.
	ErrExceedsTotalLength = errors.New("number of digits and symbols must be less than total length")
	// ErrLettersExceedsAvailable is the error returned when the number of letters
	// exceeds the number of available letters and repeats are not allowed.
	ErrLettersExceedsAvailable = errors.New("number of letters exceeds available letters and repeats are not allowed")
	// ErrDigitsExceedsAvailable is the error returned when the number of digits
	// exceeds the number of available digits and repeats are not allowed.
	ErrDigitsExceedsAvailable = errors.New("number of digits exceeds available digits and repeats are not allowed")
	// ErrSymbolsExceedsAvailable is the error returned when the number of symbols
	// exceeds the number of available symbols and repeats are not allowed.
	ErrSymbolsExceedsAvailable = errors.New("number of symbols exceeds available symbols and repeats are not allowed")
)

// Generator is the stateful generator which can be used to customize the list
// of letters, digits, and/or symbols.
type Generator struct {
	lowerLetters string
	upperLetters string
	digits       string
	symbols      string
}

// GeneratorInput is used as input to the NewGenerator function.
type GeneratorInput struct {
	LowerLetters string
	UpperLetters string
	Digits       string
	Symbols      string
}

/*
Function which creates a new generator from a specified configuration.
	Parameters:
	-----------
		i (*GeneratorInput): specified configuration
			Note: if i == nil, we use default values

	Returns:
	--------
		*Generator - a generator pointor
*/
func NewGenerator(i *GeneratorInput) *Generator {
	// Put the default values
	if i == nil {
		i = new(GeneratorInput)
	}

	// Create the Generator (we save here the pointer to access easily attributes of the object)
	g := &Generator{
		lowerLetters: i.LowerLetters,
		upperLetters: i.UpperLetters,
		digits:       i.Digits,
		symbols:      i.Symbols,
	}

	// If the value is "", we put the default associated value
	if g.lowerLetters == "" {
		g.lowerLetters = LowerLetters
	}
	if g.upperLetters == "" {
		g.upperLetters = UpperLetters
	}
	if g.digits == "" {
		g.digits = Digits
	}
	if g.symbols == "" {
		g.symbols = Symbols
	}

	return g
}

/*
Function to generate a password with the required arguments.
	Method of Generator type

	Parameters:
	-----------
		length (int): total number of characters
		numDigits (int): number of digits to include
		numSymbols (int): number of symbols to include
		allowUpper (bool): include uppercase
		allowRepeat (bool): allows repeat characters

	Returns:
	--------
		string, error - password and the error if the password was not generated
*/
func (g *Generator) Generate(length, numDigits, numSymbols int, allowUpper, allowRepeat bool) (string, error) {
	// Get all possibles letters
	letters := g.lowerLetters
	if allowUpper {
		letters += g.upperLetters
	}

	// Verify if it is possible to generate a password
	chars := length - numDigits - numSymbols
	if chars < 0 {
		return "", ErrExceedsTotalLength
	}
	if !allowRepeat && chars > len(letters) {
		return "", ErrLettersExceedsAvailable
	}
	if !allowRepeat && numDigits > len(g.digits) {
		return "", ErrDigitsExceedsAvailable
	}
	if !allowRepeat && numSymbols > len(g.symbols) {
		return "", ErrSymbolsExceedsAvailable
	}

	// Creation of the password
	var result string

	// Characters
	for i := 0; i < chars; i++ {
		// Choice a letter
		ch, err := randomElement(letters)
		if err != nil {
			return "", err
		}
		// Not add the choiced letter if is already there (only if allowRepeat is false)
		// Cancel of the insertion
		if !allowRepeat && strings.Contains(result, ch) {
			i--
			continue
		}
		// Insertion
		result, err = randomInsert(result, ch)
		if err != nil {
			return "", err
		}
	}

	// Digits
	for i := 0; i < numDigits; i++ {
		// Choice a digit
		d, err := randomElement(g.digits)
		if err != nil {
			return "", err
		}
		// Not add the choiced digit if is already there (only if allowRepeat is false)
		// Cancel of the insertion
		if !allowRepeat && strings.Contains(result, d) {
			i--
			continue
		}
		// Insertion
		result, err = randomInsert(result, d)
		if err != nil {
			return "", err
		}
	}

	// Symbols
	for i := 0; i < numSymbols; i++ {
		// Choice a symbol
		sym, err := randomElement(g.symbols)
		if err != nil {
			return "", err
		}
		// Not add the choiced symbol if is already there (only if allowRepeat is false)
		// Cancel of the insertion
		if !allowRepeat && strings.Contains(result, sym) {
			i--
			continue
		}
		// Insertion
		result, err = randomInsert(result, sym)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

/*
Function which randomly insert the given value into the given string
	Parameters:
	-----------
		str (int): string to use for insertion
		val (string): value to insert

	Returns:
	--------
		string, error - string where the given value was inserted and the error if value not inserted
*/
func randomInsert(str, val string) (string, error) {
	// Verify empty string value
	if str == "" {
		return val, nil
	}

	// Initialize the random system and get a random value
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(str)+1)))
	if err != nil {
		return "", err
	}
	i := n.Int64()

	// Insertion of the given value
	return str[0:i] + val + str[i:], nil
}

/*
Function which randomly return a value from a given string
	Parameters:
	-----------
		str (int): string to use

	Returns:
	--------
		string, error - extracted value was inserted and the error if value not inserted
*/
func randomElement(str string) (string, error) {
	// Initialize the random system and get a random value
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(str))))
	if err != nil {
		return "", err
	}
	// Directly return the choiced value
	return string(str[n.Int64()]), nil
}

func main() {
	// Initialize variables
	var length, numDigits, numSymbols int64
	var allowUpper, allowRepeat bool = true, true
	var err error
	scanner := bufio.NewScanner(os.Stdin)

	// Get the positionned arguments
	args := os.Args[1:]

	// Open interactive program
	if len(args) == 0 {
		print("Length of the password : ")
		scanner.Scan()
		length, err = strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			panic(err.Error())
		}
		print("Total number of digits : ")
		scanner.Scan()
		numDigits, err = strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			panic(err.Error())
		}
		print("Total number of symbols : ")
		scanner.Scan()
		numSymbols, err = strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			panic(err.Error())
		}
		print("Activate the uppercase (false for NO, true for YES) : ")
		scanner.Scan()
		allowUpper, err = strconv.ParseBool(scanner.Text())
		if err != nil {
			panic(err.Error())
		}
		print("Activate the character repeat (false for NO, true for YES) : ")
		scanner.Scan()
		allowRepeat, err = strconv.ParseBool(scanner.Text())
		if err != nil {
			panic(err.Error())
		}
	} else { // Not use an interactive program
		// Use arguments and verify if all the arguments are specified
		if len(args) != 3 && len(args) != 5 {
			fileWithoutExt := strings.Split(os.Args[0], ".")[0]
			fmt.Printf("Usage : %s.exe <length> <minimum_number_of_digits> <minimum_number_of_symbols> <allow_uppercase:(false|true)> <allow_repeat:(false|true)>", fileWithoutExt)
			fmt.Println("allow_uppercase and allow_repeat are optional (default is true)")
			os.Exit(2)
		}

		// Convert the arguments
		length, err = strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			panic(err.Error())
		}
		numDigits, err = strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			panic(err.Error())
		}
		numSymbols, err = strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			panic(err.Error())
		}
		if len(args) == 5 {
			allowUpper, err = strconv.ParseBool(args[3])
			if err != nil {
				panic(err.Error())
			}
			allowRepeat, err = strconv.ParseBool(args[4])
			if err != nil {
				panic(err.Error())
			}
		}
	}

	// Generate the password
	gen := NewGenerator(nil)
	pwd, err := gen.Generate(int(length), int(numDigits), int(numSymbols), allowUpper, allowRepeat)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	// Show the generated password
	fmt.Println(pwd)
	if len(args) == 0 {
		print("Please press ENTER to quit the program ...")
		scanner.Scan()
	}
}
