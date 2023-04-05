func exchangeHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    amount, _ := strconv.ParseFloat(vars["amount"], 64)
    from := vars["from"]
    to := vars["to"]
    rate, _ := strconv.ParseFloat(vars["rate"], 64)

    var result float64

    switch {
    case from == "BRL" && to == "USD":
        result = amount / rate
    case from == "USD" && to == "BRL":
        result = amount * rate
    case from == "BRL" && to == "EUR":
        result = amount / rate
    case from == "EUR" && to == "BRL":
        result = amount * rate
    case from == "BTC" && to == "USD":
        result = amount * rate
    case from == "USD" && to == "BTC":
        result = amount / rate
    case from == "BTC" && to == "BRL":
        result = amount * rate
    case from == "BRL" && to == "BTC":
        result = amount / rate
    default:
        http.Error(w, "Conversão não suportada", http.StatusBadRequest)
        return
    }

    response := Response{
        ValorConvertido: result,
        SimboloMoeda:    getSymbol(to),
    }

    saveConversion(amount, from, to, rate, result)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
