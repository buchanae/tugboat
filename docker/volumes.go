package docker

// AddVolume adds a mapped volume to the stage. A corresponding Volume record
// is added to stage.Volumes.
//
// If the volume paths are invalid or can't be mapped, an error is returned.
/*
func AddVolume(hostPath string, mountPoint string, readonly bool) error {
	vol := Volume{
		HostPath:      hostPath,
		ContainerPath: mountPoint,
		Readonly:      readonly,
	}

	for i, v := range stage.Volumes {
		// check if this volume is already present in the stage
		if vol == v {
			return nil
		}

		// If the proposed RW Volume is a subpath of an existing RW Volume
		// do not add it to the stage
		// If an existing RW Volume is a subpath of the proposed RW Volume, replace it with
		// the proposed RW Volume
		if !vol.Readonly && !v.Readonly {
			if stage.IsSubpath(vol.ContainerPath, v.ContainerPath) {
				return nil
			} else if stage.IsSubpath(v.ContainerPath, vol.ContainerPath) {
				stage.Volumes[i] = vol
				return nil
			}
		}
	}

	stage.Volumes = append(mapper.Volumes, vol)
	return nil
}
*/
