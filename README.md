## 示例代码
```
package main

import (
	"config-demo/config"
	"fmt"
)

func main() {
	filename := "config/testdata/more.ini"
	env := "develop"

	conf, err := config.Load(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(conf.Read(env, "name"))
	fmt.Println(conf.Read(env, "name1"))
}

```