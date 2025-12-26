package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Bank struct {
	Name    string
	BinFrom int64
	BinTo   int64
}

func loadBankData(path string) ([]Bank, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Error opening file %q: %w", path, err)
	}
	defer file.Close()

	var banks []Bank
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Split(line, ",")
		lenght := len(fields)
		if lenght != 3 {
			return nil, fmt.Errorf("Row contains an incorrect number of fields (%d): %q", lenght, line)
		}

		name := strings.TrimSpace(fields[0])
		binFromStr := strings.TrimSpace(fields[1])
		binToStr := strings.TrimSpace(fields[2])

		binFrom, err := strconv.ParseInt(binFromStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Incorrect BIN From %q: %w", binFromStr, err)
		}

		binTo, err := strconv.ParseInt(binToStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Incorrect BIN To %q: %w", binToStr, err)
		}

		banks = append(banks, Bank{
			Name:    name,
			BinFrom: binFrom,
			BinTo:   binTo,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading the file: %w", err)
	}

	return banks, nil
}

func extractBIN(cardNumber string) int64 {
	bin, err := strconv.ParseInt(string(cardNumber[:6]), 10, 64)
	if err != nil {
		return 0
	}
	return bin
}

func identifyBank(bin int64, banks []Bank) string {
	for _, bank := range banks {
		if bank.BinFrom <= bin && bin <= bank.BinTo {
			return bank.Name
		}
	}
	return "Unknown Bank"
}

func validateLuhn(cardNumber string) bool {
	sum := 0
	shouldDouble := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')

		if shouldDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		shouldDouble = !shouldDouble
	}

	return sum%10 == 0
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter a credit card number (or press Enter to quit):")
	card, err := reader.ReadString('\n')
	if err != nil {
		return "Error"
	}
	card = strings.TrimSpace(card)
	return card
}

func validateInput(cardNumber string) bool {
	length := len([]rune(cardNumber))
	if length < 13 || length > 19 || cardNumber == "" {
		return false
	}

	for i := 0; i < length; i++ {
		if cardNumber[i] < '0' || cardNumber[i] > '9' {
			return false
		}
	}

	return true
}

func main() {
	fmt.Println("Welcome to the card validation program!")
	banks, err := loadBankData("banks.txt")
	if err != nil {
		log.Fatalf("Failed to load bank data: %v", err)
	}
	for {
		cardNumber := getUserInput()
		if cardNumber == "" {
			fmt.Println("Program is completed")
			break
		} else if !validateInput(cardNumber) {
			fmt.Println("Invalid input. Please enter a valid credit card number.")
			continue
		}
		bin := extractBIN(cardNumber)
		bank := identifyBank(bin, banks)
		fmt.Println("Card number is valid")
		if bank == "Unknown Bank" {
			fmt.Println("The issuer has not been identified")
		} else {
			fmt.Println("Issuing Bank:", bank)
		}
	}
}
