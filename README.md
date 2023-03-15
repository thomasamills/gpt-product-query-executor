# gpt-product-query-executor
Given pdf files of catalog pages, get gpt to create you a multitude of html elements from extracting the data.

## Installation
```go mod download```

```go mod vendor```

```go build -o ./gpt-product-gen```

Now run the program with arg[1] set as the path to the csv and arg[2] set to your OpenAI secret key.
