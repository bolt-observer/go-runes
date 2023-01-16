# go-runes

## Introduction

This is a Go implementation of [runes](https://github.com/rustyrussell/runes).
Runes are authentication codes that can be extended by third parties.
They are somewhat similar to macaroons described in ["Macaroons: Cookies with Contextual Caveats for Decentralized Authorization in the Cloud"](https://research.google/pubs/pub41892/) just simpler (but read the first link for more details). There is a [macaroon go implementation](https://github.com/go-macaroon/macaroon) already, but runes are not really well supported yet unless you want to wrap some [native C code](https://github.com/ElementsProject/lightning/tree/master/ccan/ccan/rune).

## Why would you need this?

If you plan to interact with [CoreLightning](https://github.com/ElementsProject/lightning/) runes are a way to authenticate. Then you can remote control the node via commmando plugin. More [info](https://lightning.readthedocs.io/lightning-commando.7.html).

This lib is in no way tied to lightning, you can use runes in a standalone fashion for you application too.
