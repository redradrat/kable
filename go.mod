module github.com/redradrat/kable

go 1.14

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f
	github.com/fatih/color v1.13.0
	github.com/fatih/structs v1.1.0
	github.com/go-git/go-git/v5 v5.1.0
	github.com/gofiber/fiber/v2 v2.3.0
	github.com/gofiber/template v1.6.6
	github.com/google/go-querystring v1.0.0
	github.com/google/logger v1.1.0
	github.com/grafana/tanka v0.18.2
	github.com/jsonnet-bundler/jsonnet-bundler v0.4.0
	github.com/kr/pty v1.1.5 // indirect
	github.com/labstack/echo/v4 v4.1.16
	github.com/labstack/gommon v0.3.0
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/olekukonko/tablewriter v0.0.4
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.7.0
	go.etcd.io/etcd/client/v3 v3.5.1
)

replace github.com/Joker/jade v1.0.0 => github.com/Joker/jade v1.0.1-0.20200506134858-ee26e3c533bb

replace github.com/grafana/tanka v0.18.2 => github.com/redradrat/tanka v0.18.2-fix639
