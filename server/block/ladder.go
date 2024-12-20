package block

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

// Ladder is a wooden block used for climbing walls either vertically or horizontally. They can be placed only on
// the sides of other blocks.
type Ladder struct {
	transparent
	sourceWaterDisplacer

	// Facing is the side of the block the ladder is currently attached to. cube.FaceDown and cube.FaceUp
	// do not do anything in game but they are still valid states.
	Facing cube.Face
}

// NeighbourUpdateTick ...
func (l Ladder) NeighbourUpdateTick(pos, _ cube.Pos, tx *world.Tx) {
	if _, ok := tx.Block(pos.Side(l.Facing.Opposite())).(LightDiffuser); ok {
		tx.SetBlock(pos, nil, nil)
		tx.AddParticle(pos.Vec3Centre(), particle.BlockBreak{Block: l})
		dropItem(tx, item.NewStack(l, 1), pos.Vec3Centre())
	}
}

// UseOnBlock ...
func (l Ladder) UseOnBlock(pos cube.Pos, face cube.Face, _ mgl64.Vec3, tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	pos, face, used := firstReplaceable(tx, pos, face, l)
	if !used {
		return false
	}
	if face == cube.FaceUp || face == cube.FaceDown {
		return false
	}
	if _, ok := tx.Block(pos.Side(face.Opposite())).(LightDiffuser); ok {
		found := false
		for _, i := range []cube.Face{cube.FaceSouth, cube.FaceNorth, cube.FaceEast, cube.FaceWest} {
			if diffuser, ok := tx.Block(pos.Side(i)).(LightDiffuser); !ok || diffuser.LightDiffusionLevel() == 15 {
				found = true
				face = i.Opposite()
				break
			}
		}
		if !found {
			return false
		}
	}
	l.Facing = face

	place(tx, pos, l, user, ctx)
	return placed(ctx)
}

// EntityInside ...
func (l Ladder) EntityInside(_ cube.Pos, _ *world.Tx, e world.Entity) {
	if fallEntity, ok := e.(fallDistanceEntity); ok {
		fallEntity.ResetFallDistance()
	}
}

// SideClosed ...
func (l Ladder) SideClosed(cube.Pos, cube.Pos, *world.Tx) bool {
	return false
}

// BreakInfo ...
func (l Ladder) BreakInfo() BreakInfo {
	return newBreakInfo(0.4, alwaysHarvestable, axeEffective, oneOf(l))
}

// FuelInfo ...
func (Ladder) FuelInfo() item.FuelInfo {
	return newFuelInfo(time.Second * 15)
}

// EncodeItem ...
func (l Ladder) EncodeItem() (name string, meta int16) {
	return "minecraft:ladder", 0
}

// EncodeBlock ...
func (l Ladder) EncodeBlock() (string, map[string]any) {
	return "minecraft:ladder", map[string]any{"facing_direction": int32(l.Facing)}
}

// Model ...
func (l Ladder) Model() world.BlockModel {
	return model.Ladder{Facing: l.Facing}
}

// allLadders ...
func allLadders() (b []world.Block) {
	for _, f := range cube.Faces() {
		b = append(b, Ladder{Facing: f})
	}
	return
}
