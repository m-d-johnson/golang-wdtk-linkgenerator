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
	"strings"
)

func MakeMarkdownLink(linkType string, urlName string, label string) (output string, err int) {
	// Valid linkTypes are: authority_web, authority_json, feed_json, or feed_atom

	// If this is called with a missing label, we don't necessarily want
	// this to crash the program. It's to be expected that data may be
	// incomplete because of the nature of what this tool does.Instead,
	// we should log to console that there was an attempt to create a link
	// where there might be data missing, then we should return something
	// that's still valid markdown, so it doesn't break the rendering of
	// whatever the return value of this function is inserted into.
	if len(label) == 0 {
		magenta.Println("Tried to make a markdown link but was missing the label!")
		label = "Link label unknown"
	}
	// Calling this function with a missing urlName is different as the whole
	// point of making a link is that you have a functioning link, and the WDTK
	// URL Name should never be missing in normal operation. Thus, it indicates
	// that there is a logic error and so we log it and return an error.
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

// GenerateHeader is used by code that makes tables of authorities.
func GenerateHeader() string {
	var sb strings.Builder
	sb.WriteString("# Generated List of Police Forces (WhatDoTheyKnow)\n\n\n")
	sb.WriteString("# Generated List of Police Forces (WhatDoTheyKnow)\n\n\n")
	sb.WriteString("**Generated from data provided by WhatDoTheyKnow. please contact\n")
	sb.WriteString("them with corrections. This table will be corrected when the ")
	sb.WriteString("script next runs.**\n\n")
	sb.WriteString("[OPML File](police.opml)\n\n")
	sb.WriteString("| Body | Website | WDTK Page | JSON | Feed: Atom | Feed: JSON | Publication Scheme | Disclosure Log | Email |\n")
	sb.WriteString("|-|-|-|-|-|-|-|-|-|\n")
	return sb.String()
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
