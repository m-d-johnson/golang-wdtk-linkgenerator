/*
 *
 *  * Copyright (c) Mike Johnson 2024.
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

// TODO: There's too much in this file, it needs to be carved up into smaller chunks.
//       This has started with Templates stuff but perhaps:
//       - utils & markdown-related stuff
//       - ingest
//       - report generation
//       - query
// TODO: The functions are too large and need to be carved up, particularly so it's
//

// TODO: There are far too many hardcoded file paths in here, and by now I should be doing something
//       about it.

// TODO: This needs a better name, it's too ambiguous.
// Record is an entry in the manually curated file which provides data that WDTK does not.
type Record struct {
	WDTKID           string `json:"wdtk_id"`
	TelephoneGeneral string `json:"telephone_general"`
	TelephoneFOI     string `json:"telephone_foi"`
	EmailGeneral     string `json:"email_general"`
	EmailFOI         string `json:"email_foi"`
	PostalAddress    string `json:"postal_address"`
}

// Authority is a public body on the WDTK website. Here, it contains additional fields which are not
// provided by WhatDoTheyKnow. These fields are added from a manually maintained list. This struct
// is essentially how we aggregate information from different sources and bring it together to use.
// That may by outputting it in a human (Markdown) or machine-readable (JSON) format, for example.
// Of course, we can also create these structs from the JSON we generated in the first place.
// See also the `NewAuthority` function below, which is the constructor for this type.
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
	TelephoneGeneral                     string `json:"Telephone_General"`
	TelephoneFOI                         string `json:"Telephone_FOI"`
	EmailGeneral                         string `json:"Email_General"`
	PostalAddress                        string `json:"Postal_Address"`
}

// JSONResponse is a JSON object we get back when we ask for the JSON from the authority page on the
// WDTK website.  It's the API response from WDTK.
// Example: https://www.whatdotheyknow.com/body/the_met.json
type JSONResponse struct {
	Id                int
	UrlName           string     `json:"wdtk_id"`
	Name              string     `json:"name"`
	ShortName         string     `json:"WDTK_ID"`
	CreatedAt         string     `json:"created_At"` // TODO: Inconsistent capitalisation. Fix.
	UpdatedAt         string     `json:"updated_at"`
	HomePage          string     `json:"home_page"`
	Notes             string     `json:"notes"`
	PublicationScheme string     `json:"publication_scheme"`
	DisclosureLog     string     `json:"disclosure_log"`
	Tags              [][]string `json:"tags"`
}

// Convenience functions that aid (and prettify) the console output.
// Green: Successes Yellow: Warnings. Red: Error. Magenta: Informational.
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

	// Used for the experiments I've been doing around creating/writing/manipulating SQLite DBs.
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
		"Requires the wdtk url_name. Describe an authority - returns links and metadata. Intended for use on the console.")
	// TODO: This flag needs a better name, in line with the others.
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

	// Grabs a CSV file from MySociety.
	// TODO: There are two files they provide, and this needs to be clearer as to which it gets.
	if *downloadFlag {
		GetCSVDatasetFromMySociety()
		os.Exit(0)
	}

	// Generate a report of bodies that are missing FOIA Publication Schemes and Disclosure Logs.
	if *reportFlag {
		GenerateProblemReports()
		os.Exit(0)
	}

	// Describe an organisation which the user supplies on the CLI.
	//
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

	// Regenerate the dataset and overview table, includes downloading new information from WDTK.
	if *refreshFlag {
		RebuildDataset()
		MakeTableFromGeneratedDataset()
		Cleanup(*retainFlag)
		os.Exit(0)
	}
	// Just regenerate the overview table, without downloading new information from WDTK.
	if *tableFlag {
		MakeTableFromGeneratedDataset()
		Cleanup(*retainFlag)
		os.Exit(0)
	}

	red.Println("No options provided. Exiting.")

}

// DescribeAuthority shows information on a user-specified authority and creates a simple HTML page.
// Invoke with `-describe`
func DescribeAuthority(wdtkID string) {
	// TODO: This whole function needs to be violently refactored - it does things that are also
	// done elsewhere and just repeats what we're already doing when we generate a dataset.
	// It doesn't even do it as thoroughly as it's done elsewhere and it really just needs to:
	// - Check the local JSON we generate elsewhere and read details out from that, or
	// - Just create a new Authority object and let the constructor do the work.
	//
	// With that said, this was really only added as an excuse to play around with Go Templates for
	// the first time in years. At least it did that.
	//
	// It also generates console (stdout) and HTML (to a file) in the same function so at the very
	// least it needs to split that in two and have a function per output format.
	//
	// Attributes obtained from querying the site API:
	green.Println("Showing information for ", wdtkID)

	// TODO: Feels like there's repeated code here from where we get data while in the constructor
	// for Authority (NewAuthority) so this needs its own function otherwise we're eventually
	// going to change one and not the other and have odd behaviour.
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
	// TODO: Find a prettyprinter module for this, there has to be a better way.
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
	println("FOI Email:           ", p.FOIEmailAddress)
	println("WikiData Identifier: ", p.WikiDataIdentifier)
	println("LoC Authority ID:    ", p.LoCAuthorityID)
	println("ICO Reg. Identifier: ", p.DataProtectionRegistrationIdentifier)
	println("General Email:       ", p.EmailGeneral)
	println("Telephone - General: ", p.TelephoneGeneral)
	println("Telephone - FOI:     ", p.TelephoneFOI)
	println("Postal Address:      ", p.PostalAddress)

	tmpl := template.Must(template.New("Simple HTML Overview").Parse(simpleBodyOverviewPage))

	// TODO: I hate this, do something about it and don't do it again.
	var outputFilePathBuilder strings.Builder
	outputFilePathBuilder.WriteString("output/")
	outputFilePathBuilder.WriteString("summary-")
	outputFilePathBuilder.WriteString(wdtkID)
	outputFilePathBuilder.WriteString(".html")
	var outputFilePath = outputFilePathBuilder.String()
	outputHTMLFile, _ := os.Create(outputFilePath)
	defer outputHTMLFile.Close()

	// This is the bit that actually executes the template
	err = tmpl.Execute(outputHTMLFile, p)
	if err != nil {
		red.Println("Failed to output a HTML file based on this name: ", err)
	}

}

// GetCSVDatasetFromMySociety downloads a CSV dataset of all bodies that WhatDoTheyKnow tracks.
func GetCSVDatasetFromMySociety() {
	Cleanup(false)
	// create client
	client := grab.NewClient()
	req, _ := grab.NewRequest(".", "https://www.whatdotheyknow.com/body/all-authorities.csv")

	// Start download
	green.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	green.Printf("  %v\n", resp.HTTPResponse.Status)

	// Start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	// Handle errors
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
	// TODO: This whole thing needs looking at for how strings are built.
	print(Sprintf("Query MySociety dataset for a custom tag: %s", *tag))
	var csvFile, _ = os.Open("output/all-authorities.csv")
	reader := csv.NewReader(csvFile)

	// TODO: I hate this
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
		// TODO: This is duplicated functionality and also needs to use the Authority constructor.
		tagsList := strings.Split(row[3], " ")
		if slices.Contains(tagsList, *tag) && !slices.Contains(tagsList, "defunct") {
			var name = row[0]
			var urlName = row[2]
			var weblink, _ = MakeMarkdownLink("authority_web", urlName, name)
			var bodyjsonlink, _ = MakeMarkdownLink("authority_json", urlName, name)
			// TODO: Replace with string builder
			markdownRow.WriteString("| ") // Markdown row start delimiter
			markdownRow.WriteString(weblink)
			markdownRow.WriteString(" | ") // Markdown row column separator
			markdownRow.WriteString(bodyjsonlink)
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

// MakeTableFromGeneratedDataset creates a Markdown table of UK police forces from the generated
// JSON dataset and the foi-emails.txt files. The table it creates is stored in `output/`.
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
		// TODO: This is a mess, needs to be more readable. String builder.
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
	// TODO: There needs to be more input sanitisation here.
	// TODO: When the manually curated JSON is ready, this function needs to use it as input.
	// TODO: When that's done, code needs to be deleted in here.
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
	// TODO: This needs to be refactored and deleted. There needs to just be a function that takes
	// a wtdk_id and returns the email instead of passing in the foi-emails.json file all over the
	// place.

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
	// Important to rate-limit because I want to be polite with MySociety's service (and they'll
	// ratelimit me and I won't get data).
	var qps = 4
	Println(Sprintf("Rate-limiting to %d queries/sec", qps))
	var rl = ratelimit.New(qps) // per second

	// Iterate through emails. This function only creates records for authorities listed in the
	// foi-emails.json file.
	for entry := range emails {
		_ = rl.Take()
		// TODO: The Authority constructor should be reading or looking up the emails data rather
		//       having to open and parse it to pass it in every time. There's too much inefficiency
		var force = NewAuthority(entry, emails)
		listOfForces = append(listOfForces, *force)
	}

	// Sorting because it makes diffing the dataset file much easier.
	sort.Slice(listOfForces, func(i, j int) bool { return listOfForces[i].Name < listOfForces[j].Name })

	// TODO: I really want to have JSON schema against which I can validate this.
	// TODO: Make JSON schema.
	// TODO: Use JSON schema.
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
	// TODO: This never worked right and when I look at it I feel guilty of a crime. It's hack upon
	// hack to try to get around some problem with file handles not being released.
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

// GenerateProblemReports generates Markdown reports (written to a file) based on
// `data/generated-dataset.json` - it's used for finding authorities that are missing disclosure log
// and publication scheme links.
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
