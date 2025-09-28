package models

type Client struct {
    AlgoID string
    Nome   string
}

type Order struct {
    OrderID  string
    Product  string
    Amount   int
    ClientID string
}
