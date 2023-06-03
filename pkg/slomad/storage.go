package slomad

type Volume struct {
	Src   string
	Dst   string
	Mount bool
}

func NewDockerVolume(src, dst string) *Volume {
	return &Volume{
		Src: src,
		Dst: dst,
	}
}

func NewNomadVolume(src, dst string) *Volume {
	return &Volume{
		Src:   src,
		Dst:   dst,
		Mount: true,
	}
}

type StorageParams struct {
	Storage *string
	Volumes []Volume
	Mounts  []Volume
}
