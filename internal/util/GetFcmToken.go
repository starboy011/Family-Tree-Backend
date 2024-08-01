package util

import (
	"database/sql"
	"log"
)

type Tokens struct {
	Tokens []string
}

func GetFcmToken(db *sql.DB) Tokens {
	var tokens Tokens
	query := `SELECT token FROM fcmtoken`

	// Execute query
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return Tokens{}
	}
	defer rows.Close()

	// Iterate through rows
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			log.Printf("Error scanning row: %v", err)
			return Tokens{}
		}
		tokens.Tokens = append(tokens.Tokens, token)
	}

	// Check for any error encountered during iteration
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return Tokens{}
	}

	return tokens
}
