package events

import (
	"github.com/function61/eventhorizon-exampleapp-go/transaction"
)

var EventNameToApplyFn = map[string]func(*transaction.Tx, string) error{
	"CompanyCreated":  applyCompanyCreated,
	"UserCreated":     applyUserCreated,
	"UserNameChanged": applyUserNameChanged,
	"LineItemAdded":   applyLineItemAdded,
	"OrderCreated":    applyOrderCreated,
}
