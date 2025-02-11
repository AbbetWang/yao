package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/yaoapp/gou/plugin"
	"github.com/yaoapp/gou/process"
	"github.com/yaoapp/kun/exception"
	"github.com/yaoapp/kun/utils"
	"github.com/yaoapp/yao/config"
	"github.com/yaoapp/yao/engine"
	"github.com/yaoapp/yao/share"
)

var runSilent = false

var runCmd = &cobra.Command{
	Use:   "run",
	Short: L("Execute process"),
	Long:  L("Execute process"),
	Run: func(cmd *cobra.Command, args []string) {
		defer share.SessionStop()
		defer plugin.KillAll()
		defer func() {
			err := exception.Catch(recover())
			if err != nil {
				if !runSilent {
					color.Red(L("Fatal: %s\n"), err.Error())
					return
				}
				fmt.Printf("%s\n", err.Error())
			}
		}()

		Boot()

		cfg := config.Conf
		cfg.Session.IsCLI = true
		if len(args) < 1 {
			if !runSilent {
				color.Red(L("Not enough arguments\n"))
				color.White(share.BUILDNAME + " help\n")
				return
			}
			fmt.Printf(L("Not enough arguments\n"))
			return
		}

		err := engine.Load(cfg)
		if err != nil {
			if !runSilent {
				color.Red(L("Engine: %s\n"), err.Error())
				return
			}

			fmt.Printf("%s\n", err.Error())
			return
		}

		name := args[0]
		if !runSilent {
			color.Green(L("Run: %s\n"), name)
		}

		pargs := []interface{}{}
		for i, arg := range args {
			if i == 0 {
				continue
			}

			// Parse the arguments
			if strings.HasPrefix(arg, "::") {
				arg := strings.TrimPrefix(arg, "::")
				var v interface{}
				err := jsoniter.Unmarshal([]byte(arg), &v)
				if err != nil {
					color.Red(L("Arguments: %s\n"), err.Error())
					return
				}
				pargs = append(pargs, v)

				if !runSilent {
					color.White("args[%d]: %s\n", i-1, arg)
				}

			} else if strings.HasPrefix(arg, "\\::") {
				arg := "::" + strings.TrimPrefix(arg, "\\::")
				pargs = append(pargs, arg)
				if !runSilent {
					color.White("args[%d]: %s\n", i-1, arg)
				}

			} else {
				pargs = append(pargs, arg)
				if !runSilent {
					color.White("args[%d]: %s\n", i-1, arg)
				}
			}

		}

		process := process.New(name, pargs...)
		res, err := process.Exec()
		if err != nil {
			if !runSilent {
				color.Red(L("Process: %s\n"), err.Error())
				return
			}
			fmt.Printf("%s\n", err.Error())
			return
		}

		if !runSilent {
			color.White("--------------------------------------\n")
			color.White(L("%s Response\n"), name)
			color.White("--------------------------------------\n")
			utils.Dump(res)
			color.White("--------------------------------------\n")
			color.Green(L("✨DONE✨\n"))
			return
		}

		// Silent mode output
		switch res.(type) {

		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
			fmt.Printf("%v\n", res)
			return

		case string, []byte:
			fmt.Printf("%s\n", res)
			return

		default:
			txt, err := jsoniter.Marshal(res)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}
			fmt.Printf("%s\n", txt)
		}
	},
}

func init() {
	runCmd.PersistentFlags().BoolVarP(&runSilent, "silent", "s", false, L("Silent mode"))
}
