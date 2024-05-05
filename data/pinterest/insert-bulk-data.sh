#!/bin/bash
directory="bulk"

url="http://localhost:9200/_bulk"

for file in "$directory"/*.json; do
    echo "Processing file $file"
    curl -H "Content-Type: application/x-ndjson" -X POST "$url" --data-binary "@$file"
    echo " - File $file processed."
done

echo "All files have been processed."