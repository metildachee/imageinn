```bash
python3 script.py
curl -H "Content-Type: application/x-ndjson" -XPOST "localhost:9200/_bulk" --data-binary "@pinterest_es.json"
```