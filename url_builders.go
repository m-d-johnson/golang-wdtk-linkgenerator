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

import "strings"

// TODO: These can be collapsed into one function like the Markdown writer

// TODO: Need to decide what these should do with empty string input (and write tests for it)

// BuildWDTKBodyURL generates a URL for an Authority's profile page on the WDTK website.
func BuildWDTKBodyURL(wdtkID string) string {
	var sb strings.Builder
	sb.WriteString("https://www.whatdotheyknow.com/body/")
	sb.WriteString(wdtkID)
	return sb.String()
}

// BuildWDTKBodyJSONURL generates a URL for an Authority's data represented as JSON on the WDTK
// website.
func BuildWDTKBodyJSONURL(wdtkID string) string {
	var sb strings.Builder
	sb.WriteString("https://www.whatdotheyknow.com/body/")
	sb.WriteString(wdtkID)
	sb.WriteString(".json")
	return sb.String()
}

// BuildWDTKAtomFeedURL generates a URL for an Atom feed of an Authority's FOI requests.
func BuildWDTKAtomFeedURL(wdtkID string) string {
	var sb strings.Builder
	sb.WriteString("https://www.whatdotheyknow.com/feed/body/")
	sb.WriteString(wdtkID)
	return sb.String()
}

// BuildWDTKJSONFeedURL generates a URL for a JSON feed of an Authority's FOI requests.
func BuildWDTKJSONFeedURL(wdtkID string) string {
	var sb strings.Builder
	sb.WriteString("https://www.whatdotheyknow.com/feed/body/")
	sb.WriteString(wdtkID)
	sb.WriteString(".json")
	return sb.String()
}
