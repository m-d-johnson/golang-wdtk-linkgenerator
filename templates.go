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
      <div class="authname"><li>{{.Name}}</li></div>
        <ul>
          {{if .IsDefunct}}<div class"warning""><li>THIS ORGANISATION IS DEFUNCT</li></div>{{- end}}
		  <li><a href="{{.HomePageURL}}">Home Page</a></li>
		  <li><a href="mailto:{{.FOIEmailAddress}}">FOI Email</a></li>
		  <li><a href="{{.WDTKOrgPageURL}}">WDTK Page</a></li>
		  <li><a href="{{.WDTKAtomFeedURL}}">Atom Feed</a></li>
		  <li><a href="{{.WDTKJSONFeedURL}}">Updates Feed (JSON)</a></li>
		  <li><a href="{{.WDTKOrgJSONURL}}">Metadata as JSON</a></li>
		  {{if .PublicationSchemeURL}}<li><a href="{{.PublicationSchemeURL}}">Publication Scheme</a></li>{{- end}}		  
		  {{if .DisclosureLogURL}}<li><a href="{{.DisclosureLogURL}}">Disclosure Log</a></li>{{- end}}
          <li><a href="https://id.loc.gov/authorities/names/{{.LoCAuthorityID}}">LoC Authority</a></li>
          <li><a href="https://ico.org.uk/ESDWebPages/Entry/{{.WDTKID}}">DPA Registration</a></li>
          <li><a href="https://www.wikidata.org/wiki/{{.WikiDataIdentifier}}">WikiData Page</a></li>
    	</ul>
      </ul>
  </body>
</html>
`
