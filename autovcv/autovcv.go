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
package autovcv

import (
    "math/rand"
    "github.com/pommicket/autovcv/vcv"
)

type input struct {
    module int
    id int
}

func RandomVCV(nModules int, nWires int) *vcv.VCVFile {
    var v vcv.VCVFile

    // Create a slice of modules
    modules := make([]string, len(vcv.Modules))
    i := 0
    for module, _ := range vcv.Modules {
        modules[i] = module
        i++
    }


    for i := 0; i < nModules; i++ {
        model := modules[rand.Intn(len(modules))]
        info := vcv.Modules[model]
        params := make([]float64, info.Params)
        for p := 0; p < info.Params; p++ {
            b := vcv.ModuleParamBounds[model][p]

            params[p] = (b.Max - b.Min) * rand.Float64() + b.Min
        }
        v.AddModule(model, params, rand.Intn(100), rand.Intn(3))
    }

    if v.NumberOfModules() == 1 {
        // No wires
        return &v
    }

    // Which inputs have been used?
    inputsUsed := make(map[input]bool)

    for i := 0; i < nWires; i++ {
        m1id := 0
        m2id := 0
        for m1id == m2id {
            m1id = rand.Intn(v.NumberOfModules())
            m2id = rand.Intn(v.NumberOfModules())
        }
        m1outputs := vcv.Modules[v.GetModule(m1id).Model].Outputs
        m2inputs := vcv.Modules[v.GetModule(m2id).Model].Inputs

        outId := rand.Intn(m1outputs)
        inId := rand.Intn(m2inputs)
        input := input{m2id, inId}
        _, used := inputsUsed[input]
        if used { continue } // We've already used this input
        inputsUsed[input] = true
        v.AddWire(m1id, outId, m2id, inId, "#00cc00")

    }

    return &v
}
