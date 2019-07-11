package main

import (
  "unicode"
  "bytes"
  "github.com/ghodss/yaml"
)




//----------------------------------------------------------------
//Allow painless Ingesting of YAML
//----------------------------------------------------------------
func ToJSON(data []byte) ([]byte, error) {
    if hasJSONPrefix(data) {
        return data, nil
    }
    return yaml.YAMLToJSON(data)
}

var jsonPrefix = []byte("{")

// hasJSONPrefix returns true if the provided buffer starts with "{".
func hasJSONPrefix(buf []byte) bool {
    return hasPrefix(buf, jsonPrefix)
}

// Return true if the first non-whitespace bytes in buf is prefix.
func hasPrefix(buf []byte, prefix []byte) bool {
    trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
    return bytes.HasPrefix(trim, prefix)
}


