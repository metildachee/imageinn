#!/bin/bash
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <root directory> of where bulk is"
    exit 1
fi

DIRECTORY=$1

directory="$DIRECTORY"/bulk

url="http://localhost:9200/_bulk"

for file in "$directory"/*.json; do
    echo "Processing file $file"
    curl -H "Content-Type: application/x-ndjson" -X POST "$url" --data-binary "@$file"
    echo " - File $file processed."
done

echo "All files have been processed."