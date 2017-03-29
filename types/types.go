package types

// this is so the types can be reused elsewhere

type OrderPlacementItem struct {
	Product string
	Amount  int
}

type OrderPlacement struct {
	User      string
	LineItems []OrderPlacementItem
}
