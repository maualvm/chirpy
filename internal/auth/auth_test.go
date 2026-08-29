package auth

import "testing"

func TestCheckPasswordHash(t *testing.T) {
	password1 := "verySecurePASSWORD!123"
	password2 := "some-sEcur3-p4ssword"
	hashed1, _ := HashPassword(password1)
	hashed2, _ := HashPassword(password2)

	testCases := []struct {
		name            string
		password        string
		hash            string
		wantErr         bool
		matchesPassword bool
	}{
		{
			name:            "Correct password",
			password:        password1,
			hash:            hashed1,
			wantErr:         false,
			matchesPassword: true,
		},
		{
			name:            "Incorrect password",
			password:        "clearlyNotThePassword",
			hash:            hashed1,
			wantErr:         false,
			matchesPassword: false,
		},
		{
			name:            "Password does not match another hash",
			password:        password1,
			hash:            hashed2,
			wantErr:         false,
			matchesPassword: false,
		},
		{
			name:            "Password is empty",
			password:        "",
			hash:            hashed1,
			wantErr:         false,
			matchesPassword: false,
		},
		{
			name:            "Invalid hash",
			password:        password1,
			hash:            "invalidhash",
			wantErr:         true,
			matchesPassword: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			matched, err := CheckPasswordHash(testCase.password, testCase.hash)
			if err != nil && !testCase.wantErr {
				t.Errorf("CheckPasswordHash() errored: %v, wantErr %v", err, testCase.wantErr)
			}

			if !testCase.wantErr && matched != testCase.matchesPassword {
				t.Errorf("CheckPasswordHash() expected %v, got %v", testCase.matchesPassword, matched)
			}
		})
	}
}
