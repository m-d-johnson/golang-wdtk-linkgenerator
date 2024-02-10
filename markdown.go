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
	"log"
	"os"
	"strings"

	formatter "github.com/mdigger/goldmark-formatter"
)

// MakeMarkdownLink generates a snippet of Markdown with a functioning link. It can point to various
// resources on the WDTK site (indicated with linkType) and needs only the WDTK URL name (urlName)
// and the link label text (label).
func MakeMarkdownLink(linkType string, urlName string, label string) (output string, err int) {
	// Valid linkTypes are: authority_web, authority_json, feed_json, or feed_atom

	// If this is called with a missing label, we don't necessarily want
	// this to crash the program. It's to be expected that data may be
	// incomplete because of the nature of what this tool does.Instead,
	// we should log to console that there was an attempt to create a link
	// where there might be data missing, then we should return something
	// that's still valid markdown so it doesn't break the rendering of
	// whatever the returned Markdown is inserted into.
	if len(label) == 0 {
		magenta.Println("Tried to make a markdown link but no label was provided!")
		label = "Link label unknown"
	}
	// Calling this function with a missing urlName is different, as the whole
	// point of making a link is that you have a functioning link, and the WDTK
	// URL Name should never be missing in normal operation. Thus, it indicates
	// that there is a logic error, and so we log it and return an error.
	if len(urlName) == 0 {
		red.Println("Tried to make a markdown link but was missing the link!")
		return "[]()", 1
	}

	url := ""

	switch linkType {
	case "authority_web":
		url = BuildWDTKBodyURL(urlName)
	case "authority_json":
		url = BuildWDTKJSONFeedURL(urlName)
	case "feed_atom":
		url = BuildWDTKAtomFeedURL(urlName)
	case "feed_json":
		url = BuildWDTKJSONFeedURL(urlName)
	default:
		// The default shouldn't ever happen but if it does, return 1
		return "", 1
	}

	// Now that's been set, we can begin writing the Markdown snippet:
	var sb strings.Builder
	sb.WriteString("[")
	sb.WriteString(label)
	sb.WriteString("](")
	sb.WriteString(url)
	sb.WriteString(")")

	return sb.String(), 0
}

// GenerateHeader is used by code that makes tables of authorities. It's generic enough that it's
// not just limited to police forces.
func GenerateHeader() string {
	const (
		placeholderString = `# Generated List of Police Forces (WhatDoTheyKnow)


**Generated from data provided by WhatDoTheyKnow. please contact
them with corrections. This table will be corrected when the 
script next runs.**

[OPML File](police.opml)

| Body | Website | WDTK Page | JSON | Feed: Atom | Feed: JSON | Publication Scheme | Disclosure Log | Email |
|-|-|-|-|-|-|-|-|-|
`
	)
	return placeholderString
}

// GenerateReportHeader produces a very brief header of only a name, link to the page, and email.
func GenerateReportHeader(title string) string {
	var sb strings.Builder
	sb.WriteString("## ")
	sb.WriteString(title)
	sb.WriteString("\n| Name | Org Page | Email |\n")
	sb.WriteString("|-|-|-|\n")
	return sb.String()
}

// FormatMarkdownFile runs a formatter on generated Markdown files
func FormatMarkdownFile(filePath string) {
	// TODO: This never worked right and when I look at it I feel guilty of a crime.
	// It's hack upon hack to try to get around some problem with file handles not
	// being released. What should really happen is Markdown should be stored in
	// memory while it's being generated, passed to the formatter, and then the
	// formatted Markdown written to disk. However, elsewhere in this tool, we do
	// write to disk as Markdown is being generated. This is largely because the data
	// used isn't perfect and it's useful to be able to see what was generated if a
	// run fails part- way through. This means we have to read the file contents back
	// in from disk, then re-write the formatted output back to disk. Next best
	// option is probably to write the output (while generating it) to
	// `output.md.unformatted`, then read that in here, then output it to
	// `output.md`, then do some basic sanity check before deciding to delete the
	// temp file or leave it in place.
	tmpFilePath := filePath + "-tmp"

	err := os.Rename(filePath, tmpFilePath)
	if err != nil {
		log.Fatal("Failure while renaming original file: ", err)
	}

	inFile, err := os.ReadFile(tmpFilePath)
	if err != nil {
		log.Fatal("Failure while reading in original file: ", err)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal("Failure while creating output file for reformatted Markdown: ", err)
	}

	err = formatter.Format(inFile, outFile)
	if err != nil {
		log.Fatal("The Markdown formatter failed to reformat it's input: ", err)
	}

	err = outFile.Close()
	if err != nil {
		log.Fatal("Failure while closing reformatted Markdown file: ", err)
	}

	err = os.Remove(tmpFilePath)
	if err != nil {
		yellow.Println("Failure while removing temporary file:", err)
	}
}
