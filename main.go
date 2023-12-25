/*
 *
 *  * Copyright (c) Mike Johnson 2023.
 *  *
 *  * Redistribution and use in source and binary forms, with or without
 *  * modification, are permitted provided that the following conditions
 *  * are met:
 *  *
 *  * 1. Redistributions of source code must retain the above copyright
 *  *    notice, this list of conditions and the following disclaimer.
 *  *
 *  * 2. Redistributions in binary form must reproduce the above copyright
 *  *    notice, this list of conditions and the following disclaimer in
 *  *    the documentation and/or other materials provided with the
 *  *    distribution.
 *  *
 *  * 3. Neither the name of the copyright holder nor the names of its
 *  *    contributors may be used to endorse or promote products derived
 *  *    from this software without specific prior written permission.
 *  *
 *  * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 *  * “AS IS” AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 *  * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 *  * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 *  * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 *  * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 *  * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 *  * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 *  * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 *  * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 *  * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package main

import (
	"encoding/csv"
	"encoding/json"
	. "fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/fatih/color"
	formatter "github.com/mdigger/goldmark-formatter"
	flag "github.com/spf13/pflag"
	"go.uber.org/ratelimit"
)

// Record is a an entry in the manually curated file which provides data that WDTK does not.
type Record struct {
	WDTKID           string `json:"wdtk_id"`
	TelephoneGeneral string `json:"telephone_general"`
	TelephoneFOI     string `json:"telephone_foi"`
	EmailGeneral     string `json:"email_general"`
	EmailFOI         string `json:"email_foi"`
	PostalAddress    string `json:"postal_address"`
}

// Authority is a public body on the WDTK website.
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
	TelephoneGeneral                     string `json:"TelephoneGeneral"`
	TelephoneFOI                         string `json:"TelephoneFOI"`
	EmailGeneral                         string `json:"EmailGeneral"`
	PostalAddress                        string `json:"PostalAddress"`
}

// JSONResponse is a JSON object from the authority page on the WDTK website.
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

var green = color.New(color.FgHiGreen)
var red = color.New(color.FgHiRed)
var magenta = color.New(color.FgHiMagenta)
var yellow = color.New(color.FgHiYellow)

func main() {

	// Run test function ReadCSVFileAndConvertToJson (Builds JSON dataset from downloaded CSV)
	testFlag := flag.Bool(
		"test",
		false,
		"Runs ReadCSVFileAndConvertToJson.")

	createDbFlag := flag.Bool(
		"createdb",
		false,
		"Creates a SQLite Database.")

	downloadFlag := flag.Bool(
		"download",
		false,
		"Downloads the (reduced) dataset from MySociety.")

	reportFlag := flag.Bool(
		"report",
		false,
		"Generate a report of police forces which do not list their publication scheme and disclosure log URLs.")

	tableFlag := flag.Bool(
		"table",
		false,
		"Generate a table from the existing dataset (police only).")

	refreshFlag := flag.Bool(
		"refresh",
		false,
		"Rebuilds a (police only) dataset from the emails file and API, then build table.")

	retainFlag := flag.Bool(
		"retain",
		false,
		"Keep the authorities file from MySociety.")

	describeFlag := flag.String(
		"describe",
		"",
		"Requires the wdtk url_name. Describe an authority - returns links and metadata.")

	qtag := flag.String(
		"query",
		"",
		"Tag to use to generate a user-defined report.")

	flag.Parse()

	if *testFlag {
		GetCSVDatasetFromMySociety()
		ReadCSVFileAndConvertToJson("data/whatdotheyknow_authorities_dataset.csv")
		os.Exit(0)
	}

	// Create and populate a SQLite database of the Authorities data.
	if *createDbFlag {
		CreateAndPopulateSQLiteDatabaseAll()
		os.Exit(0)
	}

	// Generate a report of bodies that are missing certain metadata.
	if *downloadFlag {
		GetCSVDatasetFromMySociety()
		os.Exit(0)
	}

	// Generate a report of bodies that are missing certain metadata.
	if *reportFlag {
		GenerateProblemReports()
		os.Exit(0)
	}

	// Describe user-supplied organisation.
	if len(*describeFlag) > 0 {
		green.Println("Authority to describe: ", *describeFlag)
		DescribeAuthority(*describeFlag)
		os.Exit(0)
	}

	// Perform query by user-supplied tag.
	if len(*qtag) > 0 {
		green.Println("Tag provided for custom query: ", *qtag)
		GetCSVDatasetFromMySociety()
		RunCustomQuery(qtag)
		os.Exit(0)
	}

	// Regenerate the dataset, includes downloading new information from WDTK.
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

	red.Println("No options provided. Exiting.")

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
func BuildWDTKBodyURL(wdtkID string) string {
	var sb strings.Builder
	sb.WriteString("https://www.whatdotheyknow.com/body/")
	sb.WriteString(wdtkID)
	return sb.String()
}
func BuildWDTKBodyJSONURL(wdtkID string) string {
	var sb strings.Builder
	sb.WriteString("https://www.whatdotheyknow.com/body/")
	sb.WriteString(wdtkID)
	sb.WriteString(".json")
	return sb.String()
}
func BuildWDTKAtomFeedURL(wdtkID string) string {
	var sb strings.Builder
	sb.WriteString("https://www.whatdotheyknow.com/feed/body/")
	sb.WriteString(wdtkID)
	return sb.String()
}
func BuildWDTKJSONFeedURL(wdtkID string) string {
	var sb strings.Builder
	sb.WriteString("https://www.whatdotheyknow.com/feed/body/")
	sb.WriteString(wdtkID)
	sb.WriteString(".json")
	return sb.String()
}

// DescribeAuthority shows information on a user-specified authority and creates a simple HTML page.
// Invoke with -describe
func DescribeAuthority(wdtkID string) {
	// Attributes obtained from querying the site API:
	green.Println("Showing information for ", wdtkID)

	req, err := http.NewRequest("GET", BuildWDTKBodyJSONURL(wdtkID), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		red.Println("Error fetching WDTK data from JSON API:", err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		red.Printf("Client: could not read response body: %s\n", err)
	}

	var result JSONResponse
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		red.Println("API JSON unmarshalling failed. Malformed?: ", err)
		return
	}

	emailsData, err := os.ReadFile("data/foi-emails.json")
	if err != nil {
		red.Println("Error reading data/foi-emails.json file:", err)
		os.Exit(1)
	}

	var emails map[string]string
	err = json.Unmarshal(emailsData, &emails)
	if err != nil {
		red.Println("Error unmarshalling data/foi-emails.json:", err)
		os.Exit(1)
	}

	var p = NewAuthority(wdtkID, emails)

	println("Force:               ", p.Name)
	println("WDTK ID:             ", p.WDTKID)
	println("Defunct:             ", p.IsDefunct)
	println("WDTK Page:           ", p.WDTKOrgPageURL)
	println("Home Page:           ", p.HomePageURL)
	println("Publication Scheme:  ", p.PublicationSchemeURL)
	println("Disclosure Log:      ", p.DisclosureLogURL)
	println("Atom Feed:           ", p.WDTKAtomFeedURL)
	println("JSON Feed:           ", p.WDTKJSONFeedURL)
	println("JSON Data:           ", p.WDTKOrgJSONURL)
	println("Email:               ", p.FOIEmailAddress)
	println("WikiData Identifier: ", p.WikiDataIdentifier)
	println("LoC Authority ID:    ", p.LoCAuthorityID)
	println("ICO Reg. Identifier: ", p.DataProtectionRegistrationIdentifier)

	tmpl := template.Must(template.New("Simple HTML Overview").Parse(simpleBodyOverviewPage))

	var outputFilePathBuilder strings.Builder
	outputFilePathBuilder.WriteString("output/")
	outputFilePathBuilder.WriteString("summary-")
	outputFilePathBuilder.WriteString(wdtkID)
	outputFilePathBuilder.WriteString(".html")
	var outputFilePath = outputFilePathBuilder.String()
	outputHTMLFile, _ := os.Create(outputFilePath)
	defer outputHTMLFile.Close()

	err = tmpl.Execute(outputHTMLFile, p)
}

// GetCSVDatasetFromMySociety downloads a CSV dataset of all bodies that WhatDoTheyKnow tracks.
// Invoke with -download
func GetCSVDatasetFromMySociety() {
	Cleanup(false)
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(".", "https://www.whatdotheyknow.com/body/all-authorities.csv")

	// start download
	green.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	green.Printf("  %v\n", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	// check for errors
	if err := resp.Err(); err != nil {
		Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	green.Printf("Download saved to ./%v \n", resp.Filename)

}

// RunCustomQuery creates a Markdown table of bodies matching a user-specified tag. It uses the
// downloaded CSV file in order to have access to all the bodies they know about (and it saves a
// load of API calls).
func RunCustomQuery(tag *string) {
	print(Sprintf("Query MySociety dataset for a custom tag: %s", *tag))
	var csvFile, _ = os.Open("output/all-authorities.csv")
	reader := csv.NewReader(csvFile)

	resultsFileNameElements := []string{"output/", "custom-query-", *tag, ".md"}
	resultsFileName := strings.Join(resultsFileNameElements, "")
	resultsTable, err := os.Create(resultsFileName)
	if err != nil {
		red.Println("Error creating output file:", err)
		os.Exit(1)
	}
	defer resultsTable.Close()

	rows, _ := reader.ReadAll()
	var title = Sprintf("# Custom Listing: Authorities with tag \"%s\"", *tag)
	_, err = resultsTable.Write([]byte(title + "\n\n"))
	if err != nil {
		red.Println("Failed to write title to results file.", err)
		return
	}

	_, err = resultsTable.Write([]byte("|Name | JSON |\n"))
	if err != nil {
		red.Println("Failed to write table header to results file.", err)
		return
	}
	_, err = resultsTable.Write([]byte("|-|-|\n"))
	if err != nil {
		red.Println("Failed to write table header separator to results file.", err)
		return
	}
	/* output/all-authorities.csv file from MySociety
	   Columns and indices in this file:
	       0:  name							    string	(Unique)
	       1:  short-name						string
	       2:  url-name						    string
	       3:  tags							    string (Pipe-delimited)
	       4:  home-page						string (URL)
	       5:  publication-scheme				string (URL)
	       6:  disclosure-log					string (URL)
	       7:  notes							string
	       8:  created-at						string (Date) e.g. 2013-09-24 12:00:27 +0100
	       9: updated-at						string (Date) e.g. 2013-09-24 12:00:27 +0100
	                                                               %Y-%m-%d %H:%M:%S %z
	                                                               RFC3339
	       10: version							int
	*/

	// Can use same StringBuilder for all rows to avoid having to re-instantiate it.
	var markdownRow strings.Builder

	for _, row := range rows {
		tagsList := strings.Split(row[3], " ")
		if slices.Contains(tagsList, *tag) && !slices.Contains(tagsList, "defunct") {
			var name = row[0]
			var urlName = row[2]
			// todo: Replace with string builder
			markdownRow.WriteString("| ") // Markdown row start delimiter
			markdownRow.WriteString(MakeMarkdownLinkToWdtkBodyPage(urlName, name))
			markdownRow.WriteString(" | ") // Markdown row column separator
			markdownRow.WriteString(MakeMarkdownLinkToWdtkBodyJson(urlName, "JSON"))
			markdownRow.WriteString(" |\n") // Markdown row end delimiter
			_, err = resultsTable.WriteString(markdownRow.String())
			if err != nil {
				red.Println("Failed to write row results file.", err)
				return
			}
			markdownRow.Reset()
		}
	}

}

