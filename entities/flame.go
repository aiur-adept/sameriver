package entities

import (
    "fmt"
    "math"
    "math/rand"

    "github.com/veandco/go-sdl2/sdl"

    "github.com/dt-rush/donkeys-qquest/engine"
    "github.com/dt-rush/donkeys-qquest/engine/components"
)

func SpawnFlame (entity_manager *engine.EntityManager,
    active_component *components.ActiveComponent,
    position_component *components.PositionComponent,
    velocity_component *components.VelocityComponent,
    hitbox_component *components.HitboxComponent,
    sprite_component *components.SpriteComponent,
    logic_component *components.LogicComponent,
    initial_position [2]float64) int {

    flame_components := []engine.Component{
        engine.Component (active_component),
        engine.Component (position_component),
        engine.Component (velocity_component),
        engine.Component (hitbox_component),
        engine.Component (sprite_component),
        engine.Component (logic_component),
    }

    flame_id := entity_manager.SpawnEntity (flame_components)

    // set component values

    flame_active := true
    active_component.Set (flame_id, flame_active)

    //flame_position := [2]float64 {
    //    rand.Float64() * float64 (constants.WINDOW_WIDTH - 20) + 20,
    //    rand.Float64() * float64 (constants.WINDOW_HEIGHT - 20) + 20,
    //}
    flame_position := initial_position
    position_component.Set (flame_id, flame_position)

    flame_hitbox := [2]float64{30, 30}
    hitbox_component.Set (flame_id, flame_hitbox)

    flame_sprite := sprite_component.IndexOf ("flame.png")
    sprite_component.Set (flame_id, flame_sprite)

    // add flame logic

    player_id := entity_manager.GetTagEntityUnique ("player")
    flame_logic_name := fmt.Sprintf ("flame logic %d", flame_id)
    flame_logic_func := FlameLogic (flame_id,
                                    player_id,
                                    position_component,
                                    velocity_component,
                                    sprite_component)
    flame_logic_unit := components.LogicUnit{flame_logic_name,
                                                flame_logic_func}
    logic_component.Set (flame_id, flame_logic_unit)

    // add tag

    entity_manager.TagEntity (flame_id, "flame")

    return flame_id
}

func FlameLogic (flame_id int,
    player_id int,
    position_component *components.PositionComponent,
    velocity_component *components.VelocityComponent,
    sprite_component *components.SpriteComponent) (func (int)) {

    // closure state

    // time accumulator
    accum_200 := engine.CreateTimeAccumulator(200)
    accum_500 := engine.CreateTimeAccumulator(500)
    // current heading
    var heading []float64 = make ([]float64, 2)

    // enclosed function

    return func (dt_ms int) {

        player_pos := position_component.Get (player_id)
        flame_pos := position_component.Get (flame_id)
        flame_vel := velocity_component.Get (flame_id)

        if accum_200.Tick (dt_ms) {
            // find out how to get to the player
            vector_to_player := [2]float64{
                player_pos[0] - flame_pos[0],
                player_pos[1] - flame_pos[1],
            }
            // normalize the above vector
            scale_factor := math.Sqrt (
                (vector_to_player[0] * vector_to_player[0] + vector_to_player[1] * vector_to_player[1]))
            vector_to_player[0] /= scale_factor
            vector_to_player[1] /= scale_factor
            // ... and trigger change of heading
            heading [0] = 75 * (2 * (rand.Float64()*2 - 1) + vector_to_player [0])
            heading [1] = 75 * (2 * (rand.Float64()*2 - 1) + vector_to_player [1])
            flame_vel [0] = heading [0]
            flame_vel [1] = heading [1]

        }

        if accum_500.Tick (dt_ms) {
            // toggle flip horizontal
            current_flip := sprite_component.GetFlip (flame_id)
            var new_flip sdl.RendererFlip
            if current_flip == sdl.FLIP_NONE {
                new_flip = sdl.FLIP_HORIZONTAL
            } else {
                new_flip = sdl.FLIP_NONE
            }
            sprite_component.SetFlip (flame_id, new_flip)
        }

        // TODO replace statements like this with defer statements or at least
        // find a way to directly update map values, to avoid this weird
        // modify and replace pattern
        velocity_component.Set (flame_id, flame_vel)
    }
}

