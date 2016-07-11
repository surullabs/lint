package statictest

import "testing"

func TestPackage(t *testing.T) {
	err := Package(".")
	if err != nil {
		t.Fatal(err)
	}
}
