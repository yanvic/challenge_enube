package usecase

import (
    "strconv"

    "github.com/xuri/excelize/v2"
    "challenge_enube/internal/models"
	"challenge_enube/internal/repository"
)

type ImportUseCase struct {
    clientRepo *repository.ClientRepository
    orderRepo  *repository.OrderRepository
}

func NewImportUseCase(cRepo *repository.ClientRepository, oRepo *repository.OrderRepository) *ImportUseCase {
    return &ImportUseCase{clientRepo: cRepo, orderRepo: oRepo}
}

func (uc *ImportUseCase) ImportFromExcel(path string) error {
    f, err := excelize.OpenFile(path)
    if err != nil {
        return err
    }
    defer f.Close()

    rows, err := f.GetRows("Planilha1")
    if err != nil {
        return err
    }

    if len(rows) < 2 {
        return nil
    }

    var clients []models.Client
    var orders []models.Order

    for _, row := range rows[1:] {
        if len(row) < 5 {
            continue
        }

        clientID := row[0]
        clientName := row[1]
        orderID := row[2]
        product := row[3]
        amount, _ := strconv.Atoi(row[4])

        clients = append(clients, models.Client{
            AlgoID: clientID,
            Nome:   clientName,
        })

        orders = append(orders, models.Order{
            OrderID:  orderID,
            Product:  product,
            Amount:   amount,
            ClientID: clientID,
        })
    }

    if err := uc.clientRepo.SaveBatch(clients); err != nil {
        return err
    }

    if err := uc.orderRepo.SaveBatch(orders); err != nil {
        return err
    }

    return nil
}
