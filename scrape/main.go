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

func getPackage(pkg Package) {
	// exec packageでgo getコマンドを叩く
	// まずは手でそのコマンドを叩いてみる
	// 1. Create temp directory
	// 2. Cd to temp directory ( exec packageでcdできる）
	// 3. go mod initで初期化（最後の引数は適当）
	// 4. indexをparseしたpackage名を使ってgo get {module_path}/...
	// 5. go vet -vettool {module_path}/...
	/*
		// Get response
		path := fmt.Sprintf("https://proxy.golang.org/%s/@v/%s.info", pkg.Path, pkg.Version)
		resp, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Scan response line by line
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			// Decode each line
			fmt.Println(scanner.Text())
		}
		return nil, nil
	*/
}

func main() {

	packages, err := getPackages()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", packages[:5])
	getPackage(packages[0])
}
