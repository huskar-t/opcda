package opcda

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/huskar-t/opcda/com"
)

type OPCGroups struct {
	iServer                *com.IOPCServer
	iCommon                *com.IOPCCommon
	parent                 *OPCServer
	groupID                uint32
	defaultActive          bool
	defaultGroupUpdateRate uint32
	defaultDeadband        float32
	defaultLocaleID        uint32
	defaultGroupTimeBias   int32
	groups                 []*OPCGroup
	sync.RWMutex
}

func NewOPCGroups(opcServer *OPCServer) *OPCGroups {
	return &OPCGroups{
		parent:                 opcServer,
		iServer:                opcServer.iServer,
		defaultActive:          true,
		defaultGroupUpdateRate: uint32(1000),
		defaultDeadband:        float32(0.0),
		defaultLocaleID:        uint32(0x0400),
		defaultGroupTimeBias:   int32(0),
		iCommon:                opcServer.iCommon,
	}
}

// GetParent Returns reference to the parent OPCServer object.
func (gs *OPCGroups) GetParent() *OPCServer {
	return gs.parent
}

// GetDefaultGroupIsActive get the default active state for OPCGroups created using Groups.Add
func (gs *OPCGroups) GetDefaultGroupIsActive() bool {
	return gs.defaultActive
}

// SetDefaultGroupIsActive set the default active state for OPCGroups created using Groups.Add
func (gs *OPCGroups) SetDefaultGroupIsActive(defaultActive bool) {
	gs.defaultActive = defaultActive
}

// GetDefaultGroupUpdateRate get the default update rate (in milliseconds) for OPCGroups created using Groups.Add
func (gs *OPCGroups) GetDefaultGroupUpdateRate() uint32 {
	return gs.defaultGroupUpdateRate
}

// SetDefaultGroupUpdateRate set the default update rate (in milliseconds) for OPCGroups created using Groups.Add
func (gs *OPCGroups) SetDefaultGroupUpdateRate(defaultGroupUpdateRate uint32) {
	gs.defaultGroupUpdateRate = defaultGroupUpdateRate
}

// GetDefaultGroupDeadband get the default deadband for OPCGroups created using Groups.Add
func (gs *OPCGroups) GetDefaultGroupDeadband() float32 {
	return gs.defaultDeadband
}

// SetDefaultGroupDeadband set the default deadband for OPCGroups created using Groups.Add
func (gs *OPCGroups) SetDefaultGroupDeadband(defaultDeadband float32) {
	gs.defaultDeadband = defaultDeadband
}

// GetDefaultGroupLocaleID get the default locale for OPCGroups created using Groups.Add.
func (gs *OPCGroups) GetDefaultGroupLocaleID() uint32 {
	return gs.defaultLocaleID
}

// SetDefaultGroupLocaleID set the default locale for OPCGroups created using Groups.Add.
func (gs *OPCGroups) SetDefaultGroupLocaleID(defaultLocaleID uint32) {
	gs.defaultLocaleID = defaultLocaleID
}

// GetDefaultGroupTimeBias get the default time bias for OPCGroups created using Groups.Add.
func (gs *OPCGroups) GetDefaultGroupTimeBias() int32 {
	return gs.defaultGroupTimeBias
}

// SetDefaultGroupTimeBias set the default time bias for OPCGroups created using Groups.Add.
func (gs *OPCGroups) SetDefaultGroupTimeBias(defaultGroupTimeBias int32) {
	gs.defaultGroupTimeBias = defaultGroupTimeBias
}

// GetCount Required property for collections.
func (gs *OPCGroups) GetCount() int {
	gs.RLock()
	defer gs.RUnlock()
	return len(gs.groups)
}

// Item Returns an OPCGroup by ItemSpecifier. ItemSpecifier is the name or 0-based index into the collection
func (gs *OPCGroups) Item(index int32) (*OPCGroup, error) {
	gs.RLock()
	defer gs.RUnlock()
	if index < 0 || index >= int32(len(gs.groups)) {
		return nil, errors.New("index out of range")
	}
	return gs.groups[index], nil
}

// ItemByName Returns an OPCGroup by name
func (gs *OPCGroups) ItemByName(name string) (*OPCGroup, error) {
	gs.RLock()
	defer gs.RUnlock()
	for _, v := range gs.groups {
		if v.groupName == name {
			return v, nil
		}
	}
	return nil, errors.New("not found")
}

// Add Creates a new OPCGroup object and adds it to the collections
func (gs *OPCGroups) Add(szName string) (*OPCGroup, error) {
	gs.Lock()
	defer gs.Unlock()
	hClientGroup := atomic.AddUint32(&gs.groupID, 1)
	phServerGroup, pRevisedUpdateRate, ppUnk, err := gs.iServer.AddGroup(
		szName,
		gs.defaultActive,
		gs.defaultGroupUpdateRate,
		hClientGroup,
		&gs.defaultGroupTimeBias,
		&gs.defaultDeadband,
		gs.defaultLocaleID,
		&com.IID_IOPCGroupStateMgt,
	)
	if err != nil {
		return nil, err
	}
	opcGroup, err := NewOPCGroup(gs, ppUnk, hClientGroup, phServerGroup, szName, pRevisedUpdateRate)
	if err != nil {
		ppUnk.Release()
		return nil, err
	}
	gs.groups = append(gs.groups, opcGroup)
	return opcGroup, nil
}

// GetOPCGroupByName Returns an OPCGroup by name
func (gs *OPCGroups) GetOPCGroupByName(name string) (*OPCGroup, error) {
	return gs.ItemByName(name)
}

// GetOPCGroup Returns an OPCGroup by server handle
func (gs *OPCGroups) GetOPCGroup(serverHandle uint32) (*OPCGroup, error) {
	gs.RLock()
	defer gs.RUnlock()
	for _, v := range gs.groups {
		if v.serverGroupHandle == serverHandle {
			return v, nil
		}
	}
	return nil, errors.New("not found")
}

// Remove Removes an OPCGroup from the collection
func (gs *OPCGroups) Remove(serverHandle uint32) error {
	gs.Lock()
	defer gs.Unlock()
	for i, v := range gs.groups {
		if v.serverGroupHandle == serverHandle {
			err := gs.doRemove(serverHandle)
			if err != nil {
				return err
			}
			v.Release()
			gs.groups = append(gs.groups[:i], gs.groups[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (gs *OPCGroups) doRemove(serverHandle uint32) error {
	return gs.iServer.RemoveGroup(serverHandle, true)
}

// RemoveByName Removes an OPCGroup from the collection by name
func (gs *OPCGroups) RemoveByName(name string) error {
	gs.Lock()
	defer gs.Unlock()
	for i, v := range gs.groups {
		if v.groupName == name {
			err := gs.doRemove(v.GetServerHandle())
			if err != nil {
				return err
			}
			v.Release()
			gs.groups = append(gs.groups[:i], gs.groups[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

// RemoveAll Removes all OPCGroups from the collection
func (gs *OPCGroups) RemoveAll() error {
	gs.Lock()
	defer gs.Unlock()
	for _, v := range gs.groups {
		gs.doRemove(v.GetServerHandle())
		v.Release()
	}
	gs.groups = nil
	return nil
}

// Release Releases the resources used by the collection and the items it contains.
func (gs *OPCGroups) Release() error {
	for _, group := range gs.groups {
		group.Release()
	}
	return nil
}
