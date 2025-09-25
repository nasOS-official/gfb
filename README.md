
# Go Frame Buffer 

[![MPL 2.0 License](https://img.shields.io/badge/License-MPL%202.0-green.svg)](https://www.mozilla.org/en-US/MPL/2.0/)

Go Frame Buffer is a library for working with framebuffer in tty. For example, you can output vector images in tty using this library.


## Demo

```go
package main

import "github.com/nasOS-official/gfb"

func main(){
    var fb gfb.Gfb
	fb.InitFb(0)
	fb.ClearScreen()
	fb.DrawRectangle(40, 200, 50, 100, 255, 0, 0)
	fb.DrawCircle(300, 300, 100, 0, 255, 0)
	fb.DrawLine(100, 400, 200, 500, 0, 0, 255)
	fb.UpdateScreen()
}
```

## Documentation

[Documentation on pkg.go.dev](https://pkg.go.dev/github.com/nasOS-official/gfb)



## License

This project is licensed under [MPL 2.0](LICENSE)

