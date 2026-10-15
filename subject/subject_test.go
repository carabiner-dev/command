// SPDX-FileCopyrightText: Copyright 2026 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package subject

import (
	"os"
	"path/filepath"
	"testing"

	intoto "github.com/in-toto/attestation/go/v1"
	"github.com/spf13/cobra"
)

// sha256 of "hello\n"
const helloSHA256 = "5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03"

func writeHello(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "hello.txt")
	if err := os.WriteFile(path, []byte("hello\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestAddFlags(t *testing.T) {
	t.Parallel()
	opts := &Options{}
	cmd := &cobra.Command{Use: "test"}
	opts.AddFlags(cmd)
	if err := cmd.ParseFlags([]string{"-s", "gitCommit:" + sha1Hex, "--subject", "sha256:" + sha256Hex}); err != nil {
		t.Fatal(err)
	}
	if len(opts.Subjects) != 2 {
		t.Fatalf("expected 2 subjects, got %v", opts.Subjects)
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()
	file := writeHello(t)
	for _, tt := range []struct {
		name     string
		subjects []string
		mustErr  bool
	}{
		{"empty", nil, false},
		{"digest", []string{"gitCommit:" + sha1Hex}, false},
		{"file", []string{file}, false},
		{"mixed", []string{file, "sha256:" + sha256Hex}, false},
		{"missing file", []string{filepath.Join(t.TempDir(), "nope")}, true},
		{"broken digest never falls through to files", []string{"sha256:abc"}, true},
		{"unknown algorithm", []string{"sha3:" + sha256Hex}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := (&Options{Subjects: tt.subjects}).Validate()
			if (err != nil) != tt.mustErr {
				t.Fatalf("Validate() = %v, mustErr %v", err, tt.mustErr)
			}
		})
	}
}

func TestResourceDescriptors(t *testing.T) {
	t.Parallel()
	file := writeHello(t)
	opts := &Options{
		Subjects:   []string{"GitCommit:" + sha1Hex, file},
		Algorithms: []intoto.HashAlgorithm{intoto.AlgorithmSHA256},
	}
	rds, err := opts.ResourceDescriptors()
	if err != nil {
		t.Fatal(err)
	}
	if len(rds) != 2 {
		t.Fatalf("expected 2 descriptors, got %d", len(rds))
	}
	if got := rds[0].GetDigest()["gitCommit"]; got != sha1Hex {
		t.Errorf("digest subject: got %v", rds[0].GetDigest())
	}
	if got := rds[1].GetDigest()["sha256"]; got != helloSHA256 {
		t.Errorf("file subject: got %v", rds[1].GetDigest())
	}
	if len(rds[1].GetDigest()) != 1 {
		t.Errorf("expected only the configured algorithm, got %v", rds[1].GetDigest())
	}
}

func TestResourceDescriptorsDefaultAlgorithms(t *testing.T) {
	t.Parallel()
	opts := &Options{Subjects: []string{writeHello(t)}}
	rds, err := opts.ResourceDescriptors()
	if err != nil {
		t.Fatal(err)
	}
	if got := rds[0].GetDigest()["sha256"]; got != helloSHA256 {
		t.Errorf("got %v", rds[0].GetDigest())
	}
	if _, ok := rds[0].GetDigest()["sha512"]; !ok {
		t.Errorf("expected hasher defaults to include sha512, got %v", rds[0].GetDigest())
	}
}

func TestResourceDescriptorsInvalid(t *testing.T) {
	t.Parallel()
	if _, err := (&Options{Subjects: []string{"sha256:abc"}}).ResourceDescriptors(); err == nil {
		t.Fatal("expected error for malformed digest")
	}
	if _, err := (&Options{Subjects: []string{filepath.Join(t.TempDir(), "nope")}}).ResourceDescriptors(); err == nil {
		t.Fatal("expected error for missing file")
	}
}
