package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// GenerateProblemReports generates problem reports based on the provided dataset
func GenerateProblemReports() {
	datasetFile, err := os.ReadFile("data/generated-dataset.json")
	if err != nil {
		fmt.Println("Error reading FOI emails JSON:", err)
		os.Exit(1)
	}

	var listOfForces []map[string]interface{} // policeOrg := new(PoliceOrganisation)
	err = json.Unmarshal(datasetFile, &listOfForces)

	reportOutputFile, err := os.Create("Zreport-missing_data.txt")
	if err != nil {
		fmt.Println("Error creating report output file:", err)
		os.Exit(1)
	}
	defer reportOutputFile.Close()

	reportMarkdownFile, err := os.Create("Zmissing-data.md")
	if err != nil {
		fmt.Println("Error creating report markdown file:", err)
		os.Exit(1)
	}
	defer reportMarkdownFile.Close()

	confirm := "no" // Set to "yes" to confirm email sending
	if confirm != "yes" {
		fmt.Println("Not sending emails")
	} else {
		fmt.Println("Will send emails, Ctrl+C in the next 10s to cancel")
		time.Sleep(10 * time.Second)
	}

	// 1: Missing Publication Scheme and Disclosure Log
	_, err = reportMarkdownFile.WriteString("## Missing Pub. Scheme and Disclosure Log\n\n")
	if err != nil {
		return
	}
	_, err = reportMarkdownFile.WriteString("|Name|Org Page|Email|\n")
	if err != nil {
		return
	}
	_, err = reportMarkdownFile.WriteString("|-|-|-|\n")
	if err != nil {
		return
	}
	println("Invoked RebuildDataset")
	// Read emails from JSON file
	emailsData, err := os.ReadFile("data/foi-emails.json")
	if err != nil {
		fmt.Println("Error reading FOI emails JSON:", err)
		os.Exit(1)
	}
	println(string(emailsData))
	var emails map[string]string
	err = json.Unmarshal(emailsData, &emails)
	if err != nil {
		fmt.Println("Error decoding FOI emails JSON:", err)
		os.Exit(1)
	}

	// forceIndex := 1
	for _, force := range listOfForces {
		println(force["Name"])
		// f := NewPoliceOrganisation(x, emails)
		//println(f.Name)
		//println(f.WDTKOrgJSONURL)
	}
	reportOutputFile.Close()
	reportMarkdownFile.Close()
	//	if !force.IsDefunct && len(force.PublicationSchemeURL) < 10 && len(force.DisclosureLogURL) < 10 {
	//		addressee := force.FOIEmailAddress
	//		messageBody := fmt.Sprintf(
	//			"Subject: Publication Scheme and Disclosure Log URLs for %s\n\n"+
	//				"Dear Data Protection Officer,\n\n"+
	//				"I have noticed that on the WhatDoTheyKnow website %s "+
	//				"that your organisation is missing a link to your FOI "+
	//				"Publication Scheme and Disclosure Log. While I'm "+
	//				"aware that some information is technically "+
	//				"exempt from FOI Disclosure on the grounds that "+
	//				"it is available on the internet, I wonder if "+
	//				"you'd be kind enough to send me the links to "+
	//				"both, please?\n\nWith thanks,\nMike Johnson\nmdj@mikejohnson.xyz\n",
	//			force.Name, force.WDTKOrgPageURL)
	//
	//		findingTxt := fmt.Sprintf("1: %s is missing Pub Scheme URLs %s\n", force.Name, force.WDTKOrgPageURL)
	//		findingMd := fmt.Sprintf("|%s|%s|%s|\n", force.Name, force.WDTKOrgPageURL, force.FOIEmailAddress)
	//
	//		fmt.Print(findingTxt)
	//		reportOutputFile.WriteString(findingTxt)
	//		reportMarkdownFile.WriteString(findingMd)
	//
	//		if confirm == "yes" {
	//			fmt.Printf("Sending email to %s\n", addressee)
	//			// sendEmail(addressee, messageBody)  // Uncomment when implementing email sending
	//		} else {
	//			fmt.Println("Email sending not confirmed. Written to log anyway.")
	//		}
	//		forceIndex++
	//	}
	//}
	//reportMarkdownFile.WriteString("\n\n")
	//
	//// 2: Missing Disclosure Log URL but Publication Scheme present
	//reportMarkdownFile.WriteString("## Missing Disclosure Log only\n\n")
	//reportMarkdownFile.WriteString("|Name|Org Page|Email|\n")
	//reportMarkdownFile.WriteString("|-|-|-|\n")
	//forceIndex = 1
	//for _, force := range dataset {
	//	if !force.IsDefunct && len(force.PublicationSchemeURL) > 10 && len(force.DisclosureLogURL) < 10 {
	//		addressee := force.FOIEmailAddress
	//		messageBody := fmt.Sprintf(
	//			"Subject: Disclosure Log URLs for %s\n\n"+
	//				"Dear Data Protection Officer,\n\n"+
	//				"I have noticed that on the WhatDoTheyKnow website "+
	//				"%s that your organisation is missing a link to your "+
	//				"FOI Disclosure Log. While I'm aware that some "+
	//				"information is technically exempt from FOI "+
	//				"Disclosure on the grounds that it is available on "+
	//				"the internet, I wonder if you'd be kind enough to "+
	//				"send me the link, please?\n\nWith thanks,\n"+
	//				"Mike Johnson\nmdj@mikejohnson.xyz\n",
	//			force.Name, force.WDTKOrgPageURL)
	//
	//		findingTxt := fmt.Sprintf("2: %s is missing Disclosure Log %s\n", force.Name, force.WDTKOrgPageURL)
	//		findingMd := fmt.Sprintf("|%s|%s|%s|\n", force.Name, force.WDTKOrgPageURL, force.FOIEmailAddress)
	//
	//		fmt.Print(findingTxt)
	//		reportOutputFile.WriteString(findingTxt)
	//		reportMarkdownFile.WriteString(findingMd)
	//
	//		if confirm == "yes" {
	//			fmt.Printf("Sending email to %s\n", addressee)
	//			// sendEmail(addressee, messageBody)  // Uncomment when implementing email sending
	//		} else {
	//			fmt.Println("Email sending not confirmed. Written to log anyway.")
	//		}
	//		forceIndex++
	//	}
	//}
	//reportMarkdownFile.WriteString("\n\n")
	//
	//// 3: Missing Pubscheme URL but Disclosure Log present
	//reportMarkdownFile.WriteString("## Missing Pubscheme fields\n\n")
	//reportMarkdownFile.WriteString("|Name|Org Page|Email|\n")
	//reportMarkdownFile.WriteString("|-|-|-|\n")
	//forceIndex = 1
	//for _, force := range dataset {
	//	if !force.IsDefunct && len(force.PublicationSchemeURL) < 10 && len(force.DisclosureLogURL) > 10 {
	//		addressee := force.FOIEmailAddress
	//		messageBody := fmt.Sprintf(
	//			"Subject: Publication URL for %s\n\n"+
	//				"Dear Data Protection Officer,\n"+
	//				"I have noticed that on the WhatDoTheyKnow "+
	//				"website %s that your organisation is missing "+
	//				"a link to your FOI Publication Scheme. While "+
	//				"I'm conscious that some information is technically "+
	//				"exempt from FOI Disclosure on the grounds that "+
	//				"it is available on the internet, I wonder if "+
	//				"you'd be kind enough to send me the link, please?\n\n"+
	//				"With thanks,\nMike Johnson\nmdj@mikejohnson.xyz\n",
	//			force.Name, force.WDTKOrgPageURL)
	//
	//		findingTxt := fmt.Sprintf("3: %s is missing Pubscheme but has Disclosure Log %s\n", force.Name, force.WDTKOrgPageURL)
	//		findingMd := fmt.Sprintf("|%s|%s|%s|\n", force.Name, force.WDTKOrgPageURL, force.FOIEmailAddress)
	//
	//		fmt.Print(findingTxt)
	//	}
	//}
}
