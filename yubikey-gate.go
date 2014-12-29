package main

import (
	"code.google.com/p/goconf/conf"
	"encoding/hex"
	"flag"
	"github.com/conformal/yubikey"
	"log"
	"os"
	"strconv"
)

func getSecretKey(key string) (*yubikey.Key, error) {
	b, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	priv := yubikey.NewKey(b)

	return &priv, nil
}

func getToken(otpString string, priv *yubikey.Key) (*yubikey.Token, error) {
	_, otp, err := yubikey.ParseOTPString(otpString)
	if err != nil {
		return nil, err
	}

	t, err := otp.Parse(*priv)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func main() {
	flag.Parse()
	args := flag.Args()

	user := os.Getenv("USER")
	pass := os.Getenv("PASS")
	cf := os.Getenv("CONTEXT")
	if cf == "" {
		cf = "gate.conf"
	}
	if user == "" || pass == "" {
		if len(args) < 2 {
			log.Fatal("need more args")
		}
		user = args[0]
		pass = args[1]
	}

	c, err := conf.ReadConfigFile(cf)
	if err != nil {
		log.Fatal(err)
	}

	p, err := c.GetString(user, "key")
	if err != nil {
		log.Fatal(err)
	}
	ctr, err := c.GetInt(user, "counter")
	if err != nil {
		log.Fatal(err)
	}

	priv, err := getSecretKey(p)
	if err != nil {
		log.Fatal(err)
	}

	token, err := getToken(pass, priv)
	if err != nil {
		log.Fatal(err)
	}

	if token.Ctr < uint16(ctr) {
		log.Fatal("counter invalid")
	}

	c.AddOption(user, "counter", strconv.Itoa(int(token.Ctr)))
	c.WriteConfigFile("gate.conf", 655, "")

	log.Printf(
		"         counter: %d (0x%04x)\n"+
			" timestamp (low): %d (0x%04x)\n"+
			"timestamp (high): %d (0x%02x)\n"+
			"     session use: %d (0x%02x)\n"+
			"          random: %d (0x%02x)\n"+
			"             crc: %d (0x%04x)\n",
		token.Ctr, token.Ctr,
		token.Tstpl, token.Tstpl,
		token.Tstph, token.Tstph,
		token.Use, token.Use,
		token.Rnd, token.Rnd,
		token.Crc, token.Crc)
}
