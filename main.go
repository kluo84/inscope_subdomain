package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func readDomainsFromFile(filename string) map[string]bool {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	domains := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domains[scanner.Text()] = true
	}
	
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return domains
}

func writeDomainsToFile(domains map[string]bool, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for domain := range domains {
		fmt.Fprintln(writer, domain)
	}
	writer.Flush()
}

func main() {
	var allDomainsFile, outOfScopeFile string
	flag.StringVar(&allDomainsFile, "all", "", "File containing all subdomains")
	flag.StringVar(&outOfScopeFile, "oos", "", "File containing out of scope subdomains")
	help := flag.Bool("help", false, "Display help information")

	flag.Parse()

	if *help || allDomainsFile == "" || outOfScopeFile == "" {
		fmt.Println("This program reads all subdomains and out of scope subdomains from files, compares them and writes the in-scope subdomains to a file.")
		fmt.Println("Usage: main -all [all_domains_file] -oos [out_of_scope_file]")
		return
	}

	allSubdomains := readDomainsFromFile(allDomainsFile)
	outOfScopeSubdomains := readDomainsFromFile(outOfScopeFile)

	inScopeSubdomains := make(map[string]bool)
	for subdomain := range allSubdomains {
		if _, ok := outOfScopeSubdomains[subdomain]; !ok {
			inScopeSubdomains[subdomain] = true
		}
	}

	writeDomainsToFile(inScopeSubdomains, "in_scope_subdomains.txt")

	fmt.Println("In scope subdomains are written to 'in_scope_subdomains.txt'")
}

