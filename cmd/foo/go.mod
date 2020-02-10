module github.com/Zenithar/go-monorepo/cmd/foo

go 1.13

replace github.com/Zenithar/go-monorepo => ../../

require (
	github.com/Zenithar/go-monorepo v0.0.0-00010101000000-000000000000
	github.com/common-nighthawk/go-figure v0.0.0-20190529165535-67e0ed34491a
	github.com/fatih/color v1.9.0
	github.com/magefile/mage v1.9.0
	github.com/spf13/cobra v0.0.5
	go.uber.org/zap v1.13.0
)
