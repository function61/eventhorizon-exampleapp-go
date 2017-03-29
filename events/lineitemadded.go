package events

import (
	"encoding/json"
	"github.com/function61/pyramid-exampleapp-go/schema"
	"github.com/function61/pyramid-exampleapp-go/transaction"
)

// LineItemAdded {"order": "66cad10b", "product": "Regular paper", "amount": 3, "ts": "2016-06-06 06:06:06"}
type LineItemAdded struct {
	Order   string `json:"order"`
	Product string `json:"product"`
	Amount  int    `json:"amount"`
	Ts      string `json:"ts"`
}

func applyLineItemAdded(tx *transaction.Tx, payload string) error {
	e := LineItemAdded{}

	if err := json.Unmarshal([]byte(payload), &e); err != nil {
		return err
	}

	db := tx.Db.WithTransaction(tx.Tx)

	var order schema.Order
	if err := db.One("ID", e.Order, &order); err != nil {
		// log.Printf("ffffuck this shouldnt happen")
		return nil
	}

	lineItem := schema.LineItem{
		Product: e.Product,
		Amount:  e.Amount,
	}

	order.LineItems = append(order.LineItems, lineItem)

	return db.Update(&order)
}

func (u *LineItemAdded) Serialize() string {
	data, _ := json.Marshal(u)
	return "LineItemAdded " + string(data)
}
