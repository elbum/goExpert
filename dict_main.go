package main

import (
	"fmt"

	"github.com/elbum/goExpert/banking"
	"github.com/elbum/goExpert/mydict"
)

func main() {
	// 생성자로 만들면 포인터로 컨트롤 해야함 (주소를 리턴하니까.)
	account := banking.NewAccount("bums")
	account.Deposit(10000)
	if err := account.Withraw(1000000); err != nil {
		fmt.Println(err)
	}

	fmt.Println(account.Balance(), account.Owner())
	fmt.Println(account) // using account string()
	// fmt.Printf("%+v\n", account)
	account.ChangeOwner("Woni")
	fmt.Println(account.Balance(), account.Owner())

	fmt.Printf("\n\n\n\n\n")

	dictionary := mydict.Dictionary{"first": "FIRST WORD"}
	dictionary["hello"] = "hi"
	fmt.Println(dictionary)

	definition, err := dictionary.Search("first")
	if err != nil {
		fmt.Println("NO WORD")
	} else {
		fmt.Println("We GOT ", definition)
	}

	dictionary.Add("third", "IAM THIRD")
	fmt.Println(dictionary)

	if definition, err := dictionary.Search("third"); err != nil {
		fmt.Println("NO WORD")
	} else {
		fmt.Println("We GOT ", definition)
	}

	word := "hello"
	if err := dictionary.Update(word, "UPDATEHELLO"); err != nil {
		fmt.Println(err)
	}

	if err := dictionary.Update("hoho", "UPDATEHELLO"); err != nil {
		fmt.Println(err)
	}

	fmt.Println(dictionary)

	if err := dictionary.Delete("first"); err != nil {
		fmt.Println(err)
	}

	fmt.Println(dictionary)
}
