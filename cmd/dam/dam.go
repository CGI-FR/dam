// Copyright (C) 2023 CGI France
//
// This file is part of dam.
//
// dam is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// dam is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with dam.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"fmt"
	"io"
)

const linebreak byte = 10

type Dam struct {
	input  *bufio.Scanner
	buffer chan []byte
	output *bufio.Writer
}

func NewDam(input io.Reader, output io.Writer, capacity int) *Dam {
	return &Dam{
		input:  bufio.NewScanner(input),
		buffer: make(chan []byte, capacity),
		output: bufio.NewWriter(output),
	}
}

func (d *Dam) GetSize() int {
	return len(d.buffer)
}

func (d *Dam) IsEmpty() bool {
	return d.GetSize() == 0
}

func (d *Dam) Drop() error {
	if _, err := d.output.Write(append(<-d.buffer, linebreak)); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (d *Dam) Drain() error {
	for !d.IsEmpty() {
		if err := d.Drop(); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	d.output.Flush()

	return nil
}

func (d *Dam) Spill() error {
	for d.GetSize() >= cap(d.buffer) {
		if err := d.Drop(); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	d.output.Flush()

	return nil
}

func (d *Dam) Add(data []byte) error {
	if err := d.Spill(); err != nil {
		return fmt.Errorf("%w", err)
	}

	d.buffer <- append(make([]byte, 0, len(data)), data...)

	return nil
}

func (d *Dam) Run() error {
	for d.input.Scan() {
		if err := d.Add(d.input.Bytes()); err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	if err := d.input.Err(); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := d.Drain(); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
