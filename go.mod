module github.com/lacework/go-sdk

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.2.12
	github.com/BurntSushi/toml v0.3.1
	github.com/Netflix/go-expect v0.0.0-20200312175327-da48e75238e2
	github.com/briandowns/spinner v1.13.0
	github.com/fatih/color v1.12.0
	github.com/fatih/structs v1.1.0
	github.com/google/btree v1.0.1 // indirect
	github.com/hinshun/vt10x v0.0.0-20180809195222-d55458df857c
	github.com/hokaccha/go-prettyjson v0.0.0-20190818114111-108c894c2c0e
	github.com/honeycombio/libhoney-go v1.15.4
	github.com/kr/pty v1.1.8 // indirect
	github.com/kyokomi/emoji/v2 v2.2.8
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/peterbourgon/diskv/v3 v3.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.18.1
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/kr/pty => github.com/creack/pty v1.1.7
