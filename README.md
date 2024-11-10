# ScrapeAI

ScrapeAI is a Go library that integrates web scraping capabilities with OpenAI's GPT models, allowing for intelligent data extraction and processing from web pages.

The aim is to provide a more flexible and robust means of extracting web content than can be achieved with traditional web scraping tools. By leveraging AI, ScrapeAI can handle complex layouts, dynamic content, and context-dependent information more effectively.

I have written this library initially to help me reliably scrape data from my own personal projects, however pulled it out into a standalone library in case it might be useful or interesting to others. I have thoughts (see the Future Work section) on how to take this library further and plan to put a bit more time into it when I have some spare time.

## Features

- Intelligent web scraping powered by GPT models
- Flexible content extraction from various web page structures
- Easy integration with Go projects

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
    scraper := scrapeai.NewScraper()
    result, err := scraper.ExtractContent("https://example.com", "Extract the main heading and first paragraph")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Println(result)
}
```

For more detailed examples, check the `examples` directory. You can run the basic usage example with:

```bash
go run examples/basic_usage.go
```

### Advanced Usage

#### Custom Schema Construction

While ScrapeAI provide a default schema which returns a list of strings, which will work for many use cases, you can define custom schemas for specific data extraction needs. However, working with custom schemas requires careful attention to OpenAI's JSON Schema requirements and can be tricky to get right. We strongly recommend reviewing the official documentation before implementing custom schemas.

If you do need a custom schema, here's a basic example:

```json
{
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
```

Key requirements:
- All object schemas must include `"additionalProperties": false`
- Properties should be explicitly defined
- Use `"required"` to specify mandatory fields

For detailed information about JSON Schema support and requirements, refer to OpenAI's [Function Calling API documentation](https://platform.openai.com/docs/guides/function-calling) and [JSON Schema specification](https://json-schema.org/understanding-json-schema/).

For more detailed examples, check the `examples` directory.

## Testing

While our test coverage is currently limited, we are actively working on improving it. There are some integration tests in the `tests/` folder. You can run these tests with:

```bash
go test ./tests -v
```

## Future Work

Current plans for future work include:

- Expanding the test suite
- Adding more examples
- Enhancing the Scrape function with more flexibility and configuration options:
    - Allowing input of either URL or HTML for more user control over scraping and filtering
    - Configurable GPT model parameters (e.g., temperature) for better output control
    - Custom prompt input for tailored output
    - Support for different structured output formats beyond JSON
- Exploring methods to improve content extraction accuracy and reliability

## License

ScrapeAI is released under the [MIT License](LICENSE).
