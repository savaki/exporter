# exporter

exports partner data from the specified site

### installation

```
go get github.com/savaki/exporter
```

### download all the partner data

```
exporter crawl --codebase {search-page-url} --pages 333 --dir target
```

grabs the first 333 pages of search results and saves the partner html to the target directory

### convert html to json

```
exporter partner {source-html} > {target-json}
```

exports the partner source file as a structured json file

