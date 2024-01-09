# Guide to the files in this directory

- `collected.json` - manually maintained dataset containing:
    - Telephone numbers
    - Email addresses
    - Postal addresses
    - Links to profile on inspectorate website
- `foi-directory.csv` - from the https://foi.directory website
- `foi-emails.json` (deprecated) - maps FOI emails to force WDTK IDs.
- `generated-dataset.json` - source file for the rest of the program. Generated from various sources.
- `manual.json` (deprecated) - manually maintained JSON containing data to augment what's returned from the WDTK API.
- `police.json` - results from querying Wikidata for UK police forces.
- `wdtk-police.csv` - mapping of WDTK IDs to names of authorities.
- `wikidata-police-forces.json` - results from wikidata.
- `whatdotheyknow_authorities_dataset.xlsx` - downloaded from MySociety - a database dump.
- `whatdotheyknow_authorities_dataset.csv` - a CSV version of the above XLSX file.
