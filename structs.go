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

// Authority is a public body on the WDTK website. Here, it contains additional fields which are not
// provided by WhatDoTheyKnow. These fields are added from a manually maintained list.
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
	InspectorateProfile                  string `json:"Inspectorate_Profile"`
}

// JSONResponse is the JSON object you receive when you request details about an authority on the
// WDTK website.
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

// Record is an entry in the manually curated file which provides data that WDTK does not.
// You'll find this as data/in/collected.json
type Record struct {
	WDTKID              string `json:"wdtk_id"`
	TelephoneGeneral    string `json:"telephone_general"`
	TelephoneFOI        string `json:"telephone_foi"`
	EmailGeneral        string `json:"email_general"`
	EmailFOI            string `json:"email_foi"`
	PostalAddress       string `json:"postal_address"`
	InspectorateProfile string `json:"inspectorate_profile"`
}
