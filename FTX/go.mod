module github.com/goroumaru/test-code/FTX

go 1.13

require (
	github.com/montanaflynn/stats v0.6.3
	github.com/thrasher-corp/gocryptotrader v0.0.0-20200804040818-0e30756e353e
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	gonum.org/v1/netlib v0.0.0-20200603212716-16abd5ac5bc7 // indirect
	gonum.org/v1/plot v0.7.0
	gopkg.in/ini.v1 v1.57.0
)

replace github.com/thrasher-corp/gocryptotrader => ../../gocryptotrader
