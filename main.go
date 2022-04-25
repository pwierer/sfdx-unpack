package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sfdxunpack/salesforce"
	"strconv"
	"strings"
)

func main() {
	var myDomainURL string
	var sessionID string

	flag.StringVar(&myDomainURL, "domain", "", "My Domain URL")
	flag.StringVar(&sessionID, "sid", "", "Session ID")
	flag.Parse()

	if myDomainURL == "" {
		myDomainURL = stringInput("My Domain URL")
	}
	if sessionID == "" {
		sessionID = stringInput("Session ID")
	}

	sf := salesforce.New(myDomainURL, sessionID)
	selectedPackage := selectUnlockedPackage(sf)
	removePackageMetadata(sf, selectedPackage)
}

// selectUnlockedPackage retrieves all installed unlocked packages and asks user to select one
func selectUnlockedPackage(sf salesforce.Salesforce) (selectedPackage salesforce.InstalledSubscriberPackage) {
	packages, err := sf.GetUnlockedPackages()
	if err != nil {
		fmt.Fprint(os.Stdout, "Could not retrieve unlocked packages")
		os.Exit(1)
	}
	if len(packages) == 0 {
		fmt.Println("No unlocked packages found")
		os.Exit(0)
	}
	fmt.Println("\nUnlocked packages: ")
	for i, pkg := range packages {
		fmt.Printf("(%d) %s\n", i+1, pkg.SubscriberPackage.Name)
	}

	packageNumber := numberInput("Select a package number", 1, len(packages))
	for i, pkg := range packages {
		if packageNumber == i+1 {
			selectedPackage = pkg
			break
		}
	}
	return selectedPackage
}

// removePackageMetadata removes all packaged metadata components of the provided package
func removePackageMetadata(sf salesforce.Salesforce, selectedPackage salesforce.InstalledSubscriberPackage) {
	confirmToken := stringInput("Confirmation Token")
	componentIds, err := sf.GetMetadataComponents(selectedPackage.SubscriberPackageId)
	if err != nil {
		fmt.Fprint(os.Stdout, "Could not retrieve package components")
		os.Exit(1)
	}
	fmt.Printf("\n# Package components: %d\n", len(componentIds))
	for i, id := range componentIds[:2] {
		err := sf.RemovePackageMetadata(selectedPackage.Id, id, confirmToken)
		if err != nil {
			fmt.Fprint(os.Stdout, "Could not remove package component")
			os.Exit(1)
		}
		fmt.Printf("%d remaining\n", len(componentIds)-i-1)
	}
}

// stringInput reads user input as string
func stringInput(label string) string {
	var input string
	for input == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s: ", label)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
	}
	return input
}

// numberInput reads user input as int
func numberInput(label string, min, max int) int {
	var num int
	for num < min || num > max {
		s := stringInput(label)
		n, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		num = n
	}
	return num
}
