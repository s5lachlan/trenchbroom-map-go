# Trenchbroom .MAP Parser for Go

A renderer agnostic helper library for Go that parses and generates meshes for .MAP files. 
> [!WARNING]
> This library is experimental and primarily designed for my own personal use which involves Raylib, but it's intended to be agnostic.

## What this library does

This library can parse Valve 220 .MAP files into structs, and generate geometry.

It does not do the following

- Parse other .MAP types like Standard.
- Include an ECS for your world or have any other prebuilt management system for func entities etc.

## Usage (Raylib)

A raylib helper package is provided at `trenchbroom-map-go/raylib`, here is an example of how to use it.
```go
package main

import (
	"fmt"
	tb "trenchbroom-map-go"
	tbrl "trenchbroom-map-go/raylib"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var screenWidth int32 = 800
var screenHeight int32 = 800
var camera = rl.Camera{
	Position:   rl.NewVector3(0.0, 2.0, 4.0),
	Target:     rl.NewVector3(0.0, 2.0, 0.0),
	Up:         rl.NewVector3(0.0, 1.0, 0.0),
	Fovy:       90.0,
	Projection: rl.CameraPerspective,
}

const WORLD_SCALE = 0.1

func main() {
	// Initialize raylib
	rl.InitWindow(screenWidth, screenHeight, "My game")
	rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyEscape)
	// Close window at the end of the program
	defer rl.CloseWindow()
	// Parse your desired map
	tbMap, err := tb.Parse("example/maps/example.map")
	// Scale it down
	tb.ScaleMap(&tbMap, WORLD_SCALE)
	// Handle any errors however you'd like
	if err != nil {
		panic(err)
	}
	textures := tbrl.LoadTexturesFromMap(&tbMap, "example/textures/", ".png")
	// Generate geometry for worldspawn only
	models := tbrl.GenModelsFromEntity(tbMap.Entities[0], textures, tbMap)
	// Unload models and textures at the end of the program
	defer tbrl.UnloadModels(models)
	defer tbrl.UnloadTextureMap(textures)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.DarkGray)
		rl.BeginMode3D(camera)
		
		// Draw every model in entity 0
		for _, m := range models {
			rl.DrawModel(m, rl.Vector3Zero(), 0.1, rl.White)
		}

		rl.EndMode3D()
		rl.EndDrawing()
	}
}
```
If you would like to render multiple entities you'll have to build your own system for it in which you can decide to render an entity based on it's data and parse it into another system that allows you to move models, etc.

## Usage (any other renderer library)

Load your renderer, then when it's loaded, parse the map
```go
// Parse your desired map
tbMap, err := tb.Parse("example/maps/mymap.map")
// Scale it down
tb.ScaleMap(&tbMap, WORLD_SCALE)
// Handle any errors
if err != nil {
	panic(err)
}
```
`tbMap` contains map information, entity data and brush information with no meshes. To generate geometry do the following
```go
// Update the map's texture dimensions so that the geometry generator can have correct UVs 
tb.UpdateTextureDimensions(&tbMap, "example/textures/", ".png")
// Build geometry for one entity
meshes := tb.BuildEntityMeshes(tbMap.Entities[0], tbMap)
```
From here you now must create a model to render in your renderer from the meshes. If you're using raylib it's recommended to make a model per entity mesh, they're easier to texture this way. If you have a brush where one side is one texture and the other sides are other textures, then depending on the texture, they should be their own model.

## Resources
Resources:
- Quake Wiki: https://quakewiki.org/wiki/Quake_Map_Format
- Stefanha's 2001 Paper: https://github.com/stefanha/map-files
- Extra resources: https://book.leveldesignbook.com/appendix/resources/formats/map

Other .MAP Parsing and Mesh Generating Implementations in other languages/engines:

- https://github.com/andyp123/blender_io_mesh_qmap
- https://github.com/skyelynwaddell/trenchbroom-raylib-c
- https://github.com/QodotPlugin/Qodot
- https://github.com/progsy/phxmap
