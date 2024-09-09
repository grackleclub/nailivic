package main

type name struct {
	id    int
	value string
}

// don't do this. start with the database model and migrations in sqlc
// https://github.com/Turkosaurus/nailivic/blob/main/database.py
var names = map[string]name{
	"dama": {1, "La Dama"},
	"mano": {2, "El Mano"},
}

type color struct {
	id    int
	value string
}

// const red color = "red"
// const black color = "black"
// const green color = "green"

type size string

const small size = "small"
const medium size = "medium"
const large size = "large"

// type Inventory struct {
// 	Items []Item
// 	Parts []Part
// }

type Item struct {
	Name name
	A    Part
	B    Part
	C    Part
}

func (i *Item) assemble(a Part, b Part, c Part) {}
func (i *Item) disassemble()                    {}
func (i *Item) sku()                            {}

// type Part struct {
// 	ItemName name
// 	Letter   string // (a, b, c)
// 	Color    color
// 	Size     size
// }

// func newPart(name name, color color) Part {
// 	return Part{
// 		ItemName: name,
// 		Color:    color,
// 	}
// }

func (p *Part) create(color string) {}
