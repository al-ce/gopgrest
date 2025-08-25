package assert

import (
	"errors"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
)

// Try fails the test if err is not nil
func Try(t *testing.T, err error) {
	if err != nil {
		file, line := callerinfo()
		log.Printf("%s:%d: %s", file, line, err.Error())
		t.FailNow()
	}
}

// IsEq fails the test if got != exp
func IsEq(t *testing.T, got any, exp any) {
	if got != exp {
		file, line := callerinfo()
		log.Printf("%sa%d: %v (%T) != %v (%T)", file, line, got, got, exp, exp)
		t.FailNow()
	}
}

// IsNotEq fails the test if got == exp
func IsNotEq(t *testing.T, got any, exp any) {
	if got == exp {
		file, line := callerinfo()
		log.Printf("%sa%d: %v (%T) == %v (%T)", file, line, got, got, exp, exp)
		t.FailNow()
	}
}

// IsTrue fails the test if condition is not true
func IsTrue(t *testing.T, condition bool) {
	if !condition {
		file, line := callerinfo()
		log.Printf("%s:%d: %v is not true", file, line, condition)
		t.FailNow()
	}
}

// ErrorsIs fails the test if errors.Is(err) is not true
func ErrorsIs(t *testing.T, err, target error) {
	if !errors.Is(err, target) {
		file, line := callerinfo()
		log.Printf("%s:%d: %v is not %v", file, line, err, target)
		t.FailNow()
	}
}

func callerinfo() (string, int) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		panic("could not get caller info")
	}
	wd, err := os.Getwd()
	if err != nil {
		panic("could not get work dir")
	}

	file = strings.Replace(file, wd+"/", "", 1)

	return file, line
}
