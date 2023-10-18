/*
 * Copyright (c) Mike Johnson 2023.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in
 *    the documentation and/or other materials provided with the
 *    distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its
 *    contributors may be used to endorse or promote products derived
 *    from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * “AS IS” AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
	formatter "github.com/mdigger/goldmark-formatter"
	"github.com/olekukonko/tablewriter"
)

// PoliceOrganisation represents a WDTK body
type PoliceOrganisation struct {
	IsDefunct                            bool   `json:"Is_Defunct"`
	WDTKID                               string `json:"WDTK_ID"`
	WDTKOrgJSONURL                       string `json:"WDTK_Org_JSON_URL"`
	WDTKAtomFeedURL                      string `json:"WDTK_Atom_Feed_URL"`
	WDTKOrgPageURL                       string `json:"WDTK_Org_Page_URL"`
	WDTKJSONFeedURL                      string `json:"WDTK_JSON_Feed_URL"`
	DisclosureLogURL                     string `json:"Disclosure_Log_URL"`
	HomePageURL                          string `json:"Home_Page_URL"`
	Name                                 string `json:"Name"`
	PublicationSchemeURL                 string `json:"Publication_Scheme_URL"`
	DataProtectionRegistrationIdentifier string `json:"Data_Protection_Registration_Identifier"`
	WikiDataIdentifier                   string `json:"WikiData_Identifier"`
	LoCAuthorityID                       string `json:"LoC_Authority_ID"`
	FOIEmailAddress                      string `json:"FOI_Email_Address"`
}
type JSONResponse struct {
	Id                int
	UrlName           string     `json:"wdtk_id"`
	Name              string     `json:"name"`
	ShortName         string     `json:"WDTK_ID"`
	CreatedAt         string     `json:"created_At"`
	UpdatedAt         string     `json:"updated_at"`
	HomePage          string     `json:"home_page"`
	Notes             string     `json:"notes"`
	PublicationScheme string     `json:"publication_scheme"`
	DisclosureLog     string     `json:"disclosure_log"`
	Tags              [][]string `json:"tags"`
}

func main() {
	// Parse command line arguments
	reportFlag := flag.Bool(
		"report",
		false,
		"Generate a report of missing Pub Scheme and Disc. Log URLs.")

	generateFlag := flag.Bool(
		"generate",
		true,
		"Generate a dataset from the emails file, then build table.")

	refreshFlag := flag.Bool(
		"refresh",
		false,
		"Rebuilds a dataset from the emails file, then build table.")

	retainFlag := flag.Bool("retain", false, "Keep the authorities file from MySociety.")
	// wikidataFlag := flag.Bool("wikidata", false, "Get a listing of local forces from wikidata.")

	mysocietyFlag := flag.Bool(
		"mysociety",
		false,
		"Get all the mysociety data and create a table.")

	singleFlag := flag.Bool(
		"single",
		false,
		"Run a single function (for testing).")

	flag.Parse()
	if *reportFlag {
		GenerateProblemReports()
		os.Exit(0)
	}
	if *singleFlag {
		TestFunction("the_met")
		os.Exit(0)
	}

	if *mysocietyFlag {
		GetCSVDatasetFromMySociety()
		ProcessMySocietyDataset()
		os.Exit(0)
	}

	if *refreshFlag {
		RebuildDataset()
		MakeTableFromGeneratedDataset()
		Cleanup(*retainFlag)
		os.Exit(0)
	}

	if *generateFlag {
		MakeTableFromGeneratedDataset()
		Cleanup(*retainFlag)
		os.Exit(0)
	}

}

func TestFunction(wdtkID string) {

	// Attributes obtained from querying the site API:
	req, err := http.NewRequest("GET", "https://www.whatdotheyknow.com/body/the_met.json", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error fetching WDTK data:", err)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", resp.StatusCode)

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	var result JSONResponse
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		// println(err)
		return
	}
	// str, _ := json.MarshalIndent(result, "", "\t")
	// fmt.Println(string(str))

	p := new(PoliceOrganisation)

	p.DisclosureLogURL = result.DisclosureLog
	p.HomePageURL = result.HomePage
	p.Name = result.Name
	p.PublicationSchemeURL = result.PublicationScheme

	print("Force: ", p.Name)

}

func GetCSVDatasetFromMySociety() {

	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(".", "https://www.whatdotheyknow.com/body/all-authorities.csv")

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

}

func ProcessMySocietyDataset() {
	fmt.Println("Process MySociety dataset")
	csvFile, _ := os.Open("authorities.csv")
	reader := csv.NewReader(csvFile)
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		os.Exit(1)
	}

	outputRows := make([][]string, 0)
	rowHeaders := []string{"Name", "WDTK ID", "Tags"}

	for _, row := range rows {
		var thisRow []string
		tagsList := strings.Split(row[2], "|")
		tagsListFlattened := fmt.Sprint(tagsList)
		thisRow = append(thisRow, row[0], row[1], tagsListFlattened)
		outputRows = append(outputRows, thisRow)
	}

	allOrgsFile, err := os.Create("output/all-mysociety.md")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer allOrgsFile.Close()

	table := tablewriter.NewWriter(allOrgsFile)
	table.SetHeader(rowHeaders)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	table.AppendBulk(outputRows)
	table.Render()
}

func MakeTableFromGeneratedDataset() {
	println("Generating table from the generated JSON dataset...")

	os.Create("overview2.md")
	markdownOutputFile, err := os.OpenFile("overview2.md", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println("Error opening output markdown file:", err)
		os.Exit(1)
	}

	jsonInputFile, err := os.Open("data/generated-dataset.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		os.Exit(1)
	}
	defer jsonInputFile.Close()

	var dataset []map[string]interface{}
	err = json.NewDecoder(jsonInputFile).Decode(&dataset)
	if err != nil {
		fmt.Println("Error decoding JSON file:", err)
		os.Exit(1)
	}

	markdownOutputFile.WriteString(GenerateHeader())
	results := make([]string, 0)
	for _, force := range dataset {
		if force["Is_Defunct"].(bool) {
			continue
		}
		markup := fmt.Sprintf("|%v | [Website](%v)| [wdtk page](%v)| [wdtk json](%v)| [atom feed](%v)| [json feed](%v)|",
			force["Name"], force["Home_Page_URL"], force["WDTK_Org_Page_URL"],
			force["WDTK_Org_JSON_URL"], force["WDTK_Atom_Feed_URL"], force["WDTK_JSON_Feed_URL"])
		if force["Publication_Scheme_URL"] != nil {
			markup += fmt.Sprintf("| [Link](%v)|", force["Publication_Scheme_URL"])
		} else {
			markup += "| Missing|"
		}
		if force["Disclosure_Log_URL"] != nil {
			markup += fmt.Sprintf(" [Link](%v)|", force["Disclosure_Log_URL"])
		} else {
			markup += "| Missing|"
		}
		markup += fmt.Sprintf(" [Email](mailto:%v)|\n", force["FOI_Email_Address"])
		results = append(results, markup)

		// For ease of reading
		sort.Strings(results)
	}
	// Finally write all the results to the file
	for _, rowOfMarkup := range results {
		markdownOutputFile.WriteString(rowOfMarkup)
	}

	markdownOutputFile.Close()
	FormatMarkdownFile("overview2.md")
}

func Cleanup(retain bool) {

	if fileInfo, err := os.Stat("output/authorities.csv"); err == nil && fileInfo.Mode().IsRegular() && retain {
		fmt.Println("As requested, not deleting the authorities file.")
	} else if err == nil {
		println("Removing authorities.csv file.")
		os.Remove("output/authorities.csv")
	} else {
		fmt.Println("The WDTK CSV file does not exist, so could not be deleted.")
	}
}

// NewPoliceOrganisation creates a new PoliceOrganisation instance
func NewPoliceOrganisation(wdtkID string, emails map[string]string) *PoliceOrganisation {
	var policeOrg = new(PoliceOrganisation)
	println("Constructing a new PoliceOrganisation object: ", wdtkID)
	// Defaults
	policeOrg.IsDefunct = false
	// Attributes derived from the WDTK Url Name:
	policeOrg.WDTKOrgJSONURL = fmt.Sprintf("https://www.whatdotheyknow.com/body/%s.json", wdtkID)
	policeOrg.WDTKAtomFeedURL = fmt.Sprintf("https://www.whatdotheyknow.com/feed/body/%s", wdtkID)
	policeOrg.WDTKOrgPageURL = fmt.Sprintf("https://www.whatdotheyknow.com/body/%s", wdtkID)
	policeOrg.WDTKJSONFeedURL = fmt.Sprintf("https://www.whatdotheyknow.com/feed/body/%s.json", wdtkID)

	// Attributes obtained from querying the site API:
	req, err := http.NewRequest("GET", policeOrg.WDTKOrgJSONURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	resp, err := http.Get(policeOrg.WDTKOrgJSONURL)
	if err != nil {
		fmt.Println("Error fetching WDTK data:", err)
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	println(policeOrg.WDTKOrgJSONURL)
	println(string(bodyBytes))
	//
	//var wdtkData map[string]interface{}
	//err = json.NewDecoder(resp.Body).Decode(wdtkData)
	//if err != nil {
	//	fmt.Println("Error decoding WDTK data:", err)
	//	return nil
	//}
	//
	//policeOrg.DisclosureLogURL = wdtkData["disclosure_log"].(string)
	//policeOrg.HomePageURL = wdtkData["home_page"].(string)
	//policeOrg.Name = wdtkData["name"].(string)
	//policeOrg.PublicationSchemeURL = wdtkData["publication_scheme"].(string)
	//
	//// Process tags
	//for _, tag := range wdtkData["tags"].([]interface{}) {
	//	tagData := tag.([]interface{})
	//	switch tagData[0] {
	//	case "dpr":
	//		policeOrg.DataProtectionRegistrationIdentifier = tagData[1].(string)
	//	case "wikidata":
	//		policeOrg.WikiDataIdentifier = tagData[1].(string)
	//	case "lcnaf":
	//		policeOrg.LoCAuthorityID = tagData[1].(string)
	//	case "defunct":
	//		policeOrg.IsDefunct = true
	//	}
	//}

	return policeOrg
}

func RebuildDataset() {
	// Generates a JSON dataset from two pieces of information - the WDTK ID and
	// the list of FOI email addresses. The rest of the information can be either
	// derived from the WDTK ID or scraped using the WDTK ID when the object is
	// created.
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
	//
	var listOfForces []PoliceOrganisation
	//
	//// Iterate through emails
	for entry := range emails {
		println("Read ", entry)
		time.Sleep(2000)
		var force = NewPoliceOrganisation(entry, emails)
		//	println(force.PublicationSchemeURL)
		//	//println(force.Name)
		listOfForces = append(listOfForces, *force)
	}

	//// Write dataset to JSON file
	outFile, err := os.Create("data/generated-dataset2.json")
	if err != nil {
		fmt.Println("Error creating output JSON file:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	jsonData, err := json.MarshalIndent(listOfForces, "", "    ")
	if err != nil {
		fmt.Println("Error encoding JSON data:", err)
		os.Exit(1)
	}
	//
	_, err = outFile.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to output JSON file:", err)
		os.Exit(1)
	}
}

func FormatMarkdownFile(filePath string) {
	tmpfilename := "tmp-" + filePath

	err := os.Rename(filePath, tmpfilename)
	if err != nil {
		log.Fatal(err)
	}

	inFile, err := os.ReadFile(tmpfilename)
	if err != nil {
		log.Fatal(err)
	}

	outFile, _ := os.Create(filePath)
	_ = formatter.Format(inFile, outFile)
	err = outFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	os.Remove(tmpfilename)

}

func ReadCSVFile(filePath string) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	headers := records[0]
	rows := make([]map[string]string, 0)

	for _, record := range records[1:] {
		row := make(map[string]string)
		for i, value := range record {
			if i < len(headers) {
				row[headers[i]] = value
			}
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func GenerateHeader() string {
	body := "# Generated List of Police Forces (WhatDoTheyKnow)\n\n\n"
	body += "**Generated from data provided by WhatDoTheyKnow. please contact\n"
	body += "them with corrections. This table will be corrected when the "
	body += "script next runs.**\n\n"
	body += "[OPML File](police.opml)\n\n"
	body += "|Body|Website|WDTK Page|JSON|Feed: Atom|Feed: JSON|Publication Scheme|Disclosure Log|Email|\n"
	body += "|-|-|-|-|-|-|-|-|-|\n"
	return body
}

func MakeDataset() {
	fmt.Println("Function to make dataset not implemented.")
}
func GenerateReportHeader(title string) string {
	// 1: Missing Publication Scheme and Disclosure Log
	hdr := fmt.Sprintf("## %s\n|Name|Org Page|Email|\n|-|-|-|\n", title)
	return hdr
}

// GenerateProblemReports generates problem reports based on the provided dataset
func GenerateProblemReports() {
	// Opening the Dataset of Police Forces.
	datasetFile, err := os.ReadFile("data/generated-dataset.json")
	if err != nil {
		fmt.Println("Error reading FOI emails JSON:", err)
		os.Exit(1)
	}

	// and unmarshal that dataset into slice containing PoliceOrganisation objects
	var listOfForces []PoliceOrganisation
	err = json.Unmarshal(datasetFile, &listOfForces)

	// Prepare output files: Markdown page
	reportMarkdownFile, err := os.Create("missing-data.md")
	if err != nil {
		fmt.Println("Error creating report markdown file:", err)
		os.Exit(1)
	}
	defer reportMarkdownFile.Close()

	// Read emails from JSON file.
	emailsData, err := os.ReadFile("data/foi-emails.json")
	if err != nil {
		fmt.Println("Error reading FOI emails JSON:", err)
		os.Exit(1)
	}
	// and unmarshal Forces/Emails into a slice.
	var emails map[string]string
	err = json.Unmarshal(emailsData, &emails)
	if err != nil {
		fmt.Println("Error decoding FOI emails JSON:", err)
		os.Exit(1)
	}

	// Query for Police and Crime Commissioners
	reportMarkdownFile.WriteString(GenerateReportHeader("Police and Crime Commissioners"))
	for _, force := range listOfForces {
		if strings.Contains(force.Name, "Commissioner") {
			row := "| "
			row += force.Name
			row += " | "
			row += force.WDTKOrgPageURL
			row += " | "
			row += force.FOIEmailAddress
			row += " |"
			_, err = reportMarkdownFile.WriteString(row + "\n")
		}
	}
	// Query for Police and Crime Commissioners
	reportMarkdownFile.WriteString("\n")
	reportMarkdownFile.WriteString(GenerateReportHeader("Police and Crime Panels"))
	for _, force := range listOfForces {
		if strings.Contains(force.Name, "Panel") {
			row := "| "
			row += force.Name
			row += " | "
			row += force.WDTKOrgPageURL
			row += " | "
			row += force.FOIEmailAddress
			row += " |\n"
			_, err = reportMarkdownFile.WriteString(row)
		}
	}

	// Query for Missing Publication Scheme and Disclosure Logs
	reportMarkdownFile.WriteString(GenerateReportHeader("Disclosure Log and Publication Scheme Missing"))
	for _, force := range listOfForces {
		if force.IsDefunct == false && force.DisclosureLogURL == "" && force.PublicationSchemeURL == "" {
			row := "| "
			row += force.Name
			row += " | "
			row += force.WDTKOrgPageURL
			row += " | "
			row += force.FOIEmailAddress
			row += " |"
			_, err = reportMarkdownFile.WriteString(row + "\n")
		}
	}

	// Query for Missing Publication Scheme but Disclosure Log Present
	reportMarkdownFile.WriteString(GenerateReportHeader("Missing Publication Scheme but Disclosure Log Present"))
	for _, force := range listOfForces {
		if force.IsDefunct == false && strings.Contains(force.DisclosureLogURL, "http") && force.PublicationSchemeURL == "" {
			row := "| "
			row += force.Name
			row += " | "
			row += force.WDTKOrgPageURL
			row += " | "
			row += force.FOIEmailAddress
			row += " |"
			_, err = reportMarkdownFile.WriteString(row + "\n")
		}
	}
	// Query for Missing Disclosure Log but Publication Scheme Present
	reportMarkdownFile.WriteString(GenerateReportHeader("Missing Publication Scheme but Disclosure Log Present"))
	for _, force := range listOfForces {
		if force.IsDefunct == false && strings.Contains(force.PublicationSchemeURL, "http") && force.DisclosureLogURL == "" {
			row := "| "
			row += force.Name
			row += " | "
			row += force.WDTKOrgPageURL
			row += " | "
			row += force.FOIEmailAddress
			row += " |"
			_, err = reportMarkdownFile.WriteString(row + "\n")
		}
	}

	// Housekeeping
	reportMarkdownFile.Close()
	FormatMarkdownFile(reportMarkdownFile.Name())
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
