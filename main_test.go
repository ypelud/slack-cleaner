package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}

func TestCtv(t *testing.T) {
	ctv("20180101")
}

func TestUsage(t *testing.T) {
	usage()
}
