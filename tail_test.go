package tail

import (
	"io/ioutil"
	"testing"
)

func TestTailf(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	_ = f

}
