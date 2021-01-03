package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var version string = "1.0.1 ~ 03/01/2021"

func main() {

	fmt.Println("GameMaker Studio AudioGroup Extractor v" + version)
	fmt.Println("Author: Jonathan Hecl ~ https://www.jonathanhecl.com")
	fmt.Println()
	fmt.Println("USAGE: " + filepath.Base(os.Args[0]) + " audiogroup1.dat")
	fmt.Println()

	extract := true
	audiogroup := "audiogroup1.dat"
	if len(os.Args) == 2 || len(audiogroup) > 0 {
		if len(audiogroup) == 0 {
			audiogroup = os.Args[1]
		}
		if len(audiogroup) > 0 {
			started := time.Now()

			if _, err := os.Stat(audiogroup); err != nil {
				fmt.Println(audiogroup, "not found.")
				return
			}

			fmt.Println("Processing", audiogroup, "...")
			f, err := os.Open(audiogroup)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			stats, statsErr := f.Stat()
			if statsErr != nil {
				panic(statsErr)
			}
			// Get size of the file
			var fileSize int64 = stats.Size()
			fmt.Println("File size:", fileSize, "bytes")

			//Read the entire file into a buffer. We'll be doing this all in memory.
			file := bytes.NewBuffer(make([]byte, 0))
			buf := make([]byte, 10*1024)
			for {
				n, err := f.Read(buf)
				if n > 0 {
					file.Write(buf[:n])
				}
				if err == io.EOF {
					break
				}
			}
			bytesFile := file.Bytes()

			//Skip to the byte read in at 0x14.
			//That's where the first file is.
			tracks := int32(binary.LittleEndian.Uint32(bytesFile[0x10 : 0x10+4]))
			firstTrack := int32(binary.LittleEndian.Uint32(bytesFile[0x14 : 0x14+4]))
			fmt.Println(fmt.Sprintf("Number of tracks: %d tracks", tracks))
			fmt.Println(fmt.Sprintf("First Track Addr: 0x%08x", firstTrack))

			offset := firstTrack
			n := 1
			for {
				trackName := fmt.Sprintf("extract%03d.ogg", n)
				trackSize := int32(binary.LittleEndian.Uint32(bytesFile[offset : offset+4]))
				fmt.Println(fmt.Sprintf("File %s at 0x%08x (Size %d)... ", trackName, offset, trackSize))
				offset += 4
				if extract {
					//Extract the file
					fo, err := os.Create(trackName)
					if err != nil {
						panic(err)
					}
					defer fo.Close()
					if _, err := fo.Write(bytesFile[offset : offset+trackSize]); err != nil {
						panic(err)
					}
					fo.Close()
					time.Sleep(time.Microsecond)
				}

				//The file size must be at an address of a multiple of 4.
				rem := trackSize % 4
				if rem == 0 {
					rem = 4
				}
				offset += trackSize + 4 - rem
				if offset > int32(fileSize) {
					break
				}
				n++
			}

			fmt.Println("Processed in", time.Since(started).String())
		}
	}

}
