package tb

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

const (
	EntityMode  = "ENTITY"
	BrushMode   = "BRUSH"
	OutsideMode = "OUTSIDE"

	OPENING_TAG = "{"
	CLOSING_TAG = "}"
)

type Entity struct {
	Data    map[string]string
	Brushes []Brush
}

type Map struct {
	// A slice of every entity in this map.
	Entities []Entity
	// A list of every texture this map uses.
	// Use this to load textures with your renderer if you want to dynamically load them for each level.
	Textures []string
	// A map of strings to a Vector2 where X is texture width and Y is texture height.
	// This is to be used in conjunction with Map.Textures and is necessary to define for generating
	// level geometry. This map will be used to calculate Texture UVs correctly.
	//
	// By running UpdateTextureDimensions() the library can fill this for you.
	TextureDimensions map[string]Vector2
	// The path of the loaded map.
	Path string
	// The name of the game attached to the map.
	GameName string
	// The name of the format attached to the map.
	FormatName string
	// The brush count of the entire map
	BrushCount int
	// The entity count of the entire map
	EntityCount int
	// Epsilon for generating meshes. 0.5 by default.
	// Epsilon is scaled up or down with ScaleMap, to help generate more accurate geometry.
	// You can change this yourself, just before generating geometry, for differing results.
	Epsilon float
}

type Brush struct {
	FaceCount  int
	BrushFaces []BrushFace
}

type BrushFace struct {
	Pos1    Vector3
	Pos2    Vector3
	Pos3    Vector3
	Texture string

	UV_S        Vector4
	UV_T        Vector4
	UV_Rotation float
	U_Scale     float
	V_Scale     float
}

var PrintDebug = false

func Parse(FileName string) (Map, error) {
	newMap := Map{Path: FileName}
	newMap.Epsilon = 0.5
	currentBrush := Brush{}
	currentEntity := Entity{}
	currentEntity.Data = make(map[string]string)

	entityCount := -1
	brushCount := -1

	file, err := os.Open(FileName)
	if err != nil {
		fmt.Println("[MAPMANAGER] Error occurred in map loading: ", err)
		return Map{}, err
	}

	defer file.Close()

	mode := OutsideMode

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Trim all spaces.
		currentLine := strings.TrimSpace(scanner.Text())

		// Set the game name
		if gameName, ok := strings.CutPrefix(currentLine, "// Game:"); ok {
			gameName = strings.TrimSpace(gameName)
			newMap.GameName = gameName
		}
		// Set the format name
		if formatName, ok := strings.CutPrefix(currentLine, "// Format:"); ok {
			formatName = strings.TrimSpace(formatName)
			newMap.FormatName = formatName
			switch formatName {
			case "Valve":
				continue
			default:
				return Map{}, errors.New("Map type: " + formatName + " is not supported. Only Valve type is supported")

			}
		}

		if strings.HasPrefix(currentLine, "//") ||
			string(currentLine[0]) == "\n" ||
			string(currentLine[0]) == "" {
			continue
		}

		if string(currentLine[0]) == OPENING_TAG {
			switch mode {
			case OutsideMode:
				mode = EntityMode
				currentEntity = Entity{}
				currentEntity.Data = make(map[string]string)
				entityCount += 1
			case EntityMode:
				mode = BrushMode
				currentBrush = Brush{}
				brushCount += 1
			}
			continue
		}

		if string(currentLine[0]) == CLOSING_TAG {
			switch mode {
			case BrushMode:
				mode = EntityMode
				// Append brush to entity and update map brush count
				currentEntity.Brushes = append(currentEntity.Brushes, currentBrush)
				newMap.BrushCount += 1
			case EntityMode:
				mode = OutsideMode
				newMap.Entities = append(newMap.Entities, currentEntity)
				newMap.EntityCount += 1
			}
			continue
		}

		if mode == BrushMode {
			brushface := BrushFace{}

			count, err := fmt.Sscanf(string(currentLine),
				"( %f %f %f ) ( %f %f %f ) ( %f %f %f ) %s [ %f %f %f %f ] [ %f %f %f %f ] %f %f %f",
				&brushface.Pos1.X, &brushface.Pos1.Y, &brushface.Pos1.Z,
				&brushface.Pos2.X, &brushface.Pos2.Y, &brushface.Pos2.Z,
				&brushface.Pos3.X, &brushface.Pos3.Y, &brushface.Pos3.Z,
				&brushface.Texture,
				&brushface.UV_S.X, &brushface.UV_S.Y, &brushface.UV_S.Z, &brushface.UV_S.W,
				&brushface.UV_T.X, &brushface.UV_T.Y, &brushface.UV_T.Z, &brushface.UV_T.W,
				&brushface.UV_Rotation, &brushface.U_Scale, &brushface.V_Scale,
			)

			if err == nil && count == 21 {
				currentBrush.FaceCount = len(currentBrush.BrushFaces)
				if !slices.Contains(newMap.Textures, brushface.Texture) {
					newMap.Textures = append(newMap.Textures, brushface.Texture)
				}
				currentBrush.BrushFaces = append(currentBrush.BrushFaces, brushface)
			} else {
				fmt.Println("[MAP GENERATION] Error has occured!", err, "Count", count, "in ", brushface)
				return Map{}, err
			}
		}

		if mode == EntityMode {
			var EntityArg1, EntityArg2 string
			count, err := fmt.Sscanf(string(currentLine), "%q %q", &EntityArg1, &EntityArg2)
			if err == nil && count == 2 {
				if PrintDebug {
					fmt.Println("Entity Properties:", EntityArg1, EntityArg2)
				}
				currentEntity.Data[EntityArg1] = EntityArg2
				continue
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("[MAPMANAGER] Line failed to read. ", err)
	}

	return newMap, nil
}

// Scale an entire map up or down.
// Recommended: ScaleMap(map, 0.1)
func ScaleMap(tbMap *Map, scale float) {
	tbMap.Epsilon *= scale
	for e, entity := range tbMap.Entities {
		for b, brush := range entity.Brushes {
			for bf, face := range brush.BrushFaces {
				tbMap.Entities[e].Brushes[b].BrushFaces[bf].Pos1 = Vector3Scale(face.Pos1, scale)
				tbMap.Entities[e].Brushes[b].BrushFaces[bf].Pos2 = Vector3Scale(face.Pos2, scale)
				tbMap.Entities[e].Brushes[b].BrushFaces[bf].Pos3 = Vector3Scale(face.Pos3, scale)
				tbMap.Entities[e].Brushes[b].BrushFaces[bf].U_Scale *= scale
				tbMap.Entities[e].Brushes[b].BrushFaces[bf].V_Scale *= scale
			}
		}
	}
}
