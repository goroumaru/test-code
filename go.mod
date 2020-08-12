module github.com/goroumaru/test-code

go 1.13

require (
	github.com/google/logger v1.1.0
	github.com/montanaflynn/stats v0.6.3
	github.com/thrasher-corp/gocryptotrader v0.0.0-20200804040818-0e30756e353e
	go.uber.org/zap v1.15.0
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	gonum.org/v1/plot v0.7.0
	gopkg.in/ini.v1 v1.57.0
)

replace github.com/thrasher-corp/gocryptotrader => ../gocryptotrader
