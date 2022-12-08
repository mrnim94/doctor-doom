package doom

import (
	"os"
	"testing"
)

func prepare() {
	// Create folder structure like this:
	// ./tmp
	// - text.txt
	// - /folder1
	//   - text.txt
	//   - /folder2
	//     - text.txt

	// Create root folder
	os.Mkdir("./tmp", 0755)

	// Create file in root folder
	os.Create("./tmp/text.txt")

	// Create folder1
	os.Mkdir("./tmp/folder1", 0755)

	// Create file in folder1
	os.Create("./tmp/folder1/text.txt")

	// Create folder2
	os.Mkdir("./tmp/folder1/folder2", 0755)

	// Create file in folder2
	os.Create("./tmp/folder1/folder2/not_text__.txt")
}

func cleanup() {
	// Remove folder structure
	os.RemoveAll("./tmp")
}

func TestGetDoomVictims(t *testing.T) {

	prepare()
	defer cleanup()

	// DoomOptions
	doomOptions := DoomOptions{
		DoomPath:   "./tmp",
		Circle:     "0 0 0 * * 0", // Every Sunday at 00:00:00
		DoomExport: "/var/log",    // Export log path folder file will be named as doom-*.log
		Rule: DoomDestroyRules{
			Age:  "1s",
			Size: "0B",
			Name: "text.txt",
		},
	}

	// Doom
	doom := DoctorDoom{
		DoomOptions: doomOptions,
	}

	// Run
	victims := doom.GetDoomVictims()

	if len(victims) != 3 {
		t.Errorf("Expected 1 victim, got %d", len(victims))
	}
}

func BenchmarkGetDoomVictims(b *testing.B) {

	prepare()
	defer cleanup()

	// DoomOptions
	doomOptions := DoomOptions{
		DoomPath:   "./tmp",
		Circle:     "0 0 0 * * 0", // Every Sunday at 00:00:00
		DoomExport: "/var/log",    // Export log path folder file will be named as doom-*.log
		Rule: DoomDestroyRules{
			Age:  "1s",
			Size: "0B",
			Name: "text.txt",
		},
	}

	// Doom
	doom := DoctorDoom{
		DoomOptions: doomOptions,
	}

	// Run
	for i := 0; i < b.N; i++ {
		doom.GetDoomVictims()
	}
}

func TestDoomDestroy(t *testing.T) {

	prepare()
	defer cleanup()

	// DoomOptions
	doomOptions := DoomOptions{
		DoomPath:   "./tmp",
		Circle:     "0 0 0 * * 0", // Every Sunday at 00:00:00
		DoomExport: "/var/log",    // Export log path folder file will be named as doom-*.log
		Rule: DoomDestroyRules{
			Age:  "1s",
			Size: "0B",
			Name: "text.txt",
		},
	}

	// Doom
	doom := DoctorDoom{
		DoomOptions: doomOptions,
	}

	// Run
	victims := doom.GetDoomVictims()
	doom.DestroyDoomVictims(victims)

	if len(victims) != 3 {
		t.Errorf("Expected 1 victim, got %d", len(victims))
	}

	checkVictims := doom.GetDoomVictims()
	if len(checkVictims) != 0 {
		t.Errorf("Expected 0 victim, got %d", len(checkVictims))
	}
}
