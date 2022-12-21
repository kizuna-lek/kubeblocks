package authprovider

import (
	"github.com/99designs/keyring"

	"github.com/apecloud/kubeblocks/internal/cli/util"
)

const (
	PassEnv = "KUBEBLOCKS_AUTH_CREDENCIAL_PASSWORD"
)

func init() {
	keyring.Debug = false
}

type KeyRingReader interface {
	Get(string) (keyring.Item, error)
}

var kr keyring.Keyring

func keyringSetup() (keyring.Keyring, error) {
	if kr != nil {
		return kr, nil
	}
	fileDir, err := util.GetCliHomeDir()
	if err != nil {
		return nil, err
	}
	kr, err = keyring.Open(keyring.Config{
		ServiceName:              "KubeblocksOIDCAuth",
		KeychainName:             "kubeblocks-oidc-auth",
		KeychainTrustApplication: true,
		FilePasswordFunc:         fixedStringPrompt("kubeblocks-oidc-auth"),
		FileDir:                  fileDir,
		AllowedBackends:          []keyring.BackendType{keyring.FileBackend},
	})
	return kr, err
}

func fixedStringPrompt(value string) keyring.PromptFunc {
	return func(_ string) (string, error) {
		return value, nil
	}
}

// func terminalPrompt(prompt string) (string, error) {
// 	pass := os.Getenv(PassEnv)
// 	if len(pass) == 0 {
// 		fmt.Fprintf(os.Stderr, "in the future you can set %s to bypass this prompt\n", PassEnv)
// 		fmt.Fprintf(os.Stderr, "%s: ", prompt)
// 		b, err := term.ReadPassword(int(os.Stdin.Fd()))
// 		if err != nil {
// 			return "", err
// 		}
// 		pass = string(b)
// 	} else {
// 		fmt.Fprintf(os.Stderr, "using %s\n", PassEnv)
// 	}

// 	return pass, nil
// }
