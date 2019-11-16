package debug

import (
	"errors"
	"debug/elf"
)

func GetSymbolOffset(cmd, symbol string) (uint64, error) {
	elfFile, err := elf.Open(cmd)
	if err != nil {
		return 0, err
	}
	defer elfFile.Close()

	symbols, err := elfFile.DynamicSymbols()
	if err != nil {
		return 0, err
	}
	var offset uint64
	found := false
	for _, sym := range symbols {
		if sym.Name == symbol {
			found = true
			offset = sym.Value
			break
		}
	}
	if !found {
		return 0, errors.New("symbol not found")
	}
	base, err := getBaseAddress(elfFile)
	return offset - base, err
}

func getBaseAddress(f *elf.File) (uint64, error) {
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_LOAD {
			return prog.Vaddr, nil
		}
	}
	return 0, errors.New("failed to get base address")
}
