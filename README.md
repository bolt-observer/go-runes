# go-runes

## Introduction

This is a Go implementation of [runes](https://github.com/rustyrussell/runes).
Runes are contextual authentication tokens like cookies with the additional feature that they can be extended by third parties. (Imagine a JWT where anybody can add a further restriction.)
They are somewhat similar to macaroons described in ["Macaroons: Cookies with Contextual Caveats for Decentralized Authorization in the Cloud"](https://research.google/pubs/pub41892/) just simpler (but read the first link for more details). There is a [macaroon go implementation](https://github.com/go-macaroon/macaroon) already, but runes are not really well supported yet unless you want to wrap some [native C code](https://github.com/ElementsProject/lightning/tree/master/ccan/ccan/rune).

Just to avoid confusion `rune` is also the name of a unicode character in Go. This has nothing to do with string manipulation except that the rune restrictions can also be represented as a string.

## Why would you need this?

If you plan to interact with Bitcoin Lightning Daemon [CoreLightning](https://github.com/ElementsProject/lightning/) runes are a way to authenticate. Then you can remote control a node via [commando plugin](https://lightning.readthedocs.io/lightning-commando.7.html).
The restriction language of CoreLightning offers a few well-known [field names](https://lightning.readthedocs.io/lightning-commando-rune.7.html) like `time`. Go-runes supports any kind of field name and restriction, but `time` by itself has no meening (and you need application logic to enforce that the field means current UNIX time).

This library is in no way tied to Lightning, you can use runes in a standalone fashion for your application too.

## Examples

Create a master rune, derive a restricted rune and check it:

```
import "github.com/bolt-observer/go-runes/runes"

secret := make([]byte, 55)
_, err := rand.Read(secret)
if err != nil {
    // handle error
}

master := runes.MustMakeMasterRune(secret)
restricted := master.MustGetRestrictedFromString("method^list|method^get&time>100")

fmt.Printf("%v\n", master.IsRuneAuthorized(&restricted)) // returns true
if peers.Check(map[string]any{"method": "listpeers", "time": 1674742049}) == nil {
    fmt.Printf("Check succeeded\n")
}
```

Obtain existing rune, restrict it and evaluate it:

```
import "github.com/bolt-observer/go-runes/runes"

rune := runes.MustGetFromBase64("EMXekLFLz2z-I7bEOBkfQmR5bR_V78iaf-L-LeFu8Mc9MA")
restricted := rune.MustGetRestrictedFromString("method^list|method^get|method=summary&method/listdatastore")
fmt.Printf("Restricted rune: %s (string representation %s)\n", restricted.ToBase64(), restricted.String())

ok, msg := restricted.Evaluate(map[string]any{"method": "listdatastore", "time": 1674742049}) // ok will be false and msg will be a verbose error
```
