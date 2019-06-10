/*
Copyright (C) 2019 Leo Tenenbaum

This file is part of AutoVCV.

AutoVCV is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

AutoVCV is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with AutoVCV.  If not, see <https://www.gnu.org/licenses/>.
*/
package vcv

import (
    "io"
    "fmt"
)

type ModuleInfo struct {
    Inputs int
    Outputs int
    Params int
}

type Module struct {
    Model string
    Params []float64
    X int // x position
    Y int
}

type Wire struct {
    OutModule int
    OutId int
    InModule int
    InId int
    Color string // in hex notation
}

type VCVFile struct {
    Modules []Module
    Wires []Wire
}

type ParamBounds struct {
    Min float64
    Max float64
}


var Modules = map[string]ModuleInfo{
    "VCO": {4, 4, 7},
    "LFO": {4, 4, 7},
    "VCA-1": {2, 1, 2},
    "Delay": {5, 1, 4},
    "AudioInterface": {2, 2 /* Only use 2 channels */, 0},
}

var ModuleParamBounds = map[string][]ParamBounds {
    "VCO": {{0, 1}, {0, 1}, {-54, 54}, {-1, 1}, {0, 1}, {0, 1}, {0, 1}},
    "LFO": {{0, 1}, {0, 1}, { -8,  8}, {0, 1},  {0, 1}, {0, 1}, {0, 1}},
    "VCA-1": {{0, 1}, {0, 1}},
    "Delay": {{0, 1}, {0, 1}, {0, 1}, {0, 1}},
}

var plugins = map[string]string {
    "VCO": "Fundamental",
    "LFO": "Fundamental",
    "VCA-1": "Fundamental",
    "Delay": "Fundamental",
}

var Versions = map[string]string {
    "Core": "0.6.2c",
    "Fundamental": "0.6.2",
}

func (vcv *VCVFile) NumberOfModules() int {
    return len(vcv.Modules)
}

func (vcv *VCVFile) NumberOfWires() int {
    return len(vcv.Wires)
}

// Adds a module with the given settings to the file
// Returns an ID used to connect modules
func (vcv *VCVFile) AddModule(model string, params []float64, x int, y int) int {
    m := Module {
        model, params, x, y,
    }
    vcv.Modules = append(vcv.Modules, m)
    return len(vcv.Modules)-1
}

// Adds a wire from the module with ID to ID id2
// Returns the wire ID
func (vcv *VCVFile) AddWire(id1 int, whichOutput int, id2 int, whichInput int, color string) int {
    w := Wire {
        id1, whichOutput, id2, whichInput, color,
    }
    vcv.Wires = append(vcv.Wires, w)
    return len(vcv.Wires)-1
}

// Returns the module associated with the given ID
func (vcv *VCVFile) GetModule(id int) *Module {
    return &vcv.Modules[id]
}

func (vcv *VCVFile) GetWire(id int) *Wire {
    return &vcv.Wires[id]
}

func (m *Module) write(out io.Writer, deviceName string) error {
    if m.Model == "AudioInterface" {
        _, err := fmt.Fprintf(out, `{
  "plugin": "Core",
  "version": "%v",
  "model": "AudioInterface",
  "params": [],
  "data": {
    "audio": {
      "driver": 1,
      "deviceName": "%v",
      "offset": 0,
      "maxChannels": 8,
      "sampleRate": 44100,
      "blockSize": 256
    }
  },
  "pos": [
    %v,
    %v
  ]
}`, Versions["Core"], deviceName, m.X, m.Y)

        return err
    }
    _, err := fmt.Fprintf(out, `{
        "plugin": "%v",
        "version": "%v",
        "model": "%v",
        "pos": [%v, %v],
        "params": [`, plugins[m.Model], Versions[plugins[m.Model]], m.Model,
                      m.X, m.Y)
    if err != nil { return err }
    for i, param := range m.Params {
        var maybeComma byte
        if i == 0 {
            maybeComma = ' '
        } else {
            maybeComma = ','
        }
        _, err := fmt.Fprintf(out,  `%c{
    "paramId": %v,
    "value": %v
}`, maybeComma, i, param)
        if err != nil { return err }
    }
    _, err = io.WriteString(out, "]}")
    if err != nil { return err }

    return nil
}

func (w *Wire) write(out io.Writer) error {
    _, err := fmt.Fprintf(out, `{
    "color": "%v",
    "outputModuleId": %v,
    "outputId": %v,
    "inputModuleId": %v,
    "inputId": %v
}`, w.Color, w.OutModule, w.OutId, w.InModule, w.InId)
    return err
}

func (vcv *VCVFile) Write(out io.Writer, deviceName string) error {
    _, err := fmt.Fprintf(out, `{
        "version": "%v",
        "modules": [`, Versions["Core"])
    if err != nil { return err }
    for i, module := range vcv.Modules {
        if i != 0 {
            _, err := io.WriteString(out, ",")
            if err != nil { return err }
        }
        err := module.write(out, deviceName)
        if err != nil { return err }
    }
    _, err = io.WriteString(out, `],
    "wires": [`)
    if (err != nil) { return err }
    for i, wire := range vcv.Wires {
        if i != 0 {
            _, err := io.WriteString(out, ",")
            if err != nil { return err }
        }
        err := wire.write(out)
        if err != nil { return err }

    }
    _, err = io.WriteString(out, "]}")
    return err
}
