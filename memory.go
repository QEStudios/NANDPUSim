package main

type MemoryRegion interface {
	Read(addr uint16) byte
	Write(addr uint16, value byte)
}

type RAM struct {
	data []byte
	base uint16
}

func NewRAM(base, size uint16) *RAM        { return &RAM{data: make([]byte, size), base: base} }
func (r *RAM) Read(addr uint16) byte       { return r.data[addr-r.base] }
func (r *RAM) Write(addr uint16, val byte) { r.data[addr-r.base] = val }

type ROM struct {
	data []byte
	base uint16
}

func NewROM(base, size uint16) *ROM        { return &ROM{data: make([]byte, size), base: base} }
func (r *ROM) Init(data []byte)            { copy(r.data, data) }
func (r *ROM) Read(addr uint16) byte       { return r.data[addr-r.base] }
func (r *ROM) Write(addr uint16, val byte) {}

type MemMap struct {
	regions []MemoryRegionEntry
}

type MemoryRegionEntry struct {
	start, end uint16
	region     MemoryRegion
}

func (m *MemMap) AddRegion(start, end uint16, region MemoryRegion) {
	m.regions = append(m.regions, MemoryRegionEntry{start, end, region})
}

func (m *MemMap) Read(addr uint16) byte {
	for _, entry := range m.regions {
		if addr >= entry.start && addr <= entry.end {
			return entry.region.Read(addr)
		}
	}
	return 0xFF // Unmapped
}

func (m *MemMap) Write(addr uint16, val byte) {
	for _, entry := range m.regions {
		if addr >= entry.start && addr <= entry.end {
			entry.region.Write(addr, val)
			return
		}
	}
}
