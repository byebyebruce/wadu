package edgetts

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/websocket"
)

type turnContext struct {
	ServiceTag string `json:"serviceTag"`
}

type turnAudio struct {
	Type     string `json:"type"`
	StreamID string `json:"streamId"`
}

// {"context": {"serviceTag": "743d56a9126e4649b2af1660975e3520"}}
type turnStart struct {
	Context turnContext `json:"context"`
}

// {"context":{"serviceTag":"743d56a9126e4649b2af1660975e3520"},"audio":{"type":"inline","streamId":"8D6F2A03213641159BE3476B36473521"}}
type turnResp struct {
	Context turnContext `json:"context"`
	Audio   turnAudio   `json:"audio"`
}

type turnMetaInnerText struct {
	Text         string `json:"Text"`
	Length       int    `json:"Length"`
	BoundaryType string `json:"BoundaryType"`
}

type turnMetaInnerData struct {
	Offset   int               `json:"Offset"`
	Duration int               `json:"Duration"`
	Text     turnMetaInnerText `json:"text"`
}

type turnMetadata struct {
	Type string            `json:"Type"`
	Data turnMetaInnerData `json:"Data"`
}

type turnMeta struct {
	Metadata []turnMetadata `json:"Metadata"`
}

type communicateChunk struct {
	Type     string
	Data     []byte
	Offset   int
	Duration int
	Text     string
	Error    error
}

type CommunicateTextTask struct {
	//id     int
	//text   string
	option Option

	//chunk      chan communicateChunk
	//speechData []byte
}

type Option struct {
	voice  string
	rate   string
	volume string
}

func openWs() (*websocket.Conn, error) {
	headers := http.Header{}
	headers.Add("Pragma", "no-cache")
	headers.Add("Cache-Control", "no-cache")
	headers.Add("Origin", "chrome-extension://jdiccldimpdaibmpdkjnbmckianbfold")
	headers.Add("Accept-Encoding", "gzip, deflate, br")
	headers.Add("Accept-Language", "en-US,en;q=0.9")
	headers.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36 Edg/91.0.864.41")

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(fmt.Sprintf("%s&ConnectionId=%s", WSS_URL, uuidWithOutDashes()), headers)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func TTS(ctx context.Context, str string, writer io.Writer, option Option) error {
	//text.chunk = make(chan communicateChunk)
	// texts := splitTextByByteLength(removeIncompatibleCharacters(c.text), calcMaxMsgSize(c.voice, c.rate, c.volume))
	conn, err := openWs()
	if err != nil {
		return err
	}
	defer conn.Close()
	date := dateToString()
	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("X-Timestamp:%s\r\nContent-Type:application/json; charset=utf-8\r\nPath:speech.config\r\n\r\n{\"context\":{\"synthesis\":{\"audio\":{\"metadataoptions\":{\"sentenceBoundaryEnabled\":false,\"wordBoundaryEnabled\":true},\"outputFormat\":\"audio-24khz-48kbitrate-mono-mp3\"}}}}\r\n", date)))
	conn.WriteMessage(websocket.TextMessage, []byte(ssmlHeadersPlusData(uuidWithOutDashes(), date, mkssml(
		str, option.voice, option.rate, option.volume,
	))))

	// download indicates whether we should be expecting audio data,
	// this is so what we avoid getting binary data from the websocket
	// and falsely thinking it's audio data.
	downloadAudio := false

	// audio_was_received indicates whether we have received audio data
	// from the websocket. This is so we can raise an exception if we
	// don't receive any audio data.
	// audioWasReceived := false

	chunk := make(chan communicateChunk)
	// finalUtterance := make(map[int]int)
	go func() error {
		defer close(chunk)
		for {
			// 读取消息
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return err
			}
			switch messageType {
			case websocket.TextMessage:
				parameters, data, _ := getHeadersAndData(data)
				path := parameters["Path"]
				if path == "turn.start" {
					downloadAudio = true
				} else if path == "turn.end" {
					downloadAudio = false
					select {
					case <-ctx.Done():
						return nil
					case chunk <- communicateChunk{Type: ChunkTypeEnd}:
					}
					return nil
				} else if path == "audio.metadata" {
					meta := &turnMeta{}
					if err := json.Unmarshal(data, meta); err != nil {
						log.Fatalf("We received a text message, but unmarshal failed.")
					}
					for _, v := range meta.Metadata {
						if v.Type == ChunkTypeWordBoundary {
							cc := communicateChunk{
								Type:     v.Type,
								Offset:   v.Data.Offset,
								Duration: v.Data.Duration,
								Text:     v.Data.Text.Text,
							}
							select {
							case <-ctx.Done():
								return nil
							case chunk <- cc:
							}
						} else if v.Type == ChunkTypeSessionEnd {
							continue
						} else {
							log.Fatalf("Unknown metadata type: %s", v.Type)
						}
					}
				} else if path != "response" {
					log.Fatalf("The response from the service is not recognized.\n%s", data)
				}
			case websocket.BinaryMessage:
				if !downloadAudio {
					log.Fatalf("We received a binary message, but we are not expecting one.")
				}
				if len(data) < 2 {
					log.Fatalf("We received a binary message, but it is missing the header length.")
				}
				headerLength := int(binary.BigEndian.Uint16(data[:2]))
				if len(data) < headerLength+2 {
					log.Fatalf("We received a binary message, but it is missing the audio data.")
				}
				cc := communicateChunk{
					Type: ChunkTypeAudio,
					Data: data[headerLength+2:],
				}
				select {
				case chunk <- cc:
				case <-ctx.Done():
					return nil
				}
				// audioWasReceived = true
			}
		}
	}()

	for v := range chunk {
		if v.Type == ChunkTypeAudio {
			if _, err := writer.Write(v.Data); err != nil {
				return err
			}

			//t.speechData = append(t.speechData, v.Data...)
			// } else if v.Type == ChunkTypeWordBoundary {
		} else if v.Type == ChunkTypeEnd {
			//close(t.chunk)
			return nil
		}
	}
	return nil

}

func isValidVoice(voice string) bool {
	return regexp.MustCompile(`^Microsoft Server Speech Text to Speech Voice \(.+,.+\)$`).MatchString(voice)
}

func isValidRate(rate string) bool {
	if rate == "" {
		return false
	}
	return regexp.MustCompile(`^[+-]\d+%$`).MatchString(rate)
}

func isValidVolume(volume string) bool {
	if volume == "" {
		return false
	}
	return regexp.MustCompile(`^[+-]\d+%$`).MatchString(volume)
}
