package lwconfig

import "testing"

func TestProfile_Verify(t *testing.T) {
	tests := []struct {
		profile Profile
		desc    string
		error   bool
	}{
		{
			Profile{
				Account:   "test.test.corp",
				ApiKey:    "0000000000000000000000000000000000000000000000000000000",
				ApiSecret: "000000000000000000000000000000",
			},
			"fields present and valid",
			false,
		},
		{
			Profile{
				Account:   "test.test.corp.lacework.net",
				ApiKey:    "0000000000000000000000000000000000000000000000000000000",
				ApiSecret: "000000000000000000000000000000",
			},
			"fields present and valid full URL",
			false,
		},
		{
			Profile{
				Account:   "",
				ApiKey:    "0000000000000000000000000000000000000000000000000000000",
				ApiSecret: "000000000000000000000000000000",
			},
			"account empty",
			true,
		},
		{
			Profile{
				Account:   "test.test.corp.lacework.net",
				ApiKey:    "",
				ApiSecret: "000000000000000000000000000000",
			},
			"apikey empty",
			true,
		},
		{
			Profile{
				Account:   "test.test.corp.lacework.net",
				ApiKey:    "0000000000000000000000000000000000000000000000000000000",
				ApiSecret: "",
			},
			"apisecret empty",
			true,
		},
		{
			Profile{
				Account:   "test.test.corp.lacework.net",
				ApiKey:    "000000000000000000000000000000000000000000000000000000",
				ApiSecret: "000000000000000000000000000000",
			},
			"apikey length less than 55 characters",
			true,
		},
		{
			Profile{
				Account:   "test.test.corp.lacework.net",
				ApiKey:    "0000000000000000000000000000000000000000000000000000000",
				ApiSecret: "00000000000000000000000000000",
			},
			"api secret length less than 30 characters",
			true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.profile.Verify()
			if (err != nil) != tc.error {
				t.Error("Incorrect result")
			}
		})
	}

}
