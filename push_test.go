package cyfe

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	options := PushOptions{
		ReplaceInstead: true,
	}

	send, err := Prepare("Test", "1", "", &options)
	assert.Nil(t, err)
	fmt.Printf("\n===============\n %+v\n", send)
}
