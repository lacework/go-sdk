package lwupdater

import (
	"testing"
	"time"
)

var now = time.Now()

func sampleVersion() Version {
	return Version{
		Project:        "default-project",
		CurrentVersion: "v1.0.0",
		LatestVersion:  "v1.0.0",
		LastCheckTime:  now,
		Outdated:       false,
		ComponentsLastCheck: map[string]time.Time{
			"sca": now,
			"iac": now,
		},
	}
}

func noComponentsVersion() Version {
	return Version{
		Project:             "default-project",
		CurrentVersion:      "v1.0.0",
		LatestVersion:       "v1.0.0",
		LastCheckTime:       time.Now().AddDate(0, 0, -2),
		Outdated:            false,
		ComponentsLastCheck: nil,
	}
}

func TestVersion_CheckComponentBefore(t *testing.T) {
	type args struct {
		component string
		checkTime time.Time
	}
	tests := []struct {
		name  string
		cache Version
		args  args
		want  bool
	}{
		{
			name:  "missing0",
			cache: sampleVersion(),
			args: args{
				component: "sast",
				checkTime: now,
			},
			want: true,
		},
		{
			name:  "missing1",
			cache: noComponentsVersion(),
			args: args{
				component: "sca",
				checkTime: now,
			},
			want: true,
		},
		{
			name:  "should not check",
			cache: sampleVersion(),
			args: args{
				component: "sca",
				checkTime: now.AddDate(0, 0, -2),
			},
			want: false,
		},
		{
			name:  "should check",
			cache: sampleVersion(),
			args: args{
				component: "iac",
				checkTime: now.AddDate(0, 0, 2),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cache.CheckComponentBefore(tt.args.component, tt.args.checkTime); got != tt.want {
				t.Errorf("CheckComponentBefore() = %v, want %v", got, tt.want)
			}
		})
	}
}
