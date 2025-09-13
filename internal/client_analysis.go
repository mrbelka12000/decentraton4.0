package internal

import (
	"database/sql"
	"fmt"
)

type ClientAnalysis struct {
	ClientCode        int                `json:"client_code"`
	Name              string             `json:"name"`
	AvgMonthlyBalance float64            `json:"avg_monthly_balance"`
	SpendByCategory   map[string]float64 `json:"spend_by_category"`
	CategoryShares    map[string]float64 `json:"category_shares"`
	TotalSpend        float64            `json:"total_spend"`
	OnlineSpend       float64            `json:"online_spend"`
	OnlineShare       float64            `json:"online_share"`
	TaxiSpend         float64            `json:"taxi_spend"`
	TaxiCount         int                `json:"taxi_count"`
	Inflow            float64            `json:"inflow"`
	Outflow           float64            `json:"outflow"`
	CashGap           float64            `json:"cash_gap"`
}

func CalculateClientAnalysis(db *sql.DB, clientCode int) (*ClientAnalysis, error) {
	analysis := &ClientAnalysis{
		ClientCode:      clientCode,
		SpendByCategory: make(map[string]float64),
		CategoryShares:  make(map[string]float64),
	}

	err := db.QueryRow(`
		SELECT name, avg_monthly_balance_KZT 
		FROM clients 
		WHERE client_code = $1
	`, clientCode).Scan(&analysis.Name, &analysis.AvgMonthlyBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to get client info: %w", err)
	}

	// Траты по категориям за последние 3 месяца
	rows, err := db.Query(`
		SELECT category, SUM(amount) as total_amount, COUNT(*) as count
		FROM transactions 
		WHERE client_code = $1 
		  AND date >= (CURRENT_DATE - INTERVAL '3 months')
		GROUP BY category
	`, clientCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get category spending: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var amount float64
		var count int

		if err := rows.Scan(&category, &amount, &count); err != nil {
			continue
		}

		analysis.SpendByCategory[category] = amount
		analysis.TotalSpend += amount

		if category == "Такси" {
			analysis.TaxiSpend = amount
			analysis.TaxiCount = count
		}

		if isOnlineCategory(category) {
			analysis.OnlineSpend += amount
		}
	}

	// Доли категорий
	if analysis.TotalSpend > 0 {
		for category, amount := range analysis.SpendByCategory {
			analysis.CategoryShares[category] = amount / analysis.TotalSpend
		}
		analysis.OnlineShare = analysis.OnlineSpend / analysis.TotalSpend
	}

	// Inflow / outflow
	err = db.QueryRow(`
		SELECT 
			COALESCE(SUM(CASE WHEN direction = 'in' THEN amount ELSE 0 END), 0) as inflow,
			COALESCE(SUM(CASE WHEN direction = 'out' THEN amount ELSE 0 END), 0) as outflow
		FROM transfers 
		WHERE client_code = $1 
		  AND date >= (CURRENT_DATE - INTERVAL '3 months')
	`, clientCode).Scan(&analysis.Inflow, &analysis.Outflow)
	if err != nil {
		return nil, fmt.Errorf("failed to get cash flows: %w", err)
	}

	// Кассовый разрыв
	if analysis.Inflow > 0 {
		analysis.CashGap = analysis.Outflow / analysis.Inflow
	}

	return analysis, nil
}

func isOnlineCategory(category string) bool {
	onlineCategories := map[string]bool{
		"Едим дома":    true,
		"Смотрим дома": true,
		"Играем дома":  true,
	}
	return onlineCategories[category]
}

func GetAllClientAnalysis(db *sql.DB) (map[int]*ClientAnalysis, error) {
	rows, err := db.Query(`SELECT client_code FROM clients`)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}
	defer rows.Close()

	allFeatures := make(map[int]*ClientAnalysis)

	for rows.Next() {
		var clientCode int
		if err := rows.Scan(&clientCode); err != nil {
			continue
		}

		features, err := CalculateClientAnalysis(db, clientCode)
		if err != nil {
			fmt.Printf("Warning: failed to calculate features for client %d: %v\n", clientCode, err)
			continue
		}

		allFeatures[clientCode] = features
	}

	return allFeatures, nil
}

func (cf *ClientAnalysis) String() string {
	return fmt.Sprintf("Code %d: %s, Balance=%.0f₸, Spend=%.0f₸, Online=%.1f%%, Gap=%.2f",
		cf.ClientCode, cf.Name, cf.AvgMonthlyBalance, cf.TotalSpend, cf.OnlineShare*100, cf.CashGap)
}
