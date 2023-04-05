type Response struct {
    ValorConvertido float64 `json:"valorConvertido"`
    SimboloMoeda    string  `json:"simboloMoeda"`
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/exchange/{amount}/{from}/{to}/{rate}", exchangeHandler).Methods("GET")
    log.Fatal(http.ListenAndServe(":8000", r))
}
