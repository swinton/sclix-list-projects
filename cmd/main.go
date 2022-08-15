package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	cli_extension_lib_go "github.com/snyk/cli-extension-lib-go"
)

func main() {
	// Bootstrap Extension
	_, extensionInput, err := cli_extension_lib_go.InitExtension()
	if err != nil {
		fmt.Println("Error initializing extension")
		fmt.Println(err)
		os.Exit(1)
	}
	t, _ := extensionInput.GetHttpTransport()
	httpClient := &http.Client{
		Transport: t,
	}

	// Check if we have token available for API Requests
	apiToken := extensionInput.Token
	if len(apiToken) == 0 {
		fmt.Println("Missing Snyk API token. Run snyk auth first.")
		fmt.Println(err)
		os.Exit(1)
	}
	apiAuthHeader := "Token " + apiToken

	// Check if user defined an OrgId they want to use
	var orgId string
	orgId, err = extensionInput.Command.StringOptionValue("org")
	if err != nil {
		fmt.Println("Invalid input")
		fmt.Println(err)
		os.Exit(1)
	}

	if len(orgId) == 0 {
		// Query API for default Org
		type OrgsAPIResponse struct {
			Data []struct {
				ID string `json:"id"`
			} `json:"data"`
		}

		req, _ := http.NewRequest("GET", "https://api.snyk.io/rest/orgs/?version=2022-04-06~experimental", nil)
		// Headers
		req.Header.Add("Authorization", apiAuthHeader)
		req.Header.Add("Content-Type", "application/vnd.api+json")
		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Println("Request failure : ", err)
		}
		respBody, _ := ioutil.ReadAll(resp.Body)
		var result OrgsAPIResponse
		if err := json.Unmarshal(respBody, &result); err != nil {
			fmt.Println("Can not unmarshal JSON")
		}

		if len(result.Data) == 0 {
			fmt.Println("Invalid input")
			fmt.Println(err)
			os.Exit(1)
		}
		orgId = result.Data[0].ID
	}

	type ProjectsAPIResponse struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Name string `json:"name"`
			}
		} `json:"data"`
	}
	fmt.Println("Checking projects in Org with ID: " + orgId + "\n")
	req, _ := http.NewRequest("GET", "https://api.snyk.io/rest/orgs/"+orgId+"/projects?version=2022-04-06~experimental", nil)
	req.Header.Add("Authorization", apiAuthHeader)
	req.Header.Add("Content-Type", "application/vnd.api+json")
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Request failure : ", err)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	var result ProjectsAPIResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	if len(result.Data) == 0 {
		fmt.Println("No projects to list")
		os.Exit(0)
	}

	fmt.Printf("Found projects: %d\n\n", len(result.Data))

	for _, project := range result.Data {
		fmt.Println("- " + project.Attributes.Name)
	}
	fmt.Println(result.Data[0].Attributes.Name)
}
