package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"heyuanlong/blockchain-step/node"
	_ "heyuanlong/blockchain-step/log"
	"os"
	"sort"
)



var (
	app *cli.App
)

func init() {
	app = cli.NewApp()
	app.UseShortOptionHandling = true
	//app.Action = geth
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright The dgo Authors"
	app.Commands = []*cli.Command{
		{
			Name:        "account",
			Aliases:     []string{"account"},
			Usage:       "step账户系统",
			Description: "创建/解析/校验账户信息",
			Subcommands: []*cli.Command{
				{
					Name:      "create",
					Aliases: []string{"c"},
					Usage:   "创建一个新的账户",
					UsageText: "exe account create -password 123456",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "password",
							Usage: "密码",
							Value: "123456",
						},
					},
					Action: accountCreate,
				},
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	app.Before = func(ctx *cli.Context) error {
		return nil
	}
	app.After = func(ctx *cli.Context) error {
		return nil
	}

	app.UsageText="exe -config ./config.json"
	app.Flags= []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: "Load configuration from `FILE`",
			Value: "./config.json",
		},
	}
	app.Action=nodeRun
}

func accountCreate(c *cli.Context) error {
	return nil
}


func nodeRun(c *cli.Context) error {
	confFile := c.String("config")
	fmt.Println("confFile:", confFile)

	node.New(confFile).Run()
	return nil
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}