package types

// OrderType Define a new type for the enum
type OrderType int

// Define the constants using iota
const (
	Espp OrderType = iota
	Rsu
)
