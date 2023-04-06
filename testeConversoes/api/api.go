package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ericlagergren/decimal"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Conversion struct {
	ID              int64           `json:"id"`
	FromCurrency    string          `json:"fromCurrency"`
	ToCurrency      string          `json:"toCurrency"`
	Rate            decimal.Decimal `json:"rate"`
	Amount          decimal.Decimal `json:"amount"`
	ConvertedAmount decimal.Decimal `json:"convertedAmount"`
	CreatedAt       time.Time       `json:"createdAt"`
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
	rate, _ := decimal.NewFromString(vars["rate"])

	// Calcular a conversão de moedas
	var convertedAmount decimal.Decimal
	if from == "BRL" && to == "USD" {
		convertedAmount = decimal.NewFromFloat(amount).Mul(rate)
	} else if from == "USD" && to == "BRL" {
		convertedAmount = decimal.NewFromFloat(amount).Div(rate)
	} else if from == "BRL" && to == "EUR" {
		convertedAmount = decimal.NewFromFloat(amount).Mul(rate)
	} else if from == "EUR" && to == "BRL" {
		convertedAmount = decimal.NewFromFloat(amount).Div(rate)
	} else if from == "BTC" && to == "USD" {
		convertedAmount = decimal.NewFromFloat(amount).Mul(rate)
	} else if from == "BTC" && to == "BRL" {
		convertedAmount = decimal.NewFromFloat(amount).Mul(rate)
	} else {
		http.Error(w, "Invalid conversion", http.StatusBadRequest)
		return
	}
}

// Criar uma nova conversão
func conv(amount decimal.Decimal, from string, to string, rate decimal.Decimal) (decimal.Decimal, string) {
	// Define a quantidade de casas decimais para o resultado
	decimal.DivisionPrecision = 2

	// Realiza a conversão para a moeda desejada
	var result decimal.Decimal
	if from == "BRL" {
		result = amount.Div(rate)
	} else if to == "BRL" {
		result = amount.Mul(rate)
	} else {
		result = amount
	}

	// Retorna o valor convertido e o símbolo da moeda
	var symbol string
	if to == "USD" || from == "USD" {
		symbol = "$"
	} else if to == "EUR" || from == "EUR" {
		symbol = "€"
	} else if to == "BTC" || from == "BTC" {
		symbol = "฿"
	} else {
		symbol = "R$"
	}
	return result, symbol
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
