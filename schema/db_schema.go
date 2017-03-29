package schema

// these are the structs used by Storm ORM to read / write into the Bolt database.

type User struct {
	ID      string
	Name    string
	Company string
}

type Company struct {
	ID   string
	Name string
}

type LineItem struct {
	Product string
	Amount  int
}

type Order struct {
	ID        string
	User      string
	Ts        string
	LineItems []LineItem
}
