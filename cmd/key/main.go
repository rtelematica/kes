// Copyright 2019 - MinIO, Inc. All rights reserved.
// Use of this source code is governed by the AGPL
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

const usage = `usage: %s <command>

    server               Start a key server.

    create               Create a new master key at a key server.
    delete               Delete a master key from a key server.

    gen                  Generate a new data key from a master key.
    dec                  Decrypt a encrypted data key using a master key.

    identify             Compute the identity of a TLS certificate.
    identity             Assign policies to identities.
    policy               Manage the key server policies.

  -h, --help             Show this list of command line optios.
`

func main() {
	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cli.Usage = func() {
		fmt.Fprintf(cli.Output(), usage, cli.Name())
	}
	cli.Parse(os.Args[1:])

	args := cli.Args()
	if len(args) < 1 {
		cli.Usage()
		os.Exit(2)
	}

	switch args[0] {
	case "server":
		server(args)
	case "create":
		createKey(args)
	case "delete":
		deleteKey(args)
	case "gen":
		generateKey(args)
	case "dec":
		decryptKey(args)
	case "identify":
		identify(args)
	case "identity":
		identity(args)
	case "policy":
		policy(args)
	default:
		cli.Usage()
		os.Exit(2)
	}
}

func failf(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
	os.Exit(1)
}

func parseCommandFlags(f *flag.FlagSet, args []string) []string {
	var parsedArgs []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			f.Parse([]string{arg})
		} else {
			parsedArgs = append(parsedArgs, arg)
		}
	}
	return parsedArgs
}

func serverAddr() string {
	if addr, ok := os.LookupEnv("KEY_SERVER"); ok {
		return addr
	}
	return "https://127.0.0.1:7373"
}

func loadClientCertificates() []tls.Certificate {
	certPath := os.Getenv("KEY_CLIENT_TLS_CERT_FILE")
	keyPath := os.Getenv("KEY_CLIENT_TLS_KEY_FILE")
	if certPath != "" || keyPath != "" {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			failf(os.Stderr, "Cannot load TLS key or cert for client auth: %s", err.Error())
		}
		return []tls.Certificate{cert}
	}
	return nil
}

func isTerm(f *os.File) bool { return terminal.IsTerminal(int(f.Fd())) }
