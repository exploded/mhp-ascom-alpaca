package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

const NumOnOffSwitch = 8 // Number of on/off switches
const NumVarSwitch = 4   // Number of variable switches
const NumSwitches = 12   // Number of all switches in total

// type swm sync.RWMutex
type sw struct {
	Connected           bool       `json:"connected"`
	Focusermaxincrement int32      `json:"focucermaxincrement"`
	Focusermaxstep      int32      `json:"focucermaxstep"`
	Focucerposition     int32      `json:"focucerposition"`
	Focucerspeed        int32      `json:"focucerspeed"`
	Name                [13]string `json:"name"`       // *
	Devicetype          [13]string `json:"devicetype"` //*
	Number              [13]uint32 `json:"number"`     //**
	Uniqueid            [13]string `json:"uniqueid"`   //**
	Id                  [13]uint32 `json:"id"`
	Customname          [13]string `json:"customname"`
	Min                 [13]int64  `json:"min"`
	Max                 [13]int64  `json:"max"`
	Step                [13]int64  `json:"step"`
	Canwrite            [13]bool   `json:"canwrite"`
	Value               [13]int64  `json:"value"`
}

var s = &sw{}
var sm sync.RWMutex

func MhpSetInit() {
	s.mhpsetinit()
}

func (s *sw) mhpsetinit() {
	if !s.mhpLoadSettings() {
		sm.Lock()
		// set up the defaults
		s.Connected = false
		s.Focusermaxincrement = 150
		s.Focusermaxstep = 65535
		s.Focucerposition = 1000
		s.Focucerspeed = 50 // Range is 0 to 100%
		s.Name = [13]string{"Focuser", "Switch 1", "Switch 2", "Switch 3", "Switch 4", "Switch 5", "Switch 6", "Switch 7", "Switch 8", "Dew Heater 1", "Dew Heater 2", "Dew Heater 3", "Dew Heater 4"}
		s.Devicetype = [13]string{"Focuser", "Switch", "Switch", "Switch", "Switch", "Switch", "Switch", "Switch", "Switch", "Switch", "Switch", "Switch", "Switch"}
		s.Number = [13]uint32{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		s.Id = [13]uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		s.Customname = [13]string{"", "", "", "", "", "", "", "", "", "", "", "", ""}
		s.Min = [13]int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		s.Max = [13]int64{65535, 1, 1, 1, 1, 1, 1, 1, 1, 100, 100, 100, 100}
		s.Step = [13]int64{150, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		s.Canwrite = [13]bool{true, true, true, true, true, true, true, true, true, true, true, true, true}
		s.Value = [13]int64{1000, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		s.Uniqueid = [13]string{"6fd5bae2-40ed-489f-b6f3-a562822e48e9", "86c4b6ea-650d-45cd-ad5d-1771c86edee6", "b96a0f0d-3b3f-4240-a7dc-807645a91a9a", "5cf95480-14ed-49c6-b992-a5eb8c4c9fb2", "9e2090fa-a793-4d4e-9302-3d97ba5566d2", "16eef02f-e1f0-4b94-8a66-a45ca005246f", "8d5641ce-34d4-4750-a0d7-210603a4ea33", "177ed90e-f8f4-46c3-8bd3-b00e4c7dedb5", "c0babc9b-4403-4b8f-9d5a-f53a848f7aa2", "96730903-e921-4a0d-8f45-f76597cf6259", "0c466cbc-a363-40fb-825c-07e33f0c696f", "41ce27ac-de5c-472a-a2a2-c37e3490c627", "5e95431e-6a38-4cd3-8cc0-65dfdb087e82"}
		sm.Unlock()
		s.mhpSaveSettings()
	}
}

func (s *sw) mhpSaveSettings() {
	sm.Lock()
	defer sm.Unlock()

	data, err := json.MarshalIndent(&s, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("settings.json", data, 0644) //os.ModeExclusive)
	if err != nil {
		panic(err)
	}
}

func (s *sw) mhpLoadSettings() (result bool) {
	sm.Lock()
	defer sm.Unlock()
	data, err := os.ReadFile("settings.json")
	if err != nil {
		return false
	}

	err = json.Unmarshal(data, &s)
	return err == nil
}

func MhpGetInit() []DeviceConfiguration {
	return s.mhpgetinit()
}

func (s *sw) mhpgetinit() []DeviceConfiguration {
	sm.Lock()
	defer sm.Unlock()
	var val []DeviceConfiguration
	for i := range s.Name {
		val = append(val, DeviceConfiguration{
			DeviceName:   s.Name[i],
			DeviceType:   s.Devicetype[i],
			DeviceNumber: s.Number[i], //s.Id[i],
			UniqueID:     s.Uniqueid[i],
		})
	}
	return val
}

func MhpSetName(id int32, CustomName string) (err error) {
	if id < 0 || id >= NumSwitches {
		err = errors.New("invalid device number")
	}
	s.setname(id, CustomName)
	return
}

func (s *sw) setname(id int32, CustomName string) {
	sm.Lock()
	s.Customname[id] = CustomName
	sm.Unlock()
	s.mhpSaveSettings()
}

func MhpSetConnect(c bool) {
	s.setconnect(c)
}

func (s *sw) setconnect(c bool) {
	sm.Lock()
	s.Connected = c
	sm.Unlock()
	s.mhpSaveSettings()
}

func MhpGetConnected() bool {
	return s.getconnected()
}

func (s *sw) getconnected() bool {
	sm.Lock()
	defer sm.Unlock()
	return s.Connected // needs to be fixed
}

func MhpGetName(id int32) string {
	return s.getname(id)
}

func (s *sw) getname(id int32) string {
	sm.Lock()
	defer sm.Unlock()
	if s.Customname[id] != "" {
		return s.Customname[id]
	}
	return s.Name[id]
}

func MhpGetType(id int32) string {
	return s.gettype(id)
}

func (s *sw) gettype(id int32) string {
	sm.Lock()
	defer sm.Unlock()
	return s.Devicetype[id]
}

func MhpGetNumber(id uint32) uint32 {
	return s.getnumber(id)
}

func (s *sw) getnumber(id uint32) uint32 {
	sm.Lock()
	defer sm.Unlock()
	return s.Number[id]
}

func MhpGetUniqueID(id int32) string {
	return s.getuniqueid(id)
}

func (s *sw) getuniqueid(id int32) string {
	sm.Lock()
	defer sm.Unlock()
	return s.Uniqueid[id]
}

func MhpGetOnOff(id int32) (result bool, err error) {
	return s.getonoff(id)
}

func (s *sw) getonoff(id int32) (result bool, err error) {
	sm.Lock()
	defer sm.Unlock()
	if s.Max[id] > 1 {
		err = errors.New("device is not just an on off switch")
		return
	} else {
		if s.Value[id] == 0 {
			result = false
		} else {
			result = true
		}
	}
	return
}

func MhpGetValue(id int32) int64 {
	return s.getvalue(id)
}

func (s *sw) getvalue(id int32) int64 {
	sm.Lock()
	defer sm.Unlock()
	return s.Value[id]
}

func MhpGetMax(id int32) int64 {
	return s.getmax(id)
}

func (s *sw) getmax(id int32) int64 {
	sm.Lock()
	defer sm.Unlock()
	return s.Max[id]
}

func MhpGetMin(id int32) int64 {
	return s.getmin(id)
}

func (s *sw) getmin(id int32) int64 {
	sm.Lock()
	defer sm.Unlock()
	return s.Min[id]
}

func MhpGetStep(id int32) int64 {
	return s.getstep(id)
}

func (s *sw) getstep(id int32) int64 {
	sm.Lock()
	defer sm.Unlock()
	return s.Step[id]
}

// Function sends the command to set the 4 variable switches (i.e. dew heater controllers).
// id is from 8 to 11. range / value is from 0 to 100 (0x00 to 0x64)
func MhpSetValue(id int32, value int64) (err error) {
	if id < 1 || id > NumSwitches {
		err = errors.New("invalid switch number")
		return
	}

	if value < 0 || value > s.Max[id] { //100
		err = errors.New("invalid switch level")
		return
	}

	// Check for special case of on/off switches
	if id >= 1 && id <= NumOnOffSwitch {
		err = MhpSetOnOff(id, value == 1)
		return
	}
	// Case for dew heaters
	return s.setvalue(id, value)
}

func (s *sw) setvalue(id int32, value int64) (err error) {
	// Examples				Hex     Decimal
	// Switch 9 to 0 		4b 00	75 00
	// Switch 9 to 50 		4b 32	75 50
	// Switch 9 to 100 		4b 64	75 100
	// ...
	// Switch 10 to 0        4a 00  74 00
	// ...
	// Switch 12 to 00		48 00	72 00
	// Switch 12 to 100		48 64	72 100

	// Value is in the 2 most significant digits Switch number is the 2 lest significant 2 hex digits.
	var command int32 = (int32(value) * 0x100) + (0x48 + 12 - int32(id))
	log.Println("Set dew heater no:", id-8, " to:", value)
	err = hidSend(int64(command))
	if err != nil {
		return err
	}

	sm.Lock()
	s.Value[id] = value
	sm.Unlock()
	s.mhpSaveSettings()
	return
}

// Function returns the command to turn the 8 on/off switches on or off. id is from 0 to 7
func MhpSetOnOff(id int32, state bool) (err error) {
	// Examples
	// Switch 0 on 100  (0x64)
	// Switch 0 off 99	(0x63)
	// ...
	// Switch 7 on 86	(0x56)
	// Switch 7 off 85	(0x55)
	var command int32
	if id < 1 || id > NumOnOffSwitch {
		err = errors.New("invalid switch number")
		return
	}
	command = 0x55 + (8-id)*2
	// If the switch is to be turned on, add 1
	if state {
		command++
	}
	err = hidSend(int64(command))
	if err != nil {
		return err
	}
	err = s.setonoff(id, state)
	return
}

func (s *sw) setonoff(id int32, state bool) (err error) {
	sm.Lock()
	if state {
		s.Value[id] = 1
	} else {
		s.Value[id] = 0
	}
	sm.Unlock()
	s.mhpSaveSettings()
	// fmt.Println("SetOnOff ID", id, " State ", state)
	return
}

// Move the focuser
func MhpMove(value int32) (err error) {
	if value < 0 || value > s.Focusermaxstep {
		err = errors.New("invalid focuser position")
		return
	}
	// Move the focuser
	return s.mhpmove(value)
}

func (s *sw) mhpmove(value int32) (err error) {
	// Examples				Hex
	// In 1 step 50% speed		4e 8e 00 01
	// Out 1 step 50% speed				4c 8e 00 01
	// In 2 steps 50% speed				4e 8e 00 02
	// ...
	// In 999 steps 50%	 speed			4e 8e 03 e7
	// Out 999 steps 50% speed			4c 8e 03 e7

	// Assume 'IN' direction is negative
	var part1, part2 int64

	sm.Lock()
	current := int64(s.Focucerposition)
	// Get the speed in % range from 0% to 100% from settings and convert to range 250 to 35 DEC (i.e. 0xFA to 0x35)
	speed := int64((25050 - (s.Focucerspeed * 215)) / 100)
	sm.Unlock()

	switch {
	case int64(value) == current:
		err = errors.New("no focuser movement requested")
		return
	case int64(value) < current: // In
		part1 = speed*0x100 + 0x4e // example 0x8e4e 8e = speed 50%,  4e = in
		part2 = current - int64(value)
	case int64(value) > current: // Out
		part1 = speed*0x100 + 0x4c // example  0x8e4c 8e = speed 50%, 4c = out
		part2 = int64(value) - current
	}

	value1 := int64(part2 / 0x100)
	value2 := part2 - (value1 * 0x100)

	// Value is in the 2 most significant digits Switch number is the 2 lest significant 2 hex digits.
	var command int64 = (value2 * 0x1000000) + (value1 * 0x10000) + part1

	log.Println("Value (new position):", value, "Move Steps (+ve is out,-ve is in):", int64(value)-current, "Speed ", s.Focucerspeed)
	err = hidSend(command)
	if err != nil {
		return err
	}

	sm.Lock()
	s.Focucerposition = value
	sm.Unlock()
	s.mhpSaveSettings()
	return
}

func MhpGetMaxStep() int32 {
	return s.getmaxstep()
}

func (s *sw) getmaxstep() int32 {
	sm.Lock()
	defer sm.Unlock()
	return s.Focusermaxstep
}

func MhpGetMaxIncrement() int32 {
	return s.getmaxincrement()
}

func (s *sw) getmaxincrement() int32 {
	sm.Lock()
	defer sm.Unlock()
	return s.Focusermaxincrement
}

func MhpGetPosition() int32 {
	return s.getposition()
}

func (s *sw) getposition() int32 {
	sm.Lock()
	defer sm.Unlock()
	return s.Focucerposition
}
