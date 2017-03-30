package events

import (
	"encoding/json"
	"github.com/function61/eventhorizon-exampleapp-go/schema"
	"github.com/function61/eventhorizon-exampleapp-go/transaction"
)

// OrderCreated {"id": "c3d2ff02", "user": "4874ce7a64b7", "ts": "2015-03-14 00:00:00"}
type OrderCreated struct {
	Id   string `json:"id"`
	User string `json:"user"`
	Ts   string `json:"ts"`
}

func applyOrderCreated(tx *transaction.Tx, payload string) error {
	e := OrderCreated{}

	if err := json.Unmarshal([]byte(payload), &e); err != nil {
		return err
	}

	order := &schema.Order{
		ID:   e.Id,
		User: e.User,
		Ts:   e.Ts,
	}

	return tx.Db.WithTransaction(tx.Tx).Save(order)
}

func (u *OrderCreated) Serialize() string {
	data, _ := json.Marshal(u)
	return "OrderCreated " + string(data)
}
