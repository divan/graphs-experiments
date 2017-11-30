package main

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

// NgraphBinaryOutput stores graph data as binary files, compatible
// with ngraph library: positions.bin, links.bin, labels.json and meta.json
type NgraphBinaryOutput struct {
}

func (o *NgraphBinaryOutput) Save(l Layout, data *GraphData) error {
	err := o.WritePositionsBin(l)
	if err != nil {
		return err
	}
	log.Println("Positions written into ./positions.bin")

	err = o.WriteLinksBin(data)
	if err != nil {
		return err
	}
	log.Println("Links written into ./links.bin")

	return nil
}

// WritePositionsBin writes position data into 'positions.bin' file in the
// following way: XYZXYZXYZ... where X, Y and Z are coordinates
// for each node in signed 32 bit integer Little Endian format.
func (o *NgraphBinaryOutput) WritePositionsBin(l Layout) error {
	fd, err := os.Create("positions.bin")
	if err != nil {
		return err
	}
	defer fd.Close()

	iw := NewInt32LEWriter(fd)

	nodes := l.Nodes()
	for i, _ := range nodes {
		iw.Write(nodes[i].X)
		iw.Write(nodes[i].Y)
		iw.Write(nodes[i].Z)
		if iw.err != nil {
			return err
		}
	}

	return nil
}

// WriteLinksBin writes links information into `links.bin` file in the
// following way: Sidx,L1idx,L2idx,S2idx,L1idx... where SNidx - is the
// start node index, and LNidx - is the other link end node index.
func (o *NgraphBinaryOutput) WriteLinksBin(data *GraphData) error {
	fd, err := os.Create("links.bin")
	if err != nil {
		return err
	}
	defer fd.Close()

	iw := NewInt32LEWriter(fd)
	for i, node := range data.Nodes {
		if !data.NodeHasLinks(node.ID) {
			continue
		}

		iw.Write(int32(-(i + 1)))
		for j, link := range data.Links {
			if link.Source == node.ID {
				iw.Write(int32(j + 1))
			}
		}
		if iw.err != nil {
			return err
		}
	}
	return nil
}

type Int32LEWriter struct {
	w   io.Writer
	err error
}

func NewInt32LEWriter(w io.Writer) *Int32LEWriter {
	return &Int32LEWriter{
		w: w,
	}
}

func (iw *Int32LEWriter) Write(number int32) {
	if iw.err != nil {
		return
	}

	err := binary.Write(iw.w, binary.LittleEndian, number)
	if err != nil {
		iw.err = err
	}

	return
}