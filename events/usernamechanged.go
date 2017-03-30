package events

import (
	"encoding/json"
	"github.com/function61/eventhorizon-exampleapp-go/schema"
	"github.com/function61/eventhorizon-exampleapp-go/transaction"
)

// UserNameChanged {"user_id": "66cad10b", "new_name": "Phyllis Vance", "reason": "..", "ts": "2016-06-06 06:06:06"}
type UserNameChanged struct {
	UserId  string `json:"user_id"`
	Ts      string `json:"ts"`
	NewName string `json:"new_name"`
	Reason  string `json:"reason"`
}

func applyUserNameChanged(tx *transaction.Tx, payload string) error {
	e := UserNameChanged{}

	if err := json.Unmarshal([]byte(payload), &e); err != nil {
		return err
	}

	return tx.Db.WithTransaction(tx.Tx).Update(&schema.User{
		ID:   e.UserId,
		Name: e.NewName,
	})
}

func (u *UserNameChanged) Serialize() string {
	data, _ := json.Marshal(u)
	return "UserNameChanged " + string(data)
}
