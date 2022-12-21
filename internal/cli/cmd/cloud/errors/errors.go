package errors

import (
	"github.com/spf13/cobra"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func CheckErr(err error) {
	switch err.(type) {
	case *apierrors.StatusError:
		cmdutil.CheckErr(err)
	default:
		cobra.CheckErr(err)
	}
}
