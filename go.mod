module trade-optimizer

go 1.15

replace github.com/c-bata/goptuna => ./goptuna

replace github.com/frankrap/talib => ./talib

require (
	github.com/c-bata/goptuna v0.8.1
	github.com/dmitryikh/leaves v0.0.0-20210121075304-82771f84c313
	github.com/frankrap/talib v0.0.0-20210213031721-36337d393468
	github.com/montanaflynn/stats v0.6.6
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	gorm.io/driver/mysql v1.0.3
	gorm.io/gorm v1.20.12
)
