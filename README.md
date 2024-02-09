# WhatDoTheyKnow data munge tool: Generate overviews of certain organisations that [WhatDoTheyKnow](https://whatdotheyknow.com) monitors

## Foreword

This is not a well-engineered piece of software - it was assembled organically (read: scope creep) as the needs of a now-complete project evolved.

It was also used as an excuse to learn some more Go outside of work and to develop some other skills elsewhere.

For these reasons, please be assured that I do know better than to have code like this in production :) This tool should be:

- A library that does some operations on data, and
- Some tools that either collect and assemble data into JSON, or
- use the data from that JSON to do other things with it.

## Purpose

This repo contains a tool that works with data provided by [MySociety](https://www.mysociety.org/). Among other things, it generates an overview of police forces in the United Kingdom.

There is also some code in here that can generate overview tables of the entire WDTK list of public bodies, or a subset of them.

It also generates an OPML file which can be imported to a feed reader to monitor all UK Police FOI requests and updates tracked by WDTK.

This tool was written for a specific project, and will likely not be maintained.

This is why it has odd functionality and why some parts are messy - it's very much a toy project and I have used it as a way to learn things. There are no warranties and while I'll try to make fixes, they may take some time.

## Usage

- `go build .`

## Output Files

- output/overview.md: Overview of UK Police Forces.
- output/police.opml: OPML files containing RSS feeds for all forces in the generated table.
- output/all-mysociety.md: Simple table of all public bodies on WhatDoTheyKnow.

## Input Files

- data/foi-emails.json: Maps FOI email addresses to WDTK organisations by their 'URL Names'. Manually curated and
  required to regenerate the `generated-dataset.json` file.
- data/manual.json: Manually curated file containing additional information to augment what we can get from WDTK.
- data/police.json: Data from WikiData on UK Police Forces in JSON format.
- data/wdtk-police.csv: Mapping of full names of police forces to the `URL Name` (Unmaintained).
- data/wikidata-police-forces.json: Simple JSON mapping WikiData IDs to WDTK `URL Name`s (Unmaintained).
- output/wikidata-localpolice.csv: From WikiData - Mapping homepages to WDTK `URL Name`s (Unmaintained).
- data/whatdotheyknow_authorities_dataset.csv - created by exporting from the Excel sheet MySociety makes available.

## Information derived from the WDTK 'URL Name'

Much of the information listed in the tables is simply derived from the unique ID that WDTK uses - a string they call the `URL Name` as it's used in the URL. It's a unique identifier as far as I can tell, or at least has been for this purpose. A simple substitution yields the correct information for:

- The body's page on the WDTK site: `https://www.whatdotheyknow.com/body/{name_of_org}`
- The Atom feed of requests, updates, etc for a particular body: `https://www.whatdotheyknow.com/feed/body/{name_of_org}`
- A JSON representation of metadata about the body: `https://www.whatdotheyknow.com/feed/body/{name_of_org}.json`

## Other information

- Tags are available from the JSON data available per-body and also in the CSV that's downloaded from MySociety.
- FOI Emails are not programmatically available and are kept in a file.

# A note on testing

You'll note that there aren't tests. I am aware that good software engineering means testing your code is mandatory, but also aware that this code isn't really meant to be good software engineering as its very narrow use case means that errors are tolerable and the time spent creating tests needed to be spent organising data :)
