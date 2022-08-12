package main

import (
	"fmt"
	"os"

	cli_extension_lib_go "github.com/snyk/cli-extension-lib-go"
)

func main() {
	_, extensionInput, err := cli_extension_lib_go.InitExtension()
	if err != nil {
		fmt.Println("Error initializing extension")
		fmt.Println(err)
		os.Exit(1)
	}

	orgId, err := extensionInput.Command.StringOptionValue("org")
	if err != nil {
		fmt.Println("Invalid input")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Woof in", orgId)
	fmt.Println("Woof token:", extensionInput.Token)
}
