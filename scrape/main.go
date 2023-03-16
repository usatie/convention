package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
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

func goModInit(dir string) error {
	cmd := exec.Command("go", "mod", "init", "hoge")
	fmt.Println(cmd.String())
	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("stderr: %s\n", stderr.String())
		return err
	}
	fmt.Printf("stdout: %s\n", stdout.String())
	return nil
}

func goGet(pkg Package, dir string) error {

	cmd := exec.Command("go", "get", pkg.Path+"/...")
	fmt.Println(cmd.String())
	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("stderr: %s\n", stderr.String())
		return err
	}
	fmt.Printf("stdout: %s\n", stdout.String())
	return nil
}

func goVet(pkg Package, dir string) error {
	// go vet -vettool=/Users/shunusami/Desktop/gopher-intern/convention/convention software.sslmate.com/src/go-pkcs12/...
	executablePath := "/Users/shunusami/Desktop/gopher-intern/convention/convention"
	cmd := exec.Command("go", "vet", "-vettool="+executablePath, pkg.Path+"/...")
	fmt.Println(cmd.String())
	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("stderr: %s\n", stderr.String())
		return err
	}
	fmt.Printf("stdout: %s\n", stdout.String())
	return nil
}

func analyze(pkg Package) {
	dir, err := os.MkdirTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up
	if err = goModInit(dir); err != nil {
		return
	}
	if err = goGet(pkg, dir); err != nil {
		return
	}
	if err = goVet(pkg, dir); err != nil {
		return
	}

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
	for _, pkg := range packages[:5] {
		analyze(pkg)
	}
}
