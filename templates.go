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

// Called from the describe function. Makes a simple HTML page with links to the body's data.
const simpleBodyOverviewPage = `<html>
  <head>
    <title>{{.Name}}</title>
		<style>
			body        {background-color: white;}
			title       {font-weight: bold;}
			p           {color: red;}
            .authname   {font-weight: bold;}
			.warning    {color: red;
                         font-weight: bold;
                         background-color: white;
						}
		</style>
  </head>
  <body>
    <ul>
		<li class="authname"><a href="{{.HomePageURL}}">{{.Name}}</a></li>
        <ul>
          {{if .IsDefunct}}<li class"warning">THIS ORGANISATION IS DEFUNCT</li>{{- end}}
		  <li><a href="{{.HomePageURL}}">Home Page</a></li>
		  <li><a href="mailto:{{.EmailGeneral}}">General Email</a></li>
		  <li><a href="mailto:{{.FOIEmailAddress}}">FOI Email</a></li>
		  <li><a href="{{.WDTKOrgPageURL}}">WDTK Page</a></li>
		  <li><a href="{{.WDTKAtomFeedURL}}">Atom Feed</a></li>
		  <li><a href="{{.WDTKJSONFeedURL}}">Updates Feed (JSON)</a></li>
		  <li><a href="{{.WDTKOrgJSONURL}}">Metadata as JSON</a></li>
		  {{if .PublicationSchemeURL}}<li><a href="{{.PublicationSchemeURL}}">Publication Scheme</a></li>{{- end}}	
		  {{if .DisclosureLogURL}}<li><a href="{{.DisclosureLogURL}}">Disclosure Log</a></li>{{- end}}
          {{if .LoCAuthorityID}}<li><a href="https://id.loc.gov/authorities/names/{{.LoCAuthorityID}}">LoC Authority</a></li>{{- end}}
          {{if .DataProtectionRegistrationIdentifier}}<li><a href="https://ico.org.uk/ESDWebPages/Entry/{{.DataProtectionRegistrationIdentifier}}">DPA Registration</a></li>{{- end}}
          {{if .WikiDataIdentifier}}<li><a href="https://www.wikidata.org/wiki/{{.WikiDataIdentifier}}">WikiData Page</a></li>{{- end}}
          {{if .TelephoneGeneral}}<li>General Telephone Number: {{.TelephoneGeneral}}</li>{{- end}}
          {{if .TelephoneFOI}}<li>FOI Telephone Number: {{.TelephoneFOI}}</li>{{- end}}
          {{if .PostalAddress}}<li>Postal Address: {{.PostalAddress}}</li>{{- end}}
    	</ul>
      </ul>
  </body>
</html>
`
