package main

import "fmt"

// example of composite types a template might receive page
// and parse page.nav.bar, page.head, page.footer, etc
type page struct {
	// Nav    navbar
	Head      head
	Parts     []Part
	Inventory []Inventory
	Footer    footer
}

type head struct {
	Title       string
	Stylesheets []string
}

type Part struct {
	Piece string
	Color string
	Count int
}

type Inventory struct {
	Item    string
	Color_A string
	Color_B string
	Color_C string
	Size    string
	Count   int
}

//	type navbar struct {
//		Items []string
//	}
type footer struct {
	Year string
}

func getNewPage(name string) (page, error) {
	switch name {
	case "login":
		return page{}, fmt.Errorf("not implemented")
	case "parts":
		newPage := page{
			Head: head{
				Title: "Parts",
				Stylesheets: []string{
					"/static/css/zero.css",
					"/static/css/style.css",
				},
			},
			Parts: []Part{
				{
					Piece: "foo",
					Color: "bar",
					Count: 100,
				},
				{
					Piece: "baz",
					Color: "qux",
					Count: 200,
				},
				{
					Piece: "quux",
					Color: "corge",
					Count: 300,
				},
			},
			Footer: footer{
				Year: "9876",
			},
		}
		return newPage, nil
	case "inventory":
		newPage := page{
			Head: head{
				Title: "Inventory",
				Stylesheets: []string{
					"/static/css/zero.css",
					"/static/css/style.css",
				},
			},
			Inventory: []Inventory{
				{
					Item:    "foo",
					Color_A: "bar",
					Color_B: "baz",
					Color_C: "qux",
					Size:    "small",
					Count:   100,
				},
				{
					Item:    "baz",
					Color_A: "qux",
					Color_B: "quux",
					Color_C: "corge",
					Size:    "medium",
					Count:   200,
				},
				{
					Item:    "quux",
					Color_A: "corge",
					Color_B: "grault",
					Color_C: "garply",
					Size:    "large",
					Count:   300,
				},
				{
					Item:    "garply",
					Color_A: "waldo",
					Color_B: "fred",
					Color_C: "plugh",
					Size:    "small",
					Count:   400,
				},
			},
			Footer: footer{
				Year: "9876",
			},
		}
		return newPage, nil
	case "index":
		newPage := page{
			Head: head{
				Title: "Nailivic",
				Stylesheets: []string{
					"/static/css/zero.css",
					"/static/css/style.css",
				},
			},
			Footer: footer{
				Year: "2024",
			},
		}
		return newPage, nil
	case "dash":
		newPage := page{
			Head: head{
				Title: "Nailivic Dashboard",
				Stylesheets: []string{
					"/static/css/zero.css",
					"/static/css/style.css",
				},
			},
			Footer: footer{
				Year: "2024",
			},
		}
		return newPage, nil
	case "special":
		newPage := page{
			Head: head{
				Title: "Nailivic",
				Stylesheets: []string{
					"/static/css/zero.css",
					"/static/css/style.css",
				},
			},
			Footer: footer{
				Year: "2024",
			},
		}
		return newPage, nil
	default:
		return page{}, fmt.Errorf("invalid page name: %s", name)
	}
}

// func doStuffWrapper() {
// 	page, err := getNewPage("index")
// 	if err != nil {
// 		fmt.Println("error getting new page:", err)
// 		return
// 	}
// 	fmt.Println("new page:", page)
// }
