package main

import "fmt"

// example of composite types a template might receive page
// and parse page.nav.bar, page.head, page.footer, etc
type page struct {
	nav    navbar
	head   head
	body   body
	footer footer
}
type navbar struct{}
type head struct{}
type footer struct{}
type body struct{}

func newPage() (page, error) {
	// return page{
	// 	nav:    navbar{},
	// 	head:   head{
	// 		Title: "nailivic",

	// 	},
	// 	footer: footer{
	// 		Year: 2021,
	// 	},
	// }, nil
	return page{}, fmt.Errorf("not implemented")
}
