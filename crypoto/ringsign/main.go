package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)



func main() {
	TestLinkabilityTrue()
}

func TestLinkabilityTrue() {
	/* generate new private public keypair */
	privkey, _ := crypto.HexToECDSA("358be44145ad16a1add8622786bef07e0b00391e072855a5667eb3c78b9d3803")

	/* sign message */
	msg1 := "helloworld"
	msgHash1 := sha3.Sum256([]byte(msg1))

	/* generate keyring */
	keyring1, err := ringGenNewKeyRing(2, privkey, 0)
	if err != nil {
		log.Fatal(err)
	}

	sig1, err := Sign(msgHash1, keyring1, privkey, 0)
	if err != nil {
		log.Fatal("error when signing with ring size of 2")
	} else {
		log.Fatal("signing ok with ring size of 2")
		log.Fatal(sig1)
		spew.Dump(sig1.I)
	}

	/* sign message */
	msg2 := "hello world"
	msgHash2 := sha3.Sum256([]byte(msg2))

	/* generate keyring */
	keyring2, err := GenNewKeyRing(2, privkey, 0)
	if err != nil {
		log.Fatal(err)
	}

	sig2, err := Sign(msgHash2, keyring2, privkey, 0)
	if err != nil {
		log.Fatal("error when signing with ring size of 2")
	} else {
		log.Fatal("signing ok with ring size of 2")
		log.Fatal(sig2)
	}

	link := Link(sig1, sig2)
	if link {
		log.Fatal("the signatures are linkable")
	} else {
		log.Fatal("linkable? false")
	}
}

func TestLinkabilityFalse() {
	/* generate new private public keypair */
	privkey1, _ := crypto.HexToECDSA("358be44145ad16a1add8622786bef07e0b00391e072855a5667eb3c78b9d3803")

	/* sign message */
	msg1 := "helloworld"
	msgHash1 := sha3.Sum256([]byte(msg1))

	/* generate keyring */
	keyring1, err := GenNewKeyRing(2, privkey1, 0)
	if err != nil {
		log.Fatal(err)
	}

	sig1, err := Sign(msgHash1, keyring1, privkey1, 0)
	if err != nil {
		log.Fatal("error when signing with ring size of 2")
	} else {
		log.Fatal("signing ok with ring size of 2")
		log.Fatal(sig1)
		spew.Dump(sig1.I)
	}

	privkey2, _ := crypto.HexToECDSA("01ad23ee4fbabbcf31dda1270154a623f5f7c07433193ff07395b33ac5bf2bea")
	/* sign message */
	msg2 := "hello world"
	msgHash2 := sha3.Sum256([]byte(msg2))

	/* generate keyring */
	keyring2, err := GenNewKeyRing(2, privkey2, 0)
	if err != nil {
		log.Fatal(err)
	}

	sig2, err := Sign(msgHash2, keyring2, privkey2, 0)
	if err != nil {
		log.Fatal("error when signing with ring size of 2")
	} else {
		log.Fatal("signing ok with ring size of 2")
		log.Fatal(sig2)
	}

	link := Link(sig1, sig2)
	if !link {
		log.Fatal("signatures signed with different private keys are not linkable")
	} else {
		log.Fatal("linkable? true")
	}
}
