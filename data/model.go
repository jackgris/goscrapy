package data

// This struct will be use for match JSON data
type Product struct {
	MainEntity  MainEntity `json:"mainEntityOfPage"`
	Name        string     `json:"name"`
	Image       string     `json:"image"`
	Description string     `json:"description"`
	Offers      Offers
}

type MainEntity struct {
	Id string `json:"@id"`
}

type Offers struct {
	Price          string `json:"price"`
	Availability   string `json:"availability"`
	InventoryLevel InventoryLevel
}

type InventoryLevel struct {
	Stock string `json:"value"`
}
