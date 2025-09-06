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
	"github.com/byebyebruce/wadu/pkg/pdfx"
	//"github.com/byebyebruce/wadu/tts/edgetts"
	"github.com/byebyebruce/wadu/tts/openaitts"
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
		pdf2imgCMD(),
		uploadImageCMD(),
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
		debug      = cmd.Flags().Bool("debug", false, "debug mode")
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
		if err := d.InitDB(); err != nil {
			return err
		}
		if err := os.MkdirAll(*assets, 0755); err != nil {
			return err
		}

		/*
			tts, err := volcanotts.NewTTSWithEnv()
			if err != nil {
				return err
			}
		*/
		//tts := edgetts.New()
    tts:=openaitts.NewTTSFromEnv()

		if len(*promptFile) > 0 {
			b, err := os.ReadFile(*promptFile)
			if err != nil {
				return err
			}
			pdfbook.Prompt = string(b)
		}
		b := biz.NewBiz(d, openaiCli, tts, *assets)
		w := server.NewServer(b, *assets)
		err = w.Run(*addr, *debug)
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
		Use: "upload [pdf] <audio>",
	}
	server := cmd.Flags().String("server", "http://localhost:8081", "server address")
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

		var audioBytes []byte
		if len(args) > 1 {
			if b, err := os.ReadFile(args[1]); err != nil {
				return err
			} else {
				audioBytes = b
			}
		}

		b, err := pdfbook.GenFromPDF(ctx, openaiCli.Client, openaiCli.Model, f)
		if err != nil {
			return err
		}
		b.MP3 = audioBytes
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
		// 用户输入回车
		fmt.Scanln()

		uploadCtx, uploadCancel := context.WithTimeout(context.Background(), time.Minute*10)
		defer uploadCancel()
		err = client.PostRawBook(uploadCtx, *server, b)
		if err != nil {
			return err
		}
		fmt.Println("上传成功")
		return err
	}
	return &cmd
}

func uploadImageCMD() *cobra.Command {
	var (
		title string
	)
	cmd := cobra.Command{
		Use: "upimg <title> <image1> <image2> ...",
	}
	cmd.Flags().StringVar(&title, "title", "", "book title")
	server := cmd.Flags().String("server", "http://localhost:8081", "server address")
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		openaiCli, err := vlm.NewClientFromEnv()
		if err != nil {
			panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
		defer cancel()
		//title:= args[0]
		images := args[:]
		imagesBytes := make([][]byte, 0, len(images))
		for _, img := range images {
			f, err := os.ReadFile(img)
			if err != nil {
				panic(err)
			}
			imagesBytes = append(imagesBytes, f)
		}

		b, err := pdfbook.GenFromImages(ctx, openaiCli.Client, openaiCli.Model, title, imagesBytes...)
		if err != nil {
			return err
		}
		//b.MP3 = audioBytes
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
		// 用户输入回车
		fmt.Scanln()

		uploadCtx, uploadCancel := context.WithTimeout(context.Background(), time.Minute*10)
		defer uploadCancel()
		err = client.PostRawBook(uploadCtx, *server, b)
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

func pdf2imgCMD() *cobra.Command {
	cmd := cobra.Command{
		Use: "pdf2img",
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		imgs, err := pdfx.ConvertPDFFileToJPEG(args[0])
		if err != nil {
			panic(err)
		}
		for i, img := range imgs {
			f, err := os.Create(fmt.Sprintf("img-%d.jpg", i))
			if err != nil {
				panic(err)
			}
			defer f.Close()
			_, err = f.Write(img)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}
	return &cmd
}
