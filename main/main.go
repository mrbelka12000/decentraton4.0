package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano())) // локальный генератор

func main() {
	os.MkdirAll("case 1", 0755)

	generateClients()
	generateTransactions()
	generateTransfers()

	fmt.Println("Test data generated in 'case 1' folder")
}

func generateClients() {
	file, _ := os.Create("case 1/clients.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"client_code", "name", "status", "age", "city", "avg_monthly_balance_KZT"})

	names := []string{"Айгерим", "Асем", "Данияр", "Нурлан", "Алия"}
	cities := []string{"Алматы", "Нур-Султан", "Шымкент", "Актобе", "Караганда"}

	for i := 1; i <= 5; i++ {
		balance := r.Intn(5000000) + 50000
		writer.Write([]string{
			strconv.Itoa(i),
			names[i-1],
			"Зарплатный клиент",
			strconv.Itoa(r.Intn(40) + 20),
			cities[r.Intn(len(cities))],
			strconv.Itoa(balance),
		})
	}
}

func generateTransactions() {
	file, _ := os.Create("case 1/transactions.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Заголовок
	writer.Write([]string{"client_code", "name", "product", "status", "city", "date", "category", "amount", "currency"})

	categories := []string{
		"Такси", "Едим дома", "Смотрим дома", "Играем дома",
		"Кафе и рестораны", "Продукты", "Одежда", "Путешествия",
		"Косметика и Парфюмерия", "Ювелирные украшения",
	}

	names := []string{"Айгерим", "Асем", "Данияр", "Нурлан", "Алия"}

	// Генерируем транзакции для каждого клиента за последние 90 дней
	for clientCode := 1; clientCode <= 5; clientCode++ {
		transactionsCount := rand.Intn(50) + 20 // 20-70 транзакций

		for i := 0; i < transactionsCount; i++ {
			// Случайная дата за последние 90 дней
			daysAgo := rand.Intn(90)
			date := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")

			category := categories[rand.Intn(len(categories))]
			amount := generateAmountForCategory(category)

			writer.Write([]string{
				strconv.Itoa(clientCode),
				names[clientCode-1],
				"Карта",
				"Активный",
				"Алматы",
				date,
				category,
				fmt.Sprintf("%.2f", amount),
				"KZT",
			})
		}
	}
}

func generateTransfers() {
	file, _ := os.Create("case 1/transfers.csv")
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Заголовок
	writer.Write([]string{"client_code", "name", "product", "status", "city", "date", "type", "direction", "amount", "currency"})

	transferTypes := []string{"salary_in", "p2p_out", "utilities_out", "card_out", "atm_withdrawal"}
	names := []string{"Айгерим", "Асем", "Данияр", "Нурлан", "Алия"}

	for clientCode := 1; clientCode <= 5; clientCode++ {
		transfersCount := rand.Intn(30) + 10 // 10-40 переводов

		for i := 0; i < transfersCount; i++ {
			daysAgo := rand.Intn(90)
			date := time.Now().AddDate(0, 0, -daysAgo).Format("2006-01-02")

			transferType := transferTypes[rand.Intn(len(transferTypes))]
			direction := "out"
			amount := float64(rand.Intn(100000) + 1000) // 1k-100k

			if transferType == "salary_in" {
				direction = "in"
				amount = float64(rand.Intn(500000) + 200000) // зарплата 200k-700k
			}

			writer.Write([]string{
				strconv.Itoa(clientCode),
				names[clientCode-1],
				"Счет",
				"Активный",
				"Алматы",
				date,
				transferType,
				direction,
				fmt.Sprintf("%.2f", amount),
				"KZT",
			})
		}
	}
}

func generateAmountForCategory(category string) float64 {
	switch category {
	case "Такси":
		return float64(rand.Intn(5000) + 500) // 500-5500
	case "Едим дома", "Смотрим дома", "Играем дома":
		return float64(rand.Intn(3000) + 200) // 200-3200
	case "Путешествия":
		return float64(rand.Intn(200000) + 10000) // 10k-210k
	case "Ювелирные украшения":
		return float64(rand.Intn(500000) + 20000) // 20k-520k
	case "Кафе и рестораны":
		return float64(rand.Intn(15000) + 1000) // 1k-16k
	default:
		return float64(rand.Intn(20000) + 500) // 500-20500
	}
}
