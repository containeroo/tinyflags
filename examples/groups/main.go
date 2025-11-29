package main

import (
	"fmt"
	"os"

	"github.com/containeroo/tinyflags"
)

func main() {
	fs := tinyflags.NewFlagSet("groups", tinyflags.ExitOnError)

	api := fs.GetAllOrNoneGroup("api")
	apiFlag := fs.Bool("api", false, "Enable API").AllOrNone("api").Value()

	db := fs.GetAllOrNoneGroup("db")
	dbFlag := fs.Bool("db", false, "Enable database").AllOrNone("db").Value()

	cache := fs.GetAllOrNoneGroup("cache")
	cacheFlag := fs.Bool("cache", false, "Enable cache").AllOrNone("cache").Value()

	fs.GetOneOfGroup("stack").
		Required().
		AddGroup(api).
		AddGroup(db).
		AddGroup(cache)

	// ensure cache and db settings are treated as a block
	fs.AttachGroupToAllOrNone("db", "cache")

	if err := fs.Parse(nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("api enabled:", *apiFlag)
	fmt.Println("db enabled:", *dbFlag)
	fmt.Println("cache enabled:", *cacheFlag)
}
