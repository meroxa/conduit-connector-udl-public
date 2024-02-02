package destination

import "fmt"

const udlTimeLayout = "06002150405.000"

type UDLReport struct {
	ID      string
	Entries []UDLEntry
}

type UDLEntry struct {
	Timestamp string
	Position  UDLPosition
	Velocity  UDLVelocity
}

type UDLPosition struct {
	// X coords in km
	X string
	// Y coords in km
	Y string
	// Z coords in km
	Z string
}

type UDLVelocity struct {
	// X velocity in km/s
	X string
	// Y velocity in km/s
	Y string
	// Z velocity in km/s
	Z string
}

func (r UDLReport) String() string {
	var out string
	for _, e := range r.Entries {
		out += fmt.Sprintf("%s %s %s %s %s %s %s\n",
			e.Timestamp,
			e.Position.X,
			e.Position.Y,
			e.Position.Z,
			e.Velocity.X,
			e.Velocity.Y,
			e.Velocity.Z)
	}
	return out
}
