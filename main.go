package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"

	// driver "gitlab.com/gomidi/midi/testdrv"
	driver "gitlab.com/gomidi/rtmididrv"

	"gitlab.com/gomidi/midi/writer"
)

// This example reads from the first input and and writes to the first output port
func main() {
	// drv := driver.New("fake cables: messages written to output port 0 are received on input port 0")

	drv, err := driver.New()
	checkErr(err)

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	checkErr(err)

	outs, err := drv.Outs()
	checkErr(err)

	in, out := ins[0], outs[0]

	checkErr(in.Open())
	checkErr(out.Open())

	defer in.Close()
	defer out.Close()

	// the writer we are writing to
	wr := writer.New(out)

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(
		reader.NoLogger(),
		// write every message to the out port
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			fmt.Printf("got %s\n", msg)
		}),
	)

	// listen for MIDI
	err = rd.ListenTo(in)
	checkErr(err)

	err = writer.NoteOn(wr, 60, 100)
	checkErr(err)

	time.Sleep(1)
	err = writer.NoteOff(wr, 60)

	checkErr(err)
	// Output: got channel.NoteOn channel 0 key 60 velocity 100
	// got channel.NoteOff channel 0 key 60
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
