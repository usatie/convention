//line /a.go:1
package a

func doSomething() error {
	return nil
}

func makePointer() *struct{} {
	return nil
}

func f() {
	// Good:
	if err := doSomething(); err != nil {
		return
	}

	// Good:
	if err := doSomething(); err == nil { // if NO error
		return
	}

	// Good:
	if err := doSomething(); nil == err { // if NO error
		return
	}

	// Bad: missing comment
	if err := doSomething(); err == nil {
		return
	}

	// Bad: missing comment
	if err := doSomething(); nil == err {
		return
	}

	// Good: It's not error type
	if err := makePointer(); err == nil {
		return
	}
}
