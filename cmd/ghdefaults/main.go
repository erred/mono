package main

import (
	"go.seankhliao.com/mono/cmd/ghdefaults/internal/ghdefaults"
	"go.seankhliao.com/mono/internal/httpsvc"
)

func main() {
	httpsvc.Run(&ghdefaults.Server{})
}
