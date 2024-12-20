package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/byebyebruce/wadu/internal/biz"
	"github.com/byebyebruce/wadu/internal/client"
	"github.com/byebyebruce/wadu/internal/dao"
	"github.com/byebyebruce/wadu/internal/server"
	"github.com/byebyebruce/wadu/pdfbook"
	"github.com/byebyebruce/wadu/tts/volcanotts"
	"github.com/byebyebruce/wadu/vlm"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
	godotenv.Overload()

	rootCMD := serverCMD()
	rootCMD.AddCommand(
		uploadCMD(),
		promptCMD(),
	)
	if err := rootCMD.Execute(); err != nil {
		panic(err)
	}
}

func serverCMD() *cobra.Command {
	cmd := &cobra.Command{}
	var (
		addr       = cmd.Flags().String("addr", ":8081", "server address")
		db         = cmd.Flags().String("db", "db.bolt", "db file")
		assets     = cmd.Flags().String("assets", "assets", "assets folder path")
		promptFile = cmd.Flags().String("prompt", "", "vlm prompt file")
	)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		openaiCli, err := vlm.NewClientFromEnv()
		if err != nil {
			panic(err)
		}

		d, err := dao.New(*db)
		if err != nil {
			return err
		}
		defer d.Close()
		if err := os.MkdirAll(*assets, 0755); err != nil {
			return err
		}

		tts, err := volcanotts.NewTTSWithEnv()
		if err != nil {
			return err
		}

		if len(*promptFile) > 0 {
			b, err := os.ReadFile(*promptFile)
			if err != nil {
				return err
			}
			pdfbook.Prompt = string(b)
		}
		b := biz.NewBiz(d, openaiCli, tts, *assets)
		w := server.NewServer(b, *assets)
		err = w.Run(*addr)
		if err != nil {
			return err
		}
		return nil
	}
	return cmd
}

// 命令行上传
func uploadCMD() *cobra.Command {
	cmd := cobra.Command{
		Use: "upload",
	}
	var (
		server = cmd.Flags().String("server", "http://localhost:8081", "server address")
	)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		openaiCli, err := vlm.NewClientFromEnv()
		if err != nil {
			panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
		defer cancel()
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()

		b, err := pdfbook.GenFromPDF(ctx, openaiCli.Client, openaiCli.Model, f)
		if err != nil {
			return err
		}
		fmt.Println()
		fmt.Println(b.Title)
		for i, p := range b.Pages {
			fmt.Println("page", i)
			for _, s := range p.Sentences {
				fmt.Println("", "sentence", s)
			}
			fmt.Println()
		}
		fmt.Println("输入回车上传")
		//用户输入回车
		fmt.Scanln()

		err = client.PostRawBook(ctx, *server, b)
		if err != nil {
			return err
		}
		fmt.Println("上传成功")
		return err
	}
	return &cmd
}

// prompt
func promptCMD() *cobra.Command {
	cmd := cobra.Command{
		Use: "prompt",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Println(pdfbook.Prompt)
		return nil
	}
	return &cmd
}