// MakeTableFromGeneratedDataset creates a table of UK police forces from the generated JSON dataset
// and the foi-emails.txt files. The table it creates is rendered in Markdown and stored in output/.
func MakeTableFromGeneratedDataset() {
	magenta.Println("Generating table from the generated JSON dataset...")

	os.Create("output/overview.md")
	markdownOutputFile, err := os.OpenFile(
		"output/overview.md", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		red.Println("Error opening output markdown file:", err)
		os.Exit(1)
	}

	jsonInputFile, err := os.Open("data/generated-dataset.json")
	if err != nil {
		red.Println("Error opening JSON file:", err)
		os.Exit(1)
	}
	defer jsonInputFile.Close()

	var dataset []map[string]interface{}
	err = json.NewDecoder(jsonInputFile).Decode(&dataset)
	if err != nil {
		red.Println("Error decoding JSON file:", err)
		os.Exit(1)
	}

	_, err = markdownOutputFile.WriteString(GenerateHeader())
	if err != nil {
		red.Println("Couldn't write header to output file:", err)
		os.Exit(1)
	}

	results := make([]string, 0)
	for _, force := range dataset {
		if force["Is_Defunct"].(bool) {
			continue // Skip defunct organisations. Move to next record.
		}
		markup := Sprintf("|%v | [Website](%v)| [wdtk page](%v)| [wdtk json](%v)| [atom feed](%v)| [json feed](%v)|",
			force["Name"], force["Home_Page_URL"], force["WDTK_Org_Page_URL"],
			force["WDTK_Org_JSON_URL"], force["WDTK_Atom_Feed_URL"], force["WDTK_JSON_Feed_URL"])

		if len(force["Publication_Scheme_URL"].(string)) > 0 {
			markup += Sprintf(" [PS Link](%v)|", force["Publication_Scheme_URL"])
		} else {
			markup += " Missing |"
		}

		if len(force["Disclosure_Log_URL"].(string)) > 0 {
			markup += Sprintf(" [DL Link](%v)|", force["Disclosure_Log_URL"])
		} else {
			markup += " Missing |"
		}
		markup += Sprintf(" [Email](mailto:%v)|", force["FOI_Email_Address"])
		markup += "\n"
		results = append(results, markup)
	}

	// For ease of reading
	sort.Strings(results)

	// Finally write all the results to the file
	for _, rowOfMarkup := range results {
		markdownOutputFile.WriteString(rowOfMarkup)
	}
	green.Println("Done!")

	// Housekeeping
	err = markdownOutputFile.Close()
	if err != nil {
		red.Println("Failed to close output file: ", err)
		return
	}
}

// Cleanup deletes the existing output/all-authorities.csv file (if retain==true). If retain==false,
// then it does nothing and returns.
func Cleanup(retain bool) {
	if fileInfo, err := os.Stat("output/all-authorities.csv"); err == nil && fileInfo.Mode().IsRegular() && retain {
		yellow.Println("As requested, not deleting the output/all-authorities.csv file.")
	} else if err == nil {
		println("Removing output/all-authorities.csv file.")
		os.Remove("output/all-authorities.csv")
	} else {
		magenta.Println("output/all-authorities.csv does not exist, so could not be deleted.")
	}
}

// NewAuthority creates a new Authority instance. It uses data from the foi-emails.json file
// and some API calls.
func NewAuthority(wdtkID string, emails map[string]string) *Authority {
	var org = new(Authority)

	// Defaults
	org.IsDefunct = false

	// Attributes derived from the WDTK Url Name:
	org.WDTKID = wdtkID
	org.WDTKOrgJSONURL = BuildWDTKBodyJSONURL(wdtkID)
	org.WDTKAtomFeedURL = BuildWDTKAtomFeedURL(wdtkID)
	org.WDTKOrgPageURL = BuildWDTKBodyURL(wdtkID)
	org.WDTKJSONFeedURL = BuildWDTKJSONFeedURL(wdtkID)

	org.FOIEmailAddress = emails[wdtkID]

	// Attributes obtained from querying the site API:
	req, err := http.NewRequest("GET", org.WDTKOrgJSONURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.38 Chrome/118.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		red.Println("Error fetching WDTK data:", err)
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	println(org.WDTKOrgJSONURL)
	responsestr := string(bodyBytes)
	var wdtkData map[string]interface{}

	err = json.Unmarshal([]byte(responsestr), &wdtkData)
	if err != nil {
		red.Println("Error decoding WDTK data:", err)
		return nil
	}
	org.FOIEmailAddress = emails[wdtkID]
	org.DisclosureLogURL = wdtkData["disclosure_log"].(string)
	org.HomePageURL = wdtkData["home_page"].(string)
	org.Name = wdtkData["name"].(string)
	org.PublicationSchemeURL = wdtkData["publication_scheme"].(string)

	// Process tags
	for _, tag := range wdtkData["tags"].([]interface{}) {
		tagData := tag.([]interface{})
		switch tagData[0] {
		case "dpr":
			org.DataProtectionRegistrationIdentifier = tagData[1].(string)
		case "wikidata":
			org.WikiDataIdentifier = tagData[1].(string)
		case "lcnaf":
			org.LoCAuthorityID = tagData[1].(string)
		case "defunct":
			org.IsDefunct = true
		}
	}

	r := GetExtraDetailsFromJson()
	for _, i := range r {
		if i.WDTKID == org.WDTKID {
			if i.EmailGeneral != "" {
				org.EmailGeneral = i.EmailGeneral
			} else {
				yellow.Println("Additional data file: General Email address not present.")
			}
			if i.EmailFOI != "" {
				org.FOIEmailAddress = i.EmailFOI
			} else {
				yellow.Println("Additional data file: FOI Email address not present.")
			}
			if i.TelephoneGeneral != "" {
				org.TelephoneGeneral = i.TelephoneGeneral
			} else {
				yellow.Println("Additional data file: General Telephone number not present.")
			}
			if i.TelephoneFOI != "" {
				org.TelephoneFOI = i.TelephoneFOI
			} else {
				yellow.Println("Additional data file: FOI Telephone number not present.")
			}
			if i.PostalAddress != "" {
				org.PostalAddress = i.PostalAddress
			} else {
				yellow.Println("Additional data file: Postal Address not present.")
			}
		}
	}

	if org.IsDefunct {
		red.Println("*** This organisation is defunct ***")
	}
	return org
}

// NewAuthorityFromCSV does the same as NewAuthority, except it does it from a CSV file which is
// generated by the user. The user must download this spreadsheet and export it as CSV.
func NewAuthorityFromCSV(record []string, emails map[string]string) *Authority {
	var org = new(Authority)

	/* Converted from Excel Spreadsheet
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

	// Defaults
	if record[12] == "TRUE" {
		org.IsDefunct = true
	} else {
		org.IsDefunct = false
	}

	org.WDTKID = record[3]

	// Attributes derived from the WDTK Url Name:
	org.DisclosureLogURL = record[7]
	org.FOIEmailAddress = emails[org.WDTKID]
	org.HomePageURL = record[5]
	org.Name = record[1]
	org.PublicationSchemeURL = record[6]
	org.WDTKAtomFeedURL = BuildWDTKAtomFeedURL(org.WDTKID)
	org.WDTKJSONFeedURL = BuildWDTKJSONFeedURL(org.WDTKID)
	org.WDTKOrgJSONURL = BuildWDTKBodyJSONURL(org.WDTKID)
	org.WDTKOrgPageURL = BuildWDTKBodyURL(org.WDTKID)

	// Process tags
	tagsList := strings.Split(record[4], "|")

	for _, tag := range tagsList {
		tagData := tag

		if strings.Contains(tagData, ":") {
			tagPair := strings.Split(tagData, ":")

			switch tagPair[0] {
			case "dpr":
				org.DataProtectionRegistrationIdentifier = tagPair[1]
			case "wikidata":
				org.WikiDataIdentifier = tagPair[1]
			case "lcnaf":
				org.LoCAuthorityID = tagPair[1]
			case "defunct":
				org.IsDefunct = true
			}
		}
	}
	return org
}

// GetEmailsFromJson opens the FOI emails file and returns unmarshalled JSON mapping force names to
// their FOI email addresses.
func GetEmailsFromJson() map[string]string {

	// Read emails from JSON file
	emailsData, err := os.ReadFile("data/foi-emails.json")
	if err != nil {
		Println("Error reading FOI emails JSON:", err)
		os.Exit(1)
	}

	var emails map[string]string
	err = json.Unmarshal(emailsData, &emails)
	if err != nil {
		red.Println("Error decoding FOI emails JSON:", err)
		os.Exit(1)
	}
	return emails
}

// GetExtraDetailsFromJson opens the manually curated file and returns unmarshalled JSON,
// providing information not provided by WDTK.
func GetExtraDetailsFromJson() []Record {

	// Read emails from JSON file
	extraData, err := os.ReadFile("data/manual.json")
	if err != nil {
		Println("Error reading manual.json JSON:", err)
		os.Exit(1)
	}

	var records []Record
	err = json.Unmarshal(extraData, &records)
	if err != nil {
		red.Println("Error decoding FOI emails JSON:", err)
		os.Exit(1)
	}
	return records
}

// RebuildDataset recreates the generated-dataset.json file, which is a subset of the bodies that
// MySociety knows about -- it exists to have a source of information which includes FOI emails.
func RebuildDataset() {
	// Generates a JSON dataset from two pieces of information - the WDTK url_name and the list
	// of FOI email addresses. The rest of the information can be either be derived from the
	// WDTK ID (url_name) or scraped using the WDTK ID when the object is created. This is done
	// by the Authority object constructor. This function only creates records for authorities
	// listed in the foi-emails.json file.

	// Read emails from JSON file
	emailsData, err := os.ReadFile("data/foi-emails.json")
	if err != nil {
		Println("Error reading FOI emails JSON:", err)
		os.Exit(1)
	}

	var emails map[string]string
	err = json.Unmarshal(emailsData, &emails)
	if err != nil {
		red.Println("Error decoding FOI emails JSON:", err)
		os.Exit(1)
	}

	var listOfForces []Authority
	// Important to rate-limit for the sake of MySociety's service.
	var qps = 4
	Println(Sprintf("Rate-limiting to %d queries/sec", qps))
	var rl = ratelimit.New(qps) // per second

	// Iterate through emails. This function only creates records for authorities listed in the
	// foi-emails.json file.
	for entry := range emails {
		_ = rl.Take()
		// todo: The Authority constructors could probably better read in the emails data rather
		//       having to open and parse it to pass it in every time.
		var force = NewAuthority(entry, emails)
		listOfForces = append(listOfForces, *force)
	}

	// Sorting because it makes diffing the dataset file much easier.
	sort.Slice(listOfForces, func(i, j int) bool { return listOfForces[i].Name < listOfForces[j].Name })

	// Write dataset to JSON file
	outFile, err := os.Create("data/generated-dataset.json")
	if err != nil {
		red.Println("Error creating output JSON file:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	jsonData, err := json.MarshalIndent(listOfForces, "", "    ")
	if err != nil {
		red.Println("Error encoding JSON data:", err)
		os.Exit(1)
	}

	_, err = outFile.Write(jsonData)
	if err != nil {
		red.Println("Error writing to output JSON file:", err)
		os.Exit(1)
	} else {
		green.Println("Dataset generated and saved to", outFile.Name())
	}

}

func FormatMarkdownFile(filePath string) {

	tmpFilePath := filePath + "-tmp"

	err := os.Rename(filePath, tmpFilePath)
	if err != nil {
		log.Fatal(err)
	}

	inFile, err := os.ReadFile(tmpFilePath)
	if err != nil {
		log.Fatal(err)
	}

	outFile, _ := os.Create(filePath)
	_ = formatter.Format(inFile, outFile)
	err = outFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	os.Remove(tmpFilePath)
}

// ReadCSVFileAndConvertToJson is used to convert CSV data exported from the spreadsheet MySociety
// publishes. It provides more information than the CSV file they make available to download
// programmatically.
func ReadCSVFileAndConvertToJson(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		os.Exit(1)
	}
	var orgs []Authority
	emails := GetEmailsFromJson()

	for _, record := range records[1:] {
		org := NewAuthorityFromCSV(record, emails)
		orgs = append(orgs, *org)
	}
	// Write dataset to JSON file
	outFile, err := os.Create("data/generated-dataset-offline.json")
	if err != nil {
		red.Println("Error creating output JSON file:", err)
		os.Exit(1)
	}
	defer outFile.Close()

	jsonData, err := json.MarshalIndent(orgs, "", "    ")
	if err != nil {
		red.Println("Error encoding JSON data:", err)
		os.Exit(1)
	}

	_, err = outFile.Write(jsonData)
	if err != nil {
		red.Println("Error writing to output JSON file:", err)
		os.Exit(1)
	} else {
		green.Println("Dataset generated and saved to", outFile.Name())
	}

}

// GenerateHeader is used by the code that makes the table of police forces but could be used for
// producing more general overviews for other sets of authorities.
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

// GenerateReportHeader produces a very brief header of only a name, link to the page, and email.
func GenerateReportHeader(title string) string {
	header := Sprintf("## %s\n\n|Name|Org Page|Email|\n|-|-|-|\n", title)
	return header
}

// GenerateProblemReports generates reports based on data/generated-dataset.json - it's used for
// finding authorities that are missing disclosure log and publication scheme links.
func GenerateProblemReports() {
	// Open the dataset of police forces.
	datasetFile, err := os.ReadFile("data/generated-dataset.json")
	if err != nil {
		Println("Error reading FOI emails JSON:", err)
		os.Exit(1)
	}

	// and unmarshal that dataset into slice containing Authority objects
	var listOfForces []Authority
	err = json.Unmarshal(datasetFile, &listOfForces)

	// Prepare output files: Markdown page
	reportMarkdownFile, err := os.Create("output/missing-data.md")
	if err != nil {
		Println("Error creating report markdown file:", err)
		os.Exit(1)
	}
	defer reportMarkdownFile.Close()

	// Read emails from JSON file.
	emailsData, err := os.ReadFile("data/foi-emails.json")
	if err != nil {
		Println("Error reading FOI emails JSON:", err)
		os.Exit(1)
	}
	// and unmarshal authority/email data into a slice of maps.
	var emails map[string]string
	err = json.Unmarshal(emailsData, &emails)
	if err != nil {
		Println("Error decoding FOI emails JSON:", err)
		os.Exit(1)
	}

	// Query for Police and Crime Commissioners
	reportMarkdownFile.WriteString(GenerateReportHeader("Police and Crime Commissioners"))
	reportMarkdownFile.WriteString("")
	for _, force := range listOfForces {
		// todo: Replace with string builder
		if strings.Contains(force.Name, "Commissioner") {
			row := "| "
			row += force.Name
			row += " | [WDTK Link]("
			row += force.WDTKOrgPageURL
			row += ") | [Email]("
			row += force.FOIEmailAddress
			row += ") |"
			_, err = reportMarkdownFile.WriteString(row + "\n")
		}
	}

	// Query for Police and Crime Commissioners
	reportMarkdownFile.WriteString("\n")
	reportMarkdownFile.WriteString(GenerateReportHeader("Police and Crime Panels"))
	for _, force := range listOfForces {
		// todo: Replace with string builder
		if strings.Contains(force.Name, "Panel") {
			row := "| "
			row += force.Name
			row += " | [WDTK Link]("
			row += force.WDTKOrgPageURL
			row += ") | [Email]("
			row += force.FOIEmailAddress
			row += ") |\n"
			_, err = reportMarkdownFile.WriteString(row)
		}
	}

	// Query for Missing Publication Scheme and Disclosure Logs
	reportMarkdownFile.WriteString(GenerateReportHeader("Disclosure Log and Publication Scheme Missing"))
	for _, force := range listOfForces {
		// todo: Replace with string builder
		if force.IsDefunct == false && force.DisclosureLogURL == "" && force.PublicationSchemeURL == "" {
			row := "| "
			row += force.Name
			row += " | [WDTK Link]("
			row += force.WDTKOrgPageURL
			row += ") | [Email]("
			row += force.FOIEmailAddress
			row += ") |"
			_, err = reportMarkdownFile.WriteString(row + "\n")
		}
	}

	// Query for Missing Publication Scheme but Disclosure Log Present
	reportMarkdownFile.WriteString(GenerateReportHeader("Missing Publication Scheme but Disclosure Log Present"))
	for _, force := range listOfForces {
		// todo: Replace with string builder
		if force.IsDefunct == false && strings.Contains(force.DisclosureLogURL, "http") && force.PublicationSchemeURL == "" {
			row := "| "
			row += force.Name
			row += " | [WDTK Link]("
			row += force.WDTKOrgPageURL
			row += ") | [Email]("
			row += force.FOIEmailAddress
			row += ") |"
			_, err = reportMarkdownFile.WriteString(row + "\n")
		}
	}

	// Query for Missing Disclosure Log but Publication Scheme Present
	reportMarkdownFile.WriteString(GenerateReportHeader("Missing Publication Scheme but Disclosure Log Present"))
	for _, force := range listOfForces {
		// todo: Replace with string builder
		if force.IsDefunct == false && strings.Contains(force.PublicationSchemeURL, "http") && force.DisclosureLogURL == "" {
			row := "| "
			row += force.Name
			row += " | [WDTK Link]("
			row += force.WDTKOrgPageURL
			row += ") | [Email]("
			row += force.FOIEmailAddress
			row += ") |"
			_, err = reportMarkdownFile.WriteString(row + "\n")
		}
	}

	// Housekeeping
	reportMarkdownFile.Close()
	FormatMarkdownFile(reportMarkdownFile.Name())
	reportMarkdownFile.Close()
}
