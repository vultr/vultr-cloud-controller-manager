# metadata

A Go client that interacts with the [Vultr Metadata](https://www.vultr.com/metadata/).

## Installation

`go get -u github.com/vultr/metadata`

## Usage

Currently, there is only one available call `Metadata()` which will retrieve your entire metadata from your instance. If you want to retrieve a specific of your metadata you can do so by calling the corresponding exported field on the `metadata` struct.

```go
c := metadata.NewClient()

meta, err := c.Metadata()
if err != nil {
	fmt.Println(err)
}

fmt.Println(meta)
fmt.Println(meta.InstanceID) // will print your instance-id
```