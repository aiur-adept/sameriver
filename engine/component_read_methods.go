
//
//
// THIS FILE HAS BEEN GENERATED BY sameriver-generate
//
//
// DO NOT MODIFY BY HAND UNLESS YOU WANNA HAVE A GOOD TIME WHEN THE NEXT
// GENERATION DESTROYS WHAT YOU WROTE. UNLESS YOU KNOW HOW TO HAVE A GOOD TIME
//
//

package engine

import (
	"errors"
	"fmt"
)

func (ct *ComponentsTable) ReadBox(entity EntityToken) (Box, error) {
	ct.em.rLockEntityComponent(e, BOX_COMPONENT)
	defer ct.em.rUnlockEntityComponent(e, BOX_COMPONENT)
	if !ct.em.entityTable.genValidate(e) {
		return Box{}, errors.New.New(fmt.Sprintf("%+v no longer exists", e))
	}
	return ct.Box[e.ID], nil
}

func (ct *ComponentsTable) ReadSprite(entity EntityToken) (Sprite, error) {
	ct.em.rLockEntityComponent(e, SPRITE_COMPONENT)
	defer ct.em.rUnlockEntityComponent(e, SPRITE_COMPONENT)
	if !ct.em.entityTable.genValidate(e) {
		return Sprite{}, errors.New.New(fmt.Sprintf("%+v no longer exists", e))
	}
	return ct.Sprite[e.ID], nil
}

func (ct *ComponentsTable) ReadTagList(entity EntityToken) (TagList, error) {
	ct.em.rLockEntityComponent(e, TAGLIST_COMPONENT)
	defer ct.em.rUnlockEntityComponent(e, TAGLIST_COMPONENT)
	if !ct.em.entityTable.genValidate(e) {
		return TagList{}, errors.New.New(fmt.Sprintf("%+v no longer exists", e))
	}
	return DeepCopyTagList(ct.TagList[e.ID]), nil
}

func (ct *ComponentsTable) ReadVelocity(entity EntityToken) ([2]float32, error) {
	ct.em.rLockEntityComponent(e, VELOCITY_COMPONENT)
	defer ct.em.rUnlockEntityComponent(e, VELOCITY_COMPONENT)
	if !ct.em.entityTable.genValidate(e) {
		return [2]float32{}, errors.New.New(fmt.Sprintf("%+v no longer exists", e))
	}
	return ct.Velocity[e.ID], nil
}
