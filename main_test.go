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
	"fmt"
	"log"
	"os"
	"testing"
)

// TestCleanupWithRetainFalse calls main.Cleanup with a bool retain=false after creating a file
//
//	called `output/all-authorities.csv`. The expected behaviour is that the file will be deleted.
//
// The test must fail if the file isn't deleted.
func TestCleanupWithRetainFalse(t *testing.T) {
	// Create output/ directory
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		err := os.Mkdir("output", 0755)
		if err != nil {
			fmt.Println("Error creating an output/ directory: ", err)
		}
	}

	// Create all-authorities.csv file
	if _, err := os.Stat("output"); !os.IsNotExist(err) {
		f, err := os.Create("output/all-authorities.csv")
		if err != nil {
			log.Fatalln("Failed to create file output/all-authorities.csv", err)
		}
		f.Close()
	}

	// Call function under test
	Cleanup(false)

	// Check if file still exists
	if _, e := os.Stat("output/all-authorities.csv"); !os.IsNotExist(e) {
		t.Errorf("error path output/all-authorities still exists: %v", e)
	}
}

func TestCleanupWithRetainTrue(t *testing.T) {
	// Create output/ directory
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		err := os.Mkdir("output", 0755)
		if err != nil {
			fmt.Println("Error creating an output/ directory: ", err)
		}
	}

	// Create all-authorities.csv file
	if _, err := os.Stat("output"); !os.IsNotExist(err) {
		f, err := os.Create("output/all-authorities.csv")
		if err != nil {
			log.Fatalln("Failed to create file output/all-authorities.csv", err)
		}
		f.Close()
	}

	// Call function under test
	Cleanup(true)

	// Check if file still exists
	if _, e := os.Stat("output/all-authorities.csv"); os.IsNotExist(e) {
		t.Errorf("error path output/all-authorities.csv does not exist: %v", e)
	}

	// Cleaning up the file that we retained
	os.Remove("output/all-authorities.csv")
}
