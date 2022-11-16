package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	vault "github.com/hashicorp/vault/api"
	"github.com/wwalexander/spelunk/spelunk"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"Usage: %s [flags] [path ...]\n",
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.String("path", "", "True if the pathname being examined matches `pattern`")
	flag.String("name", "", "True if the last component of the pathname being examined matches `pattern`")
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	roots := flag.Args()
	var filter struct {
		Name *string
		Path *string
	}
	flag.Visit(func(f *flag.Flag) {
		value := f.Value.String()
		switch f.Name {
		case "path":
			filter.Path = &value
		case "name":
			filter.Name = &value
		}
	})
	client, err := vault.NewClient(nil)
	if err != nil {
		log.Fatal(err)
	}
	code := 0
	fn := func(name string, data map[string]interface{}, err error) error {
		if err != nil {
			log.Printf("%s: %v", name, err)
			code = 1
		}
		if pattern := filter.Path; pattern != nil {
			matched, err := path.Match(*pattern, name)
			if err != nil {
				return err
			}
			if !matched {
				return nil
			}
		}
		if pattern := filter.Name; pattern != nil {
			matched, err := path.Match(*pattern, path.Base(name))
			if err != nil {
				return err
			}
			if !matched {
				return nil
			}
		}
		fmt.Println(name)
		return nil
	}
	for _, root := range roots {
		if err := spelunk.Walk(client.Logical(), root, fn); err != nil {
			log.Print(err)
		}
	}
	os.Exit(code)
}
