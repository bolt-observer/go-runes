package example

import (
	"crypto/rand"
	"fmt"
	"testing"

	runes "github.com/bolt-observer/go-runes/runes"
	"github.com/stretchr/testify/assert"
)

func TestUsage(t *testing.T) {

	secret := make([]byte, 55)
	_, err := rand.Read(secret)
	if err != nil {
		assert.NoError(t, err)
	}

	master := runes.MustMakeMasterRune(secret)
	peers := master.MustGetRestrictedFromString("method^list|method^get&time>100")

	fmt.Printf("%v\n", master.IsRuneAuthorized(&peers)) // returns true
	if peers.Check(map[string]any{"method": "listpeers", "time": 1674742049}) == nil {
		fmt.Printf("Check succeeded\n")
	}

	rune := runes.MustGetFromBase64("EMXekLFLz2z-I7bEOBkfQmR5bR_V78iaf-L-LeFu8Mc9MA")
	restricted := rune.MustGetRestrictedFromString("method^list|method^get|method=summary&method/listdatastore")

	fmt.Printf("Restricted rune: %s\n", restricted.ToBase64())
}
