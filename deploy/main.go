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
	"slices"
	"sort"
	"strings"

	"github.com/fatih/color"
	formatter "github.com/mdigger/goldmark-formatter"
	"go.uber.org/ratelimit"
)

// Authority represents a WDTK body
type Authority struct {
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

func MakeMarkdownLinkToWdtkBodyPage(urlName string, label string) string {
	markupElements := []string{"[", label, "](", "https://www.whatdotheyknow.com/body/", urlName, ")"}
	markup := strings.Join(markupElements, "")
	return markup
}
func MakeMarkdownLinkToWdtkBodyJson(urlName string, label string) string {
	markupElements := []string{"[", label, "](", "https://www.whatdotheyknow.com/body/", urlName, ".json)"}
	markup := strings.Join(markupElements, "")
	return markup
}
func main() {
	// Parse command line arguments
	reportFlag := flag.Bool(
		"report",
		false,
		"Generate a report of missing Pub. Scheme and Disc. Log URLs.")

	tableFlag := flag.Bool(
		"table",
		true,
		"Generate a table from the existing dataset.")

	refreshFlag := flag.Bool(
		"refresh",
		false,
		"Rebuilds a dataset from the emails file and API, then build table.")

	retainFlag := flag.Bool("retain", false, "Keep the authorities file from MySociety.")
	// wikidataFlag := flag.Bool("wikidata", false, "Get a listing of local forces from wikidata.")

	singleFlag := flag.Bool(
		"single",
		false,
		"Run a single function (for testing).")

	// qtag := "none"
	qtag := flag.String("query", "none", "Tag to use to generate a user-defined report.")

	flag.Parse()
	fmt.Println("qtag:", qtag)
	if *reportFlag {
		GenerateProblemReports()
		os.Exit(0)
	}
	if *singleFlag {
		TestFunction("nhs")
		os.Exit(0)
	}

	// Perform query by user-supplied tag.
	if *qtag != "" {
		println("Tag provided for custom query: ", *qtag)
		GetCSVDatasetFromMySociety()
		CustomQuery(qtag)
		os.Exit(0)
	}

	if *refreshFlag {
		RebuildDataset()
		MakeTableFromGeneratedDataset()
		Cleanup(*retainFlag)
		os.Exit(0)
	}

	if *tableFlag {
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

	p := new(Authority)

	p.DisclosureLogURL = result.DisclosureLog
	p.HomePageURL = result.HomePage
	p.Name = result.Name
	p.PublicationSchemeURL = result.PublicationScheme

	print("Force: ", p.Name)

}

func GetCSVDatasetFromMySociety() {
	//
	//// create client
	//client := grab.NewClient()
	//req, _ := grab.NewRequest(".", "https://www.whatdotheyknow.com/body/all-authorities.csv")
	//
	//// start download
	//fmt.Printf("Downloading %v...\n", req.URL())
	//resp := client.Do(req)
	//fmt.Printf("  %v\n", resp.HTTPResponse.Status)

}

func CustomQuery(tag *string) {
	fmt.Sprintf("Query MySociety dataset for a custom tag: %s", *tag)
	csvFile, _ := os.Open("all-authorities.csv")
	reader := csv.NewReader(csvFile)

	resultsFileNameElements := []string{"output/", "custom-query-", *tag, ".md"}
	resultsFileName := strings.Join(resultsFileNameElements, "")
	resultsTable, err := os.Create(resultsFileName)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer resultsTable.Close()

	// // Read first line for headers.
	// rowHeaders, _ := reader.Read()
	rows, _ := reader.ReadAll()
	var title = fmt.Sprintf("# Custom Listing: Authorities with tag \"%s\"", *tag)
	resultsTable.Write([]byte(title + "\n\n"))
	resultsTable.Write([]byte("|Name | JSON |\n"))
	resultsTable.Write([]byte("|-|-|\n"))

	/* all-authorities.csv file from MySociety
	   Columns and indices in this file:
	       0:  id								string
	       1:  name							    string	(Unique)
	       2:  short-name						string
	       3:  url-name						    string
	       4:  tags							    string (Pipe-delimited)
	       5:  home-page						string (URL)
	       6:  publication-scheme				string (URL)
	       7:  disclosure-log					string (URL)
	       8:  notes							string
	       9:  created-at						string (Date) e.g. 2013-09-24 12:00:27 +0100
	       10: updated-at						string (Date) e.g. 2013-09-24 12:00:27 +0100
	                                                               %Y-%m-%d %H:%M:%S %z
	                                                               RFC3339
	       11: version							int
	       12: defunct							bool
	       13: categories						string (Pipe-delimited)
	       14: top-level-categories			    string
	       15: single-top-level-category		string
	*/

	//headers := strings.Join(rowHeaders[:], " | ")
	//fmt.Print("| " + headers + " |")
	var markdownRow strings.Builder

	for _, row := range rows {
		tagsList := strings.Split(row[3], " ")
		if slices.Contains(tagsList, *tag) && !slices.Contains(tagsList, "defunct") {
			var name = row[0]
			var urlName = row[2]

			markdownRow.WriteString("| ") // Markdown row start delimiter
			markdownRow.WriteString(MakeMarkdownLinkToWdtkBodyPage(urlName, name))
			markdownRow.WriteString(" | ") // Markdown row column separator
			markdownRow.WriteString(MakeMarkdownLinkToWdtkBodyJson(urlName, "JSON"))
			markdownRow.WriteString(" |\n") // Markdown row end delimiter
			resultsTable.WriteString(markdownRow.String())
			markdownRow.Reset()
		}
	}

	//table := tablewriter.NewWriter(resultsTable)
	//table.SetHeader(rowHeaders)
	//table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	//table.SetCenterSeparator("|")
	//table.SetColumnSeparator("|")
	//table.SetRowSeparator("-")
	//table.AppendBulk(outputRows)
	//table.Render()
}

func MakeTableFromGeneratedDataset() {
	println("Generating table from the generated JSON dataset...")

	os.Create("output/overview.md")
	markdownOutputFile, err := os.OpenFile("output/overview.md", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
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
			continue // Skip defunct organisations. Move to next record.
		}
		markup := fmt.Sprintf("|%v | [Website](%v)| [wdtk page](%v)| [wdtk json](%v)| [atom feed](%v)| [json feed](%v)|",
			force["Name"], force["Home_Page_URL"], force["WDTK_Org_Page_URL"],
			force["WDTK_Org_JSON_URL"], force["WDTK_Atom_Feed_URL"], force["WDTK_JSON_Feed_URL"])

		if strings.Contains("http", force["Publication_Scheme_URL"].(string)) {
			markup += fmt.Sprintf("| [Link](%v)|", force["Publication_Scheme_URL"])
		} else {
			markup += "| Missing |"
		}

		if strings.Contains("http", force["Disclosure_Log_URL"].(string)) {
			markup += fmt.Sprintf(" [Link](%v)|", force["Disclosure_Log_URL"])
		} else {
			markup += "| Missing |"
		}
		markup += fmt.Sprintf(" [Email](mailto:%v)|\n", force["FOI_Email_Address"])
		results = append(results, markup)
	}

	// For ease of reading
	sort.Strings(results)
	// Finally write all the results to the file
	for _, rowOfMarkup := range results {
		markdownOutputFile.WriteString(rowOfMarkup)
	}

	markdownOutputFile.Close()
	// FormatMarkdownFile("output/overview.md")
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

// NewAuthority creates a new Authority instance
func NewAuthority(wdtkID string, emails map[string]string) *Authority {
	var policeOrg = new(Authority)
	green := color.New(color.FgHiGreen)
	red := color.New(color.FgHiRed)
	magenta := color.New(color.FgHiMagenta)
	green.Println("\n\n~~ Constructing a new Authority object: ", wdtkID, " ~~")

	// Defaults
	policeOrg.IsDefunct = false
	// Attributes derived from the WDTK Url Name:
	policeOrg.WDTKID = wdtkID
	policeOrg.WDTKOrgJSONURL = fmt.Sprintf("https://www.whatdotheyknow.com/body/%s.json", wdtkID)
	policeOrg.WDTKAtomFeedURL = fmt.Sprintf("https://www.whatdotheyknow.com/feed/body/%s", wdtkID)
	policeOrg.WDTKOrgPageURL = fmt.Sprintf("https://www.whatdotheyknow.com/body/%s", wdtkID)
	policeOrg.WDTKJSONFeedURL = fmt.Sprintf("https://www.whatdotheyknow.com/feed/body/%s.json", wdtkID)

	policeOrg.FOIEmailAddress = emails[wdtkID]
	// Attributes obtained from querying the site API:
	req, err := http.NewRequest("GET", policeOrg.WDTKOrgJSONURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.38 Chrome/118.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error fetching WDTK data:", err)
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	println(policeOrg.WDTKOrgJSONURL)
	responsestr := string(bodyBytes)
	var wdtkData map[string]interface{}

	err = json.Unmarshal([]byte(responsestr), &wdtkData)
	if err != nil {
		fmt.Println("Error decoding WDTK data:", err)
		return nil
	}
	policeOrg.FOIEmailAddress = emails[wdtkID]
	print("FOI Email Address:       ")
	magenta.Println(policeOrg.FOIEmailAddress)
	policeOrg.DisclosureLogURL = wdtkData["disclosure_log"].(string)
	println("Disc. Log URL:          ", policeOrg.DisclosureLogURL)
	policeOrg.HomePageURL = wdtkData["home_page"].(string)
	println("Home Page:              ", policeOrg.HomePageURL)
	policeOrg.Name = wdtkData["name"].(string)
	println("Name:                   ", policeOrg.Name)
	policeOrg.PublicationSchemeURL = wdtkData["publication_scheme"].(string)
	println("Publication Scheme URL: ", policeOrg.PublicationSchemeURL)
	//
	// Process tags
	for _, tag := range wdtkData["tags"].([]interface{}) {
		tagData := tag.([]interface{})
		switch tagData[0] {
		case "dpr":
			policeOrg.DataProtectionRegistrationIdentifier = tagData[1].(string)
			println("Data Prot. Identifier:  ", policeOrg.DataProtectionRegistrationIdentifier)
			println("ICO Registration URL:   ", "https://ico.org.uk/ESDWebPages/Entry/"+policeOrg.DataProtectionRegistrationIdentifier)
		case "wikidata":
			policeOrg.WikiDataIdentifier = tagData[1].(string)
			println("WikiData Identifier:    ", policeOrg.WikiDataIdentifier)
			println("WikiData URL:           ", "https://www.wikidata.org/wiki/"+policeOrg.WikiDataIdentifier)
		case "lcnaf":
			policeOrg.LoCAuthorityID = tagData[1].(string)
			println("LoC Authority ID:       ", policeOrg.LoCAuthorityID)
			println("LoC Authority URL:      ", "https://lccn.loc.gov/"+policeOrg.LoCAuthorityID)
		case "defunct":
			policeOrg.IsDefunct = true
		}
	}
	if policeOrg.IsDefunct {
		red.Println("*** This organisation is defunct ***")
	}
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

	var listOfForces []Authority
	// Important to rate-limit for the sake of MySociety's service.
	rl := ratelimit.New(5) // per second
	// Iterate through emails
	for entry := range emails {
		_ = rl.Take()
		var force = NewAuthority(entry, emails)
		listOfForces = append(listOfForces, *force)
	}

	// Write dataset to JSON file
	outFile, err := os.Create("data/generated-dataset.json")
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

	// and unmarshal that dataset into slice containing Authority objects
	var listOfForces []Authority
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
}
