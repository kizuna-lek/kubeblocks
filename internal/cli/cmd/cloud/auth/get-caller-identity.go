package auth

import (
	"encoding/json"
	"fmt"

	jwt "github.com/go-jose/go-jose/v3"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	authprovider "github.com/apecloud/kubeblocks/internal/cli/cmd/cloud/auth/provider"
)

type callerIdentityOptions struct {
	genericclioptions.IOStreams
}

func newGetCallerIdentityCmd(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := &callerIdentityOptions{
		IOStreams: streams,
	}

	cmd := &cobra.Command{
		Use:   "get-caller-identity",
		Short: "Returns details about the credentialed accounts with username and organization info",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd))
			cmdutil.CheckErr(o.Run(f, cmd))
		},
	}
	cmd.AddCommand(newLoginCmd(f, streams))
	return cmd
}

func (o *callerIdentityOptions) Complete(f cmdutil.Factory, cmd *cobra.Command) error {
	return nil
}

type token struct {
	IDToken string `json:"id_token"`
}

type identity struct {
	UserID string   `json:"UserId"`
	Email  string   `json:"Email"`
	Groups []string `json:"Groups"`
}

func (i *identity) GetUID() string {
	return i.UserID
}

func (i *identity) GetEmail() string {
	return i.Email
}

func (i *identity) GetGroups() []string {
	return i.Groups
}

func (i *identity) String() string {
	bt, _ := json.MarshalIndent(i, "", "    ")
	return string(bt)
}

type Identity interface {
	GetUID() string
	GetEmail() string
	GetGroups() []string
}

type claims map[string]json.RawMessage

func (c claims) unmarshalClaim(name string, v interface{}) error {
	val, ok := c[name]
	if !ok {
		return fmt.Errorf("claim not present")
	}
	return json.Unmarshal([]byte(val), v)
}

func GetCallerIdentity() (Identity, error) {

	var reader authprovider.TokenReader

	reader, err := authprovider.GetTokenStore()

	if err != nil {
		return nil, err
	}

	tokenData, err := getTokenFromCache(reader)
	if err != nil {
		return nil, fmt.Errorf("get token failed, err=%s", err.Error())
	}
	var tk token
	err = json.Unmarshal([]byte(tokenData), &tk)
	if err != nil {
		return nil, err
	}

	tok, err := jwt.ParseSigned(tk.IDToken)
	if err != nil {
		return nil, err
	}

	var c claims
	payload := tok.UnsafePayloadWithoutVerification()

	err = json.Unmarshal(payload, &c)
	if err != nil {
		return nil, err
	}

	var id identity

	if err := c.unmarshalClaim("kubeblocks.io/uid", &id.UserID); err != nil {
		return nil, err
	}
	if err := c.unmarshalClaim("kubeblocks.io/email", &id.Email); err != nil {
		return nil, err
	}
	if err := c.unmarshalClaim("kubeblocks.io/groups", &id.Groups); err != nil {
		return nil, err
	}

	return &id, nil

}

func (o *callerIdentityOptions) Run(f cmdutil.Factory, cmd *cobra.Command) error {

	id, err := GetCallerIdentity()
	if err != nil {
		return err
	}

	fmt.Fprintln(o.IOStreams.Out, id)
	return nil
}
