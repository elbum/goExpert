package mydict

import (
	"errors"
	"fmt"
)

// Dictionary type
type Dictionary map[string]string

var errNotFount = errors.New("Not found")
var errWordExists = errors.New("WORD ALREADY EXISTS")
var errNoWordExists = errors.New("WORD DOESN'T EXISTS")

// Search for word
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errNotFount

}

// Add
func (d Dictionary) Add(word, def string) error {
	if _, err := d.Search(word); err == errNotFount {
		d[word] = def
		fmt.Println("WORD ADDED")
	} else if err == nil {
		return errWordExists
	}
	return nil
}

// UPDATE
func (d Dictionary) Update(word, definition string) error {
	switch _, err := d.Search(word); err {
	case errNotFount:
		fmt.Println("THERE'S NO WORD TO UPDATE")
		return errNoWordExists
	case nil:
		d[word] = definition
		fmt.Println("WE FOUND A WORD TO UPDATE")
		return nil
	}
	return errNoWordExists
}

// Delete
func (d Dictionary) Delete(word string) error {
	if _, err := d.Search(word); err == nil {
		fmt.Println("FOUND A WORD TO DELETE")
		delete(d, word)
		return nil
	} else {
		return errNoWordExists
	}
}
