package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	cli_extension_lib_go "github.com/snyk/cli-extension-lib-go"
)

type ProjectsAPIResponse struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Name string `json:"name"`
		}
	} `json:"data"`
}

func mainE(extensionInput *cli_extension_lib_go.ExtensionInput) (err error) {
	t, err := extensionInput.GetHttpTransport()
	if err != nil {
		return fmt.Errorf("Failed to get Transport object. %v", err)
	}

	httpClient := &http.Client{
		Transport: t,
	}

	// Check if we have token available for API Requests
	apiToken := extensionInput.Token
	if len(apiToken) == 0 {
		return fmt.Errorf("Missing Snyk API token. Run snyk auth first.")
	}
	apiAuthHeader := "Token " + apiToken

	// Check if more emoji requested
	moreEmoji, err := extensionInput.Command.BoolOptionValue("more-emoji")
	if err != nil {
		return fmt.Errorf("Invalid input %v", err)
	}

	if moreEmoji {
		fmt.Println("🚀")
	}

	// Check if user defined an OrgId they want to use
	var orgId string
	orgId, err = extensionInput.Command.StringOptionValue("org")
	if err != nil {
		return fmt.Errorf("Invalid input %v", err)
	}

	if len(orgId) == 0 {
		// Query API for default Org
		type OrgsAPIResponse struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}

		req, err := http.NewRequest("GET", "https://api.snyk.io/rest/orgs/?version=2022-04-06~experimental", nil)
		if err != nil {
			return err
		}

		// Headers
		req.Header.Add("Authorization", apiAuthHeader)
		req.Header.Add("Content-Type", "application/vnd.api+json")
		resp, err := httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("Request failure: %v", err)
		}
		respBody, _ := ioutil.ReadAll(resp.Body)
		var result OrgsAPIResponse
		if err := json.Unmarshal(respBody, &result); err != nil {
			return fmt.Errorf("Can not unmarshal JSON. %v", err)
		}

		if len(result.Data) == 0 {
			return err
		}
		orgId = result.Data[0].ID
	}

	fmt.Println("Checking projects in Org with ID: " + orgId + "\n")
	req, err := http.NewRequest("GET", "https://api.snyk.io/rest/orgs/"+orgId+"/projects?version=2022-04-06~experimental", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", apiAuthHeader)
	req.Header.Add("Content-Type", "application/vnd.api+json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Request failure : %v", err)
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	var result ProjectsAPIResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("Can not unmarshal JSON. %v", err)
	}

	if len(result.Data) == 0 {
		return fmt.Errorf("No projects to list")
	}

	fmt.Printf("Found projects: %d\n\n", len(result.Data))

	for _, project := range result.Data {
		fmt.Println("- " + project.Attributes.Name)
	}
	fmt.Println(result.Data[0].Attributes.Name)

	return nil
}

func main() {
	exitCode := 0

	// Bootstrap Extension
	_, extensionInput, err := cli_extension_lib_go.InitExtension()
	if err != nil {
		fmt.Println("[list extensions] Error:", err)
		os.Exit(1)
	}

	if extensionInput.Debug {
		fmt.Println("[list extensions] Starting")
	}

	err = mainE(extensionInput)
	if err != nil {
		fmt.Println("[list extensions] Error:", err)
		exitCode = 1
	}

	if extensionInput.Debug {
		fmt.Printf("[list extensions] Exiting with %d \n", exitCode)
	}

	os.Exit(exitCode)
}
