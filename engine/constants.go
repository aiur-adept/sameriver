/**
  *
  * This file defines constants for the game engine operation,
  * mostly flags to be set at compile-time.
  *
  *
**/

package engine

import (
	"time"
)

const VERSION = "0.3.2"
const AUDIO_ON = true

const DEBUG_ENTITY_MANAGER = true
const DEBUG_EVENTS = true
const DEBUG_UPDATED_ENTITY_LISTS = true
const DEBUG_GOROUTINES = true
const DEBUG_ENTITY_LOGIC = true
const DEBUG_ENTITY_MANAGER_UPDATE_TIMING = false
const DEBUG_DESPAWN = true
const DEBUG_ATOMIC_MODIFY = true
const DEBUG_ENTITY_CLASS = true
const DEBUG_WORLD_LOGIC = true
const DEBUG_ENTITY_LOCKS = true

const FPS = 60
const FRAME_SLEEP = (1000 / FPS) * time.Millisecond
const MAX_ENTITIES = 1024

const COLLISION_RATELIMIT_TIMEOUT_MS = 500

const SPAWN_CHANNEL_CAPACITY = 128
const ENTITY_MODIFICATION_CHANNEL_CAPACITY = MAX_ENTITIES
const ENTITY_QUERY_WATCHER_CHANNEL_CAPACITY = 16
const EVENT_CHANNEL_CAPACITY = 16
