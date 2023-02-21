/*
 * This file is part of the Atomic Stack (https://github.com/libatomic/atomic).
 * Copyright (c) 2020 Atomic Publishing.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, version 3.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	convertCmd = &cli.Command{
		Name:   "convert",
		Usage:  "convert a file to a csv",
		Action: convert,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file",
				Usage: "file we want to convert to a csv",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "output file",
			},
		},
	}
)

type (
	measurement struct {
		Time time.Time                                `json:"time,omitempty"`
		Data map[string]map[string]map[string]float64 `json:"data,omitempty"`
	}

	telemetry struct {
		Averages map[time.Time]measurement `json:"telemetry_averages"`
		Devices  []string                  `json:"device_ids,omitempty"`
	}
)

func convert(c *cli.Context) error {
	configFile, err := os.Open(c.String("file"))

	if err != nil {
		return err
	}
	defer configFile.Close()

	l := telemetry{}
	decodeErr := json.NewDecoder(configFile).Decode(&l)

	if decodeErr != nil {
		return decodeErr
	}

	f, err := os.Create(c.String("output"))
	if err != nil {
		return err
	}
	defer f.Close()

	b := make([]time.Time, 0)
	for _, j := range l.Averages {
		b = append(b, j.Time)
	}

	sort.Slice(b, func(i, j int) bool {
		return b[i].Before(b[j])
	})

	alreadyExists := make(map[string]int)
	columns := make([][]string, 0)
	for _, avg := range l.Averages {
		for name, d_data := range avg.Data {
			for depth, _ := range d_data {
				_, ok := alreadyExists[name+"_"+depth]
				if !ok {
					columns = append(columns, []string{name, depth})
					alreadyExists[name+"_"+depth] = 0
				}
			}
		}
	}

	sort.Slice(columns, func(i, j int) bool {
		return columns[i][0]+"_"+columns[i][1] < columns[j][0]+"_"+columns[j][1]
	})

	f.WriteString("time,")
	for _, v := range columns {
		f.WriteString(v[0] + "_" + v[1] + ",")
	}
	f.WriteString("\n")

	for _, avg_t := range b {
		f.WriteString(avg_t.String() + ",")
		for _, val := range columns {
			f.WriteString(fmt.Sprintf("%f,", l.Averages[avg_t].Data[val[0]][val[1]]["value"]))
		}
		f.WriteString(("\n"))
	}

	f.Sync()

	fmt.Println("SUCCESS: sensor saved")
	return nil
}
