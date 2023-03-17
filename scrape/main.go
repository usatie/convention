package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
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
		//fmt.Println(err)
		//fmt.Printf("stderr: %s\n", stderr.String())
		return err
	}
	//fmt.Printf("stdout: %s\n", stdout.String())
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
		//fmt.Println(err)
		//fmt.Printf("stderr: %s\n", stderr.String())
		return err
	}
	//fmt.Printf("stdout: %s\n", stdout.String())
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
	writeResult(pkg, stderr)
	if err != nil {
		//fmt.Println(err)
		fmt.Printf("stderr: %s\n", stderr.String())
		return err
	}
	fmt.Printf("stdout: %s\n", stdout.String())
	return nil
}

const resultDir string = "results"

func writeResult(pkg Package, stderr strings.Builder) {
	s := stderr.String()
	if s == "" {
		s = "# " + pkg.Path
	}
	filename := path.Join(resultDir, "results.txt")
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(s); err != nil {
		log.Fatal(err)
	}
}

func analyze(pkg Package) {
	// Create result dir
	err := os.Mkdir(resultDir, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	// Create temp dir
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
	// 1. Create temp directory, Create result directory
	// 2. Cd to temp directory ( exec packageでcdできる）
	// 3. go mod initで初期化（最後の引数は適当）
	// 4. indexをparseしたpackage名を使ってgo get {module_path}/...
	// 5. go vet -vettool {module_path}/...
	// 6. 結果をresultに集計
}

func main() {

	packages, err := getPackages()
	if err != nil {
		panic(err)
	}
	fmt.Println("Number of packages: ", len(packages))
	for _, pkg := range packages {
		analyze(pkg)
	}
}
