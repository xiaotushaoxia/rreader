Example
-------

```
package main

import (
	"bytes"
	"io"
	"os"

	"github.com/xiaotushaoxia/rreader"
)

func main() {
	rs := bytes.NewReader([]byte(
		"123454-3-2-1-",
	))

	nn := bytes.NewBuffer(nil)

	newReader := rreader.NewReader(rs)
	n, err := nn.ReadFrom(newReader)
	fmt.Println(n, err)
	fmt.Println(nn.String()) // 输出-1-2-3-454321
}
```