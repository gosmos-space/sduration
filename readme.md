# sduration

[![Go Reference](https://pkg.go.dev/badge/github.com/gosmos-space/sduration.svg)](https://pkg.go.dev/github.com/gosmos-space/sduration)
[![Go Report Card](https://goreportcard.com/badge/github.com/gosmos-space/sduration)](https://goreportcard.com/report/github.com/gosmos-space/sduration)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go package that provides a string-serializable duration type that wraps `time.Duration`.

## Overview

The `sduration` package implements a custom duration type (`SDuration`) that wraps the standard `time.Duration` with serialization capabilities. It enables:

- Text marshaling/unmarshaling (for encoding like YAML, TOML)
- JSON marshaling/unmarshaling
- sql package compatibility for database storage
- Human-readable duration strings in configs and data formats

## Installation

```bash
go get github.com/gosmos-space/sduration
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/gosmos-space/sduration"
)

func main() {
    // Create an SDuration from time.Duration
    d := sduration.SDuration(5 * time.Second)
    
    // Convert back to time.Duration
    timeDur := d.Duration()
    
    fmt.Printf("Duration: %v\n", timeDur) // Output: Duration: 5s
}
```

### With Text Marshaling (YAML, TOML)

```go
package main

import (
    "fmt"
    "encoding/json"
    "gopkg.in/yaml.v3"
    
    "github.com/gosmos-space/sduration"
)

type Config struct {
    Timeout sduration.SDuration `yaml:"timeout" json:"timeout"`
    Retry   sduration.SDuration `yaml:"retry" json:"retry"`
}

func main() {
    // Parse config from YAML
    yamlData := []byte(`
timeout: 30s
retry: 5m
`)
    
    var config Config
    if err := yaml.Unmarshal(yamlData, &config); err != nil {
        panic(err)
    }
    
    fmt.Printf("Timeout: %v\n", config.Timeout.Duration())
    fmt.Printf("Retry: %v\n", config.Retry.Duration())
    
    // Convert back to YAML
    out, _ := yaml.Marshal(config)
    fmt.Println(string(out))
}
```

### With JSON

```go
package main

import (
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/gosmos-space/sduration"
)

type APIConfig struct {
    RequestTimeout sduration.SDuration `json:"requestTimeout"`
    SessionTTL     sduration.SDuration `json:"sessionTTL"`
}

func main() {
    // Create a config
    config := APIConfig{
        RequestTimeout: sduration.SDuration(30 * time.Second),
        SessionTTL:     sduration.SDuration(24 * time.Hour),
    }
    
    // Marshal to JSON
    jsonData, _ := json.Marshal(config)
    fmt.Println(string(jsonData))
    // Output: {"requestTimeout":"30s","sessionTTL":"24h0m0s"}
    
    // Unmarshal from JSON
    jsonInput := []byte(`{"requestTimeout":"45s","sessionTTL":"12h"}`)
    var newConfig APIConfig
    
    if err := json.Unmarshal(jsonInput, &newConfig); err != nil {
        panic(err)
    }
    
    fmt.Printf("Request Timeout: %v\n", newConfig.RequestTimeout.Duration())
    fmt.Printf("Session TTL: %v\n", newConfig.SessionTTL.Duration())
}
```

### With SQL Databases

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    
    "github.com/gosmos-space/sduration"
    _ "github.com/mattn/go-sqlite3"
)

type Job struct {
    ID       int
    Timeout  sduration.SDuration
    Interval sduration.SDuration
}

func main() {
    // Open database connection
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // Create table
    _, err = db.Exec(`CREATE TABLE jobs (
        id INTEGER PRIMARY KEY,
        timeout TEXT,
        interval TEXT
    )`)
    if err != nil {
        log.Fatal(err)
    }
    
    // Insert job with durations
    job := Job{
        ID:       1,
        Timeout:  sduration.SDuration(30 * time.Second),
        Interval: sduration.SDuration(5 * time.Minute),
    }
    
    _, err = db.Exec(
        "INSERT INTO jobs (id, timeout, interval) VALUES (?, ?, ?)",
        job.ID, job.Timeout, job.Interval,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Query the data back
    var retrievedJob Job
    err = db.QueryRow("SELECT id, timeout, interval FROM jobs WHERE id = ?", 1).Scan(
        &retrievedJob.ID, &retrievedJob.Timeout, &retrievedJob.Interval,
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Job ID: %d\n", retrievedJob.ID)
    fmt.Printf("Timeout: %v\n", retrievedJob.Timeout.Duration())
    fmt.Printf("Interval: %v\n", retrievedJob.Interval.Duration())
}
```

## Features

- **String Serialization**: Durations are stored as human-readable strings (e.g., "5m30s") rather than nanosecond integers
- **Compatible with time.Duration**: Easy conversion between SDuration and time.Duration
- **Supports JSON**: Properly marshals/unmarshals to/from JSON as strings with units
- **Supports Text Encoding**: Works with text-based encodings like YAML and TOML
- **SQL Database Support**: Implements `sql.Scanner` and `driver.Valuer` for database storage
- **Type Safety**: Maintains type safety while ensuring serialization works correctly

## Supported Formats

SDuration supports all the same duration formats as `time.ParseDuration`:

- "ns" - nanoseconds
- "us" or "µs" - microseconds
- "ms" - milliseconds
- "s" - seconds
- "m" - minutes
- "h" - hours

Examples: "300ms", "1.5h", "2h45m", "-30s"

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
