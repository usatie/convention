package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
)

type Package struct {
	Path      string `json:"Path"`
	Version   string `json:"Version"`
	Timestamp string `json:"Timestamp"`
}

func getPackages() ([]Package, error) {
	// Get response
	resp, err := http.Get("https://index.golang.org/index")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Scan response line by line
	scanner := bufio.NewScanner(resp.Body)
	var packages []Package
	for scanner.Scan() {
		// Decode each line
		var pkg Package
		if err := json.Unmarshal(scanner.Bytes(), &pkg); err != nil {
			return nil, err
		}
		packages = append(packages, pkg)
	}
	return packages, nil
}

func main() {

	packages, err := getPackages()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", packages[:5])
}
