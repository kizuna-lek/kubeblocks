package migration

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdtesting "k8s.io/kubectl/pkg/cmd/testing"
)

var _ = Describe("list", func() {

	var (
		streams genericclioptions.IOStreams
		tf      *cmdtesting.TestFactory
	)

	It("command build", func() {
		cmd := NewMigrationListCmd(tf, streams)
		Expect(cmd).ShouldNot(BeNil())
	})

})
