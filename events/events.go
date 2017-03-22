package events

import (
	"github.com/function61/pyramid-exampleapp-go/transaction"
)

var EventNameToApplyFn = map[string]func(*transaction.Tx, string) error{
	"CompanyCreated":  applyCompanyCreated,
	"UserCreated":     applyUserCreated,
	"UserNameChanged": applyUserNameChanged,
}
