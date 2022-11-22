# gormlogger
A Logger contribe to [`gorm`](https://github.com/go-gorm/gorm). Although `gorm` provide a logger interface for user, it's not so convenient as expected. 

This repository intends to provide a more convenient way to implement gorm logger interface.


## Install
```bash
go get github.com/lawyzheng/gormlogger
```

## Quick Start
to work with [`logrus`](https://github.com/sirupsen/logrus)

```go

import (
	"github.com/lawyzheng/gormlogger"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func main(){
	// to use logrus.Logger
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.NewLogger(logrus.New()),
	})
	if err != nil {
		panic(err)
	}

	// or to use logrus.Entry
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormlogger.NewLogger(logrus.NewEntry(logrus.New())),
	})
	if err != nil {
		panic(err)
	}
}
```
