package libpeerpipe

import (
	"testing"
)

func Test(t *testing.T) {
	lpeerpipe, err := New()
	if err != nil {
		t.Errorf(err.Error())
	}

	rpeerpipe, err := New()

	rpeerpipe.Connect(lpeerpipe.GetHash())
}
