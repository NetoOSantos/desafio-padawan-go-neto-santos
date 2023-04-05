package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
    "github.com/ericlagergren/decimal"
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
    log.Fatal(http.ListenAndServe(":8000", router))
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

    // Criar uma nova conversão
    conv := Conversion{
        FromCurrency:    from,
        ToCurrency:      to,
        Rate:
