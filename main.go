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
	"flag"
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
	"go.uber.org/ratelimit"
)

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

func main() {

	// Parse command line arguments
	downloadFlag := flag.Bool(
		"download",
		false,
		"Downloads the dataset from MySociety.")

	reportFlag := flag.Bool(
		"report",
		false,
		"Generate a report of missing Pub. Scheme and Disc. Log URLs.")

	tableFlag := flag.Bool(
		"table",
		false,
		"Generate a table from the existing dataset.")

	refreshFlag := flag.Bool(
		"refresh",
		false,
		"Rebuilds a dataset from the emails file and API, then build table.")

	retainFlag := flag.Bool(
		"retain",
		false,
		"Keep the authorities file from MySociety.")

	describeFlag := flag.String(
		"describe",
		"",
		"Describe an authority.")

	qtag := flag.String(
		"query",
		"",
		"Tag to use to generate a user-defined report.")

	flag.Parse()
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
		red.Println("Error fetching WDTK data:", err)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		red.Printf("client: could not read response body: %s\n", err)
	}

	var result JSONResponse
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		red.Println("API JSON unmarshalling failed. Malformed?: ", err)
		return
	}

	emailsData, err := os.ReadFile("data/foi-emails.json")
	if err != nil {
		red.Println("Error reading foi-emails.json file:", err)
		os.Exit(1)
	}

	var emails map[string]string
	err = json.Unmarshal(emailsData, &emails)
	if err != nil {
		red.Println("Error unmarshalling foi-emails.json:", err)
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
	println("WikiData Identifier :", p.WikiDataIdentifier)
	println("LoC Authority ID    :", p.LoCAuthorityID)
	println("ICO Registration Identifier: ", p.DataProtectionRegistrationIdentifier)

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

// RunCustomQuery creates a markdown table of bodies matching a user-specified tag. It uses the
// downloaded CSV file in order to have access to all the bodies they know about (and it saves a
// load of API calls.
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
		red.Println("Failed to write title to results file.", err)
		return
	}
	_, err = resultsTable.Write([]byte("|-|-|\n"))
	if err != nil {
		red.Println("Failed to write table header to results file.", err)
		return
	}
	/* output/all-authorities.csv file from MySociety
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

	// Can use same StringBuilder for all rows to avoid having to re-instantiate it.
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
// and the foi-emails.txt files. The table it creates is rendered in markdown and stored in output/.
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

	markdownOutputFile.WriteString(GenerateHeader())
	results := make([]string, 0)
	for _, force := range dataset {
		if force["Is_Defunct"].(bool) {
			continue // Skip defunct organisations. Move to next record.
		}
		markup := Sprintf("|%v | [Website](%v)| [wdtk page](%v)| [wdtk json](%v)| [atom feed](%v)| [json feed](%v)|",
			force["Name"], force["Home_Page_URL"], force["WDTK_Org_Page_URL"],
			force["WDTK_Org_JSON_URL"], force["WDTK_Atom_Feed_URL"], force["WDTK_JSON_Feed_URL"])

		if strings.Contains("http", force["Publication_Scheme_URL"].(string)) {
			markup += Sprintf("| [Link](%v)|", force["Publication_Scheme_URL"])
		} else {
			markup += "| Missing |"
		}

		if strings.Contains("http", force["Disclosure_Log_URL"].(string)) {
			markup += Sprintf(" [Link](%v)|", force["Disclosure_Log_URL"])
		} else {
			markup += "| Missing |"
		}
		markup += Sprintf(" [Email](mailto:%v)|\n", force["FOI_Email_Address"])
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
		green.Println("As requested, not deleting the authorities file.")
	} else if err == nil {
		println("Removing authorities.csv file.")
		os.Remove("output/all-authorities.csv")
	} else {
		magenta.Println("The WDTK CSV file does not exist, so could not be deleted.")
	}
}

// NewAuthority creates a new Authority instance. It uses data from the
// foi-emails.json file and some API calls.
func NewAuthority(wdtkID string, emails map[string]string) *Authority {
	var policeOrg = new(Authority)

	// Defaults
	policeOrg.IsDefunct = false

	// Attributes derived from the WDTK Url Name:
	policeOrg.WDTKID = wdtkID
	policeOrg.WDTKOrgJSONURL = BuildWDTKBodyJSONURL(wdtkID)
	policeOrg.WDTKAtomFeedURL = BuildWDTKAtomFeedURL(wdtkID)
	policeOrg.WDTKOrgPageURL = BuildWDTKBodyURL(wdtkID)
	policeOrg.WDTKJSONFeedURL = BuildWDTKJSONFeedURL(wdtkID)

	policeOrg.FOIEmailAddress = emails[wdtkID]

	// Attributes obtained from querying the site API:
	req, err := http.NewRequest("GET", policeOrg.WDTKOrgJSONURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.38 Chrome/118.0.0.0 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		red.Println("Error fetching WDTK data:", err)
		return nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	println(policeOrg.WDTKOrgJSONURL)
	responsestr := string(bodyBytes)
	var wdtkData map[string]interface{}

	err = json.Unmarshal([]byte(responsestr), &wdtkData)
	if err != nil {
		red.Println("Error decoding WDTK data:", err)
		return nil
	}
	policeOrg.FOIEmailAddress = emails[wdtkID]
	policeOrg.DisclosureLogURL = wdtkData["disclosure_log"].(string)
	policeOrg.HomePageURL = wdtkData["home_page"].(string)
	policeOrg.Name = wdtkData["name"].(string)
	policeOrg.PublicationSchemeURL = wdtkData["publication_scheme"].(string)

	// Process tags
	for _, tag := range wdtkData["tags"].([]interface{}) {
		tagData := tag.([]interface{})
		switch tagData[0] {
		case "dpr":
			policeOrg.DataProtectionRegistrationIdentifier = tagData[1].(string)
		case "wikidata":
			policeOrg.WikiDataIdentifier = tagData[1].(string)
		case "lcnaf":
			policeOrg.LoCAuthorityID = tagData[1].(string)
		case "defunct":
			policeOrg.IsDefunct = true
		}
	}

	if policeOrg.IsDefunct {
		red.Println("*** This organisation is defunct ***")
	}
	return policeOrg
}

// RebuildDataset recreates the generated-dataset.json file, which is a subset of the bodies that
// MySociety knows about -- it exists to have a source of information which includes FOI emails.
func RebuildDataset() {
	// Generates a JSON dataset from two pieces of information - the WDTK ID and
	// the list of FOI email addresses. The rest of the information can be either
	// derived from the WDTK ID or scraped using the WDTK ID when the object is
	// created.

	// Read emails from JSON file
	// emailsData, err := ReadCSVFileAndGetRows("data/foi-emails.json")
	emailsData, err := os.ReadFile("data/foi-emails.json")
	if err != nil {
		Println("Error reading FOI emails JSON:", err)
		os.Exit(1)
	}
	println(string(emailsData))

	var emails map[string]string
	err = json.Unmarshal(emailsData, &emails)
	if err != nil {
		red.Println("Error decoding FOI emails JSON:", err)
		os.Exit(1)
	}

	var listOfForces []Authority
	// Important to rate-limit for the sake of MySociety's service.
	qps := 5
	Println(Sprintf("Rate-limiting to %d queries/sec", qps))
	rl := ratelimit.New(qps) // per second
	// Iterate through emails
	for entry := range emails {
		_ = rl.Take()
		var force = NewAuthority(entry, emails)
		listOfForces = append(listOfForces, *force)
	}

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
	}
}

func FormatMarkdownFile(fp string) {

	tmpfp := fp + "-tmp"

	err := os.Rename(fp, tmpfp)
	if err != nil {
		log.Fatal(err)
	}

	inFile, err := os.ReadFile(tmpfp)
	if err != nil {
		log.Fatal(err)
	}

	outFile, _ := os.Create(fp)
	_ = formatter.Format(inFile, outFile)
	err = outFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	os.Remove(tmpfp)
}

func ReadCSVFileAndGetRows(filePath string) ([]map[string]string, error) {
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
	header := Sprintf("## %s\n\n|Name|Org Page|Email|\n|-|-|-|\n", title)
	return header
}

// GenerateProblemReports generates problem reports based on the provided dataset
func GenerateProblemReports() {
	// Opening the Dataset of Police Forces.
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
	// and unmarshal Forces/Emails into a slice.
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
