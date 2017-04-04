package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/3d0c/gmf"
)

// Simple transcoder, it guesses source format and codecs and tries to convert it to a:mp3.
func transcode(srcFileName string, dstFileName string, bitrate int) error {
	inputCtx, err := gmf.NewInputCtx(srcFileName)
	if err != nil {
		return err
	}
	defer inputCtx.CloseInputAndRelease()

	srcAst, err := inputCtx.GetBestStream(gmf.AVMEDIA_TYPE_AUDIO)
	cc := srcAst.CodecCtx()

	fmt.Println("input channels: ", cc.Channels())
	fmt.Println("input channelLayout: ", cc.ChannelLayout())
	fmt.Println("input frameSize:", cc.FrameSize())
	fmt.Println("input sampleFmt:", cc.SampleFmt())
	fmt.Printf("input bitRate: %vk\n", cc.BitRate()/1000)

	codec, err := gmf.FindEncoder("libmp3lame")
	if err != nil {
		return err
	}

	occ := gmf.NewCodecCtx(codec)
	if occ == nil {
		return err
	}
	defer gmf.Release(occ)

	occ.SetSampleFmt(gmf.AV_SAMPLE_FMT_S16P).
		SetSampleRate(cc.SampleRate()).
		SetChannels(cc.Channels()).
		SetBitRate(bitrate)
	channelLayout := occ.SelectChannelLayout()
	occ.SetChannelLayout(channelLayout)

	fmt.Println("output channels: ", occ.Channels())
	fmt.Println("output channelLayout: ", occ.ChannelLayout())
	fmt.Println("output frameSize:", occ.FrameSize())
	fmt.Println("output sampleFmt:", occ.SampleFmt())
	fmt.Printf("output bitRate: %vk\n", occ.BitRate()/1000)

	if err := occ.Open(nil); err != nil {
		return err
	}

	/// resample
	options := []*gmf.Option{
		{"in_channel_count", cc.Channels()},
		{"out_channel_count", occ.Channels()},
		{"in_sample_rate", cc.SampleRate()},
		{"out_sample_rate", occ.SampleRate()},
		{"in_sample_fmt", gmf.SampleFmt(cc.SampleFmt())},
		{"out_sample_fmt", gmf.SampleFmt(gmf.AV_SAMPLE_FMT_S16P)},
	}

	swrCtx := gmf.NewSwrCtx(options, occ)
	if swrCtx == nil {
		return errors.New("unable to create Swr Context")
	}

	outputCtx, err := gmf.NewOutputCtx(dstFileName)
	if err != nil {
		return err
	}
	defer outputCtx.CloseOutputAndRelease()

	ost := outputCtx.NewStream(codec)
	if ost == nil {
		return errors.New(fmt.Sprintf("Unable to create stream for [%s]\n", codec.LongName()))
	}
	defer gmf.Release(ost)

	ost.SetCodecCtx(occ)

	if err := outputCtx.WriteHeader(); err != nil {
		return err
	}

	count := 0
	for packet := range inputCtx.GetNewPackets() {
		ist, err := inputCtx.GetStream(packet.StreamIndex())
		if err != nil {
			return err
		}
		if !ist.IsAudio() {
			continue
		}
		srcFrame, got, ret, err := packet.DecodeToNewFrame(ist.CodecCtx())
		gmf.Release(packet)
		if !got || ret < 0 || err != nil {
			log.Println("capture audio error:", err)
			continue
		}

		dstFrame := swrCtx.Convert(srcFrame)

		if dstFrame == nil {
			continue
		}
		writePacket, ready, _ := dstFrame.EncodeNewPacket(occ)
		for ready {
			if err := outputCtx.WritePacket(writePacket); err != nil {
				log.Println("write packet err", err.Error())
			}

			gmf.Release(writePacket)

			if count < int(cc.SampleRate())*10 {
				break
			} else { //exit
				writePacket, ready, _ = dstFrame.FlushNewPacket(occ)
			}
		}
		gmf.Release(dstFrame)
		gmf.Release(srcFrame)
	}
	return nil
}
