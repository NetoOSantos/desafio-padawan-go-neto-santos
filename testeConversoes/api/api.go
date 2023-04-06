package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Conversion struct {
	ID              int64     `json:"id"`
	FromCurrency    string    `json:"fromCurrency"`
	ToCurrency      string    `json:"toCurrency"`
	Rate            float64   `json:"rate"`
	Amount          float64   `json:"amount"`
	ConvertedAmount float64   `json:"convertedAmount"`
	CreatedAt       time.Time `json:"createdAt"`
}

func main() {
	// Configurar o roteador da API
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/exchange/{amount}/{from}/{to}/{rate}", handleExchange)

	// Configurar o banco de dados
	db, err := sql.Open("mysql", "user:password@/database")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verificar a conexão com o banco de dados
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Iniciar o servidor HTTP
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleExchange(w http.ResponseWriter, r *http.Request) {
	// Obter os parâmetros da URL
	vars := mux.Vars(r)
	amount, _ := strconv.ParseFloat(vars["amount"], 64)
	from := vars["from"]
	to := vars["to"]
	rate, _ := strconv.ParseFloat(vars["rate"], 64)

	// Calcular a conversão de moedas
	var convertedAmount float64
	if from == "BRL" && to == "USD" {
		convertedAmount = amount * rate
	} else if from == "USD" && to == "BRL" {
		convertedAmount = amount / rate
	} else if from == "BRL" && to == "EUR" {
		convertedAmount = amount * rate
	} else if from == "EUR" && to == "BRL" {
		convertedAmount = amount / rate
	} else if from == "BTC" && to == "USD" {
		convertedAmount = amount * rate
	} else if from == "BTC" && to == "BRL" {
		convertedAmount = amount * rate
	} else {
		http.Error(w, "Invalid conversion", http.StatusBadRequest)
		return
	}

	// Criar um objeto Conversion com os dados da conversão
	conv := Conversion{
		FromCurrency:    from,
		ToCurrency:      to,
		Rate:            rate,
		Amount:          amount,
		ConvertedAmount: convertedAmount,
		CreatedAt:       time.Now(),
	}

	// Salvar a conversão no banco de dados
func saveConversion(db *sql.DB, c Conversion) error {
	query := `INSERT INTO conversions(from_currency, to_currency, rate, amount, converted_amount, created_at)
	VALUES (?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, c.FromCurrency, c.ToCurrency, c.Rate, c.Amount, c.ConvertedAmount, c.CreatedAt)
	return err
}

// Recuperar as conversões do banco de dados
func getConversions(db *sql.DB) ([]Conversion, error) {
	query := "SELECT id, from_currency, to_currency, rate, amount, converted_amount, created_at FROM conversions ORDER BY id DESC"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversions := []Conversion{}
	for rows.Next() {
		c := Conversion{}
		err := rows.Scan(&c.ID, &c.FromCurrency, &c.ToCurrency, &c.Rate, &c.Amount, &c.ConvertedAmount, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		conversions = append(conversions, c)
	}
	return conversions, nil
}

/*          Desenvolva uma REST API utilizando a linguagem GO
            e orientação a objetos que faça conversão de moedas.

            Especifícações:

            A URL da requisição deve seguir o seguinte formato:
            http://localhost:8000/exchange/{amount}/{from}/{to}/{rate}
            http://localhost:8000/exchange/10/BRL/USD/4.50
            A resposta deve seguir o seguinte formato:
            {
            "valorConvertido": 45,
            "simboloMoeda": "$"
            }

            * Conversões:
            * De Real para Dólar;
            * De Dólar para Real;
            * De Real para Euro;
            * De Euro para Real;
            * De BTC para Dolar;
            * De BTC para Real;

            * Salvar os dados no banco de dados:
            * criar uma rotina para salvar o dados para consultas futuras
*/
