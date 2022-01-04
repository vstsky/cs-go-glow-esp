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

func GetEntity(entityPtr uint32) Entity {
	var health uint32
	var team uint32
	var position EntityPosition

	_ = ReadProcessMemory(handle, entityPtr+m_iHealth, &health)
	_ = ReadProcessMemory(handle, entityPtr+m_iTeamNum, &team)
	_ = ReadProcessMemory(handle, entityPtr+m_vecOrigin, &position)

	return Entity{health, team, position, entityPtr}
}
