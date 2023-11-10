一些说明
-------
1. 如果把rs io.ReadSeeker传给了NewReader，rs自身不应该有任何操作
2. rs已经读取的数据不会再被reverse reader读到。 如rs本身有1234567，在传给NewReader前读了12，那么反向读全部就只能读出76543
3. 不并发安全
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