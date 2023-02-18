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
	"time"

	"github.com/urfave/cli/v2"
)
  
 var (
	convertCmd = &cli.Command{
	Name:    "convert",
	Usage:   "convert a file to a csv",
	Action:    convert,
	Flags: []cli.Flag{
		&cli.StringFlag{ 
			Name:    "file",
			Usage:   "file we want to convert to a csv",
	},
		&cli.StringFlag{
			Name:    "output",
			Usage:   "output file",
	},
	},
	
	}
 )

 type (
	measurement struct {
		Time time.Time `json:"time,omitempty"`
		Data map[string]map[string]map[string]float64 `json:"data,omitempty"`
	} 

	telemetry struct {
		Averages map[time.Time]measurement `json:"telemetry_averages"`
		Devices []string `json:"device_ids,omitempty"`
	}
 )

 func convert(c *cli.Context) error {
	configFile, err := os.Open(c.String("file"))
	
    if err != nil {
        return err
    }
	defer configFile.Close()

	//l := make(map[string]map[string]map[string]map[string]map[string]float64)
	l := telemetry{}
	decodeErr := json.NewDecoder(configFile).Decode(&l)

	if decodeErr != nil {
		return decodeErr
	}
	
	fmt.Println("SUCCESS: sensor saved")
	return nil
 }
 