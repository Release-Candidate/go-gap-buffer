# go-gap-buffer

A gap buffer implementation in Go.

- [Usage](#usage)
  - [Example](#example)
- [License](#license)

## Usage

```go
import gap "https://github.com/Release-Candidate/go-gap-buffer"
```

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

 // This does the same in a single step.
 gapBuffer = gap.NewStr("Hello, World!")
 fmt.Println(gapBuffer.String())

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
 fmt.Printf("%s<|>%s\n", l, r)

 // From now on, "<|>" marks the current "cursor" position in the output.

 // Insert a Unicode smiley.
 gapBuffer.Insert("🙂")
 l, r = gapBuffer.StringPair()
 fmt.Printf("%s<|>%s\n", l, r)

 // Delete the Unicode smiley with a single `backspace` (delete the Unicode
 // Rune to the left of the cursor).
 gapBuffer.LeftDel()
 l, r = gapBuffer.StringPair()
 fmt.Printf("%s<|>%s\n", l, r)

 // Insert the string "funny" in a line on its own at the current cursor
 // location.
 gapBuffer.Insert("\nfunny\n")
 l, r = gapBuffer.StringPair()
 fmt.Printf("%s<|>%s\n", l, r)

 // Move the cursor up one line.
 gapBuffer.UpMv()
 l, r = gapBuffer.StringPair()
 fmt.Printf("%s<|>%s\n", l, r)

 // And down two lines.
 gapBuffer.DownMv()
 gapBuffer.DownMv()
 l, r = gapBuffer.StringPair()
 fmt.Printf("%s<|>%s\n", l, r)
```

In the directory [./example](./example) there is a simple example of how to use the gap buffer with a cli prompt.

To run it, use

```shell
go run ./example
```

## License

This library is licensed under the MIT License, see file [./LICENSE](./LICENSE).
