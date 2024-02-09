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
	"testing"
)

// This test calls BuildWDTKBodyURL with valid input and expects a well-formed result.
func TestBuildWDTKBodyURL_ValidParameters(t *testing.T) {
	want := "https://www.whatdotheyknow.com/body/valid_wdtk_id"
	// Call function under test
	result := BuildWDTKBodyURL("valid_wdtk_id")
	if result != want {
		t.Errorf("got %s but expected %s", result, want)
	}
}

// This test calls BuildWDTKBodyJSONURL with valid input and expects a well-formed result.
func TestBuildWDTKBodyJSONURL_ValidParameters(t *testing.T) {
	want := "https://www.whatdotheyknow.com/body/valid_wdtk_id.json"
	// Call function under test
	result := BuildWDTKBodyJSONURL("valid_wdtk_id")
	if result != want {
		t.Errorf("got %s but expected %s", result, want)
	}
}

// This test calls BuildWDTKAtomFeedURL with valid input and expects a well-formed result.
func TestBuildWDTKBodyAtomFeedURL_ValidParameters(t *testing.T) {
	want := "https://www.whatdotheyknow.com/feed/body/valid_wdtk_id"
	// Call function under test
	result := BuildWDTKAtomFeedURL("valid_wdtk_id")
	if result != want {
		t.Errorf("got %s but expected %s", result, want)
	}
}

// This test calls BuildWDTKJSONFeedURL with valid input and expects a well-formed result.
func TestBuildWDTKBodyJSONFeedURL_ValidParameters(t *testing.T) {
	want := "https://www.whatdotheyknow.com/feed/body/valid_wdtk_id.json"
	// Call function under test
	result := BuildWDTKJSONFeedURL("valid_wdtk_id")
	if result != want {
		t.Errorf("got %s but expected %s", result, want)
	}
}
