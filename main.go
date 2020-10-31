package main

import (
	"github.com/jackgris/goscrapy/wholesalers"
)

func main(){

	GetConfig()
	wholesalers.GetData(
		dburi,
		dbuser,
		dbpass,
		GetWholesalersData(),
	)
}
