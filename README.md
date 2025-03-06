# go-template

Simple CLI which allows to render [Go templates](https://pkg.go.dev/text/template) and pass data from JSON or YAML into it.

## Install
```
go install github.com/dvob/go-template@latest
```

## Usage
Without data:
```
go-template my-template.tmpl
```

With data and data format determined from file extension:
```
# json
go-template -d myfile.json my-template.tmpl

# yaml
go-template -d myfile.yaml my-template.tmpl
```

Specify data format:
```
go-template -d myfile -f json my-template.tmpl
```
