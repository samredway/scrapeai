# ScrapeAI

ScrapeAI is a Go library that integrates web scraping capabilities with OpenAI's GPT models, allowing for intelligent data extraction and processing from web pages.

The aim is to provide a more flexible and robust means of extracting web content than can be achieved with traditional web scraping tools. By leveraging AI, ScrapeAI can handle complex layouts, dynamic content, and context-dependent information more effectively.

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
  - [Basic Usage](#basic-usage)
  - [Advanced Usage](#advanced-usage)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)

## Features

- Intelligent web scraping powered by GPT models
- Flexible content extraction from various web page structures
- Easy integration with Go projects

## Prerequisites

- Go 1.23 or higher
- OpenAI API key

## Installation

To install ScrapeAI, use `go get`:

```bash
go get github.com/samredway/scrapeai
```

## Setting up OpenAI API Key

You will need to set up an OpenAI API key. You can get one [here](https://platform.openai.com/account/api-keys).

Once you have an API key, you can set it in your environment variables:

```bash
export OPENAI_API_KEY=<your_api_key>
```

Alternatively, you can use a .env file and the `godotenv` package to load the environment variables:

```go
import "github.com/joho/godotenv"

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}
```

## Usage

### Basic Usage

Here's a simple example of how to use ScrapeAI:

```go
package main

import (
    "fmt"
    "github.com/samredway/scrapeai"
)

func main() {
    url := "https://example.com"
    req, err := scrapeai.NewScrapeAiRequest(url, "Extract the main headline")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    result, err := scrapeai.Scrape(req)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Println(result.Results)
}

// Output:
// {"data":["Example Domain"]}
```

For more examples, check the `examples` directory:
- `examples/basic/` - Shows basic static and dynamic HTML scraping
- `examples/advanced/` - Demonstrates custom schema usage

You can run the examples with:

```bash
# Basic static and dynamic scraping
go run examples/basic/main.go

# Advanced custom schema usage
go run examples/advanced/main.go
```

### Advanced Usage

#### Custom Schema Construction

While ScrapeAI provides a default schema which returns an object with an array of strings which will look like this on return:

```json
{
    "data": [
        "string1",
        "string2"
    ]
}
```

This format will cover many different use cases and you may never need to define your own, however at some point you may want to have more control over the shape of the return data.

When working with custom schemas, you need to follow OpenAI's JSON Schema requirements:

1. The schema must be a valid JSON Schema with `type: "object"` at the root level
2. All object schemas should include `"additionalProperties": false`
3. Properties should be explicitly defined
4. Use `"required"` to specify mandatory fields

Here's a basic example:

```json
{
    "type": "object",
    "properties": {
        "data": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "headline": {"type": "string"},
                    "body": {"type": "string"}
                },
                "additionalProperties": false,
                "required": ["headline", "body"]
            }
        }
    },
    "additionalProperties": false,
    "required": ["data"]
}
```

And the corresponding Go struct would be:
```go
struct {
    Data []struct {
        Headline string `json:"headline"`
        Body     string `json:"body"`
    } `json:"data"`
}
```

You would use your go struct to unmarshal the JSON response from the GPT model.

For detailed information about JSON Schema support and requirements, refer to OpenAI's [Function Calling API documentation](https://platform.openai.com/docs/guides/function-calling) and [JSON Schema specification](https://json-schema.org/understanding-json-schema/).

For more detailed examples, check the `examples` directory.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Testing

Tests can be found in the `tests/` directory, and can be run with:

```bash
go test ./tests -v
```

## License

ScrapeAI is released under the [MIT License](LICENSE).
