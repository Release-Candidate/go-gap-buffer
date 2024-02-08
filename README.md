# go-gap-buffer

[![golangci-lint](https://github.com/Release-Candidate/go-gap-buffer/actions/workflows/lint.yml/badge.svg)](https://github.com/Release-Candidate/go-gap-buffer/actions/workflows/lint.yml)
[![Test](https://github.com/Release-Candidate/go-gap-buffer/actions/workflows/test.yml/badge.svg)](https://github.com/Release-Candidate/go-gap-buffer/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/Release-Candidate/go-gap-buffer/graph/badge.svg?token=VCMMINZQF9)](https://codecov.io/gh/Release-Candidate/go-gap-buffer)
[![GitHub Release](https://img.shields.io/github/v/release/Release-Candidate/go-gap-buffer)](https://github.com/Release-Candidate/go-gap-buffer/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/Release-Candidate/go-gap-buffer.svg)](https://pkg.go.dev/github.com/Release-Candidate/go-gap-buffer)

A gap buffer implementation in Go.

- [Usage](#usage)
  - [Example](#example)
- [Version History](#version-history)
- [License](#license)

## Usage

```go
import gap "https://github.com/Release-Candidate/go-gap-buffer"
```

Below is a short example, detailed documentation can be found at [pkg.go.dev](https://pkg.go.dev/github.com/Release-Candidate/go-gap-buffer)

### Example

```go
import (
    "fmt"

    gap "github.com/Release-Candidate/go-gap-buffer"
)

 // Create a new, empty gap buffer.
 gapBuffer := gap.New()

 // Insert "Hello, World!" at the start of the buffer.
 gapBuffer.Insert("Hello, World!")

 // Print the content of the gap buffer as a single string.
 fmt.Println(gapBuffer.String())
 // Output:
 // Hello, World!

 // This does the same in a single step.
 gapBuffer = gap.NewStr("Hello, World!")
 fmt.Println(gapBuffer.String())
 // Output:
 // Hello, World!

 // Move 6 Unicode runes to the left, before "World!"
 gapBuffer.LeftMv()
 gapBuffer.LeftMv()
 gapBuffer.LeftMv()
 gapBuffer.LeftMv()
 gapBuffer.LeftMv()
 gapBuffer.LeftMv()

 // We can also get the content of the gap buffer as a pair of strings, one
 // to the left of the "cursor" and one to the right.
 l, r := gapBuffer.StringPair()
 fmt.Printf("left: '%s' |cursor| right: '%s'\n", l, r)
 // Output:
 // left: 'Hello, ' |cursor| right: 'World!'

 fmt.Printf("%s<|>%s\n", l, r)
 // Output:
 // Hello, <|>World!

 // From now on, "<|>" marks the current "cursor" position in the output.

 // Insert a Unicode smiley.
 gapBuffer.Insert("🙂")
 l, r = gapBuffer.StringPair()
 fmt.Printf("%s<|>%s\n", l, r)
 // Output:
 // Hello, 🙂<|>World!

 // Delete the Unicode smiley with a single `backspace` (delete the Unicode
 // Rune to the left of the cursor).
 gapBuffer.LeftDel()
 l, r = gapBuffer.StringPair()
 fmt.Printf("%s<|>%s\n", l, r)
 // Output:
 // Hello, <|>World!

 // Insert the string "funny" in a line on its own at the current cursor
 // location.
 gapBuffer.Insert("\nfunny\n")
 l, r = gapBuffer.StringPair()
 fmt.Printf("%s<|>%s\n", l, r)
 // Output:
 // Hello,
 // funny
 // <|>World!

 // Move the cursor up two lines.
 gapBuffer.UpMv()
 gapBuffer.UpMv()
 l, r = gapBuffer.StringPair()
 fmt.Printf("%s<|>%s\n", l, r)
 // Output:
 // <|>Hello,
 // funny
 // World!

 // Two runes to the right and down two lines again.
 gapBuffer.RightMv()
 gapBuffer.RightMv()
 gapBuffer.DownMv()
 gapBuffer.DownMv()
 l, r = gapBuffer.StringPair()
 fmt.Printf("%s<|>%s\n", l, r)
 // Output:
 // Hello,
 // funny
 // Wo<|>rld!
```

In the directory [./example](./example) there is this simple example of how to use the gap buffer.

To run it, use

```shell
go run ./example
```

## Version History

The latest release information is a [latest release](https://github.com/Release-Candidate/go-gap-buffer/releases/latest)

See file [./CHANGELOG](./CHANGELOG.md).

## License

This library is licensed under the MIT License, see file [./LICENSE](./LICENSE).
