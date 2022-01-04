package cs

type Entity struct {
	Health   uint32
	Team     uint32
	Position EntityPosition
	Pointer  uint32
}

type EntityPosition struct {
	X float32
	Y float32
	Z float32
}

type Color struct {
	r float32
	g float32
	b float32
}

func Init() error {
	pid, _ := findProcessPid(processName)

	var err error

	handle, err = OpenProcess(PROCESS_ALL_ACCESS, false, pid)

	if err != nil {
		return err
	}

	clientDllModuleAddress, err = getModuleBaseAddress(pid, moduleName)

	if err != nil {
		return err
	}

	return nil
}

func GetEntities() []Entity {
	var entities []Entity
	for i := 0; i < 64; i++ {
		var entityPtr uint32
		_ = ReadProcessMemory(handle, uint32(int(clientDllModuleAddress+dwEntityList)+i*0x10), &entityPtr)

		if entityPtr == 0 {
			continue
		}

		entities = append(entities, GetEntity(entityPtr))
	}

	return entities
}

func (e *Entity) EnableGlowEsp() {
	var glowManager, entityGlow uint32

	_ = ReadProcessMemory(handle, uint32(clientDllModuleAddress+dwGlowObjectManager), &glowManager)
	_ = ReadProcessMemory(handle, e.Pointer+m_iGlowIndex, &entityGlow)

	_ = WriteProcessMemory(handle, glowManager+(entityGlow*0x38)+0x8, rgbToColor(255, 255, 255))
	_ = WriteProcessMemory(handle, glowManager+(entityGlow*0x38)+0x14, float32(0.5))
	_ = WriteProcessMemory(handle, glowManager+(entityGlow*0x38)+0x27, true)
	_ = WriteProcessMemory(handle, glowManager+(entityGlow*0x38)+0x28, true)
}

func GetLocalPlayer() Entity {
	var localPlayerPtr uint32
	_ = ReadProcessMemory(handle, uint32(clientDllModuleAddress+dwLocalPlayer), &localPlayerPtr)

	return GetEntity(localPlayerPtr)
}

func GetEntity(entityPtr uint32) Entity {
	var health uint32
	var team uint32
	var position EntityPosition

	_ = ReadProcessMemory(handle, entityPtr+m_iHealth, &health)
	_ = ReadProcessMemory(handle, entityPtr+m_iTeamNum, &team)
	_ = ReadProcessMemory(handle, entityPtr+m_vecOrigin, &position)

	return Entity{health, team, position, entityPtr}
}
