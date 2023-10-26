Write-Output 'Deleting files from the output/ directory.'

Remove-Item -Verbose output/*.html
Remove-Item -Verbose output/*.htm
Remove-Item -Verbose output/*.txt
Remove-Item -Verbose output/*.md

Remove-Item -Verbose all-authorities.csv
