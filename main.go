package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-tfe"
)

var client *tfe.Client

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, help)
		os.Exit(1)
	}

	var err error
	client, err = tfe.NewClient(tfe.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed initializing client: %v", err)
	}

	org := os.Args[1]
	dir := os.Args[2]

	// Check to make sure the directory exists first.
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create destination dir: %v", err)
	}

	setsResponse, err := client.PolicySets.List(context.Background(), org,
		tfe.PolicySetListOptions{})
	if err != nil {
		log.Fatalf("Failed listing policy sets: %v", err)
	}

	for _, set := range setsResponse.Items {
		// Skip empty policy sets
		if set.PolicyCount == 0 {
			continue
		}

		setDir := filepath.Join(dir, set.Name)
		log.Printf("Mirroring policy set %q into %s", set.Name, setDir)

		if err := os.Mkdir(setDir, 0755); err != nil {
			log.Fatalf("Failed creating policy set dir: %v", err)
		}

		for _, p := range set.Policies {
			policy, err := client.Policies.Read(context.Background(), p.ID)
			if err != nil {
				log.Fatalf("Failed reading policy: %v", err)
			}

			if err := downloadPolicy(setDir, policy); err != nil {
				log.Fatalf("Failed downloading policy: %v", err)
			}

			if err := writePolicyConfig(setDir, policy); err != nil {
				log.Fatalf("Failed writing policy config: %v", err)
			}
		}
	}
}

func downloadPolicy(dir string, policy *tfe.Policy) error {
	body, err := client.Policies.Download(context.Background(), policy.ID)
	if err != nil {
		return fmt.Errorf("failed downloading policy: %v", err)
	}
	body = append(body, '\n')

	path := filepath.Join(dir, fmt.Sprintf("%s.sentinel", policy.Name))
	if err := ioutil.WriteFile(path, body, 0644); err != nil {
		return fmt.Errorf("failed writing policy file: %v", err)
	}

	return nil
}

func writePolicyConfig(dir string, policy *tfe.Policy) error {
	path := filepath.Join(dir, "sentinel.hcl")
	fh, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fh.Close()

	block := fmt.Sprintf(policyTemplate, policy.Name, policy.Enforce[0].Mode)
	if _, err := fh.WriteString(block); err != nil {
		return err
	}

	return nil
}

const policyTemplate = `policy %q {
    enforcement_level = %q
}
`

const help = `Usage: ./tfc-policy-sets-migration <organization> <output-dir>

Mirrors all manually-managed policy sets in a Terraform Cloud organization to
files in a directory structure on the local disk. This directory structure may
then be checked into a version control system in order to leverage the VCS
integration with policy sets in Terraform Cloud.
`
