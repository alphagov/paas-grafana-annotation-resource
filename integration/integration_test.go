package integration_test

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = Describe("Happy path", func() {
	BeforeSuite(func() {
		buildCmd := exec.Command("docker-compose", "build")
		session, err := gexec.Start(buildCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session, 60).Should(gexec.Exit(0))
	})

	BeforeEach(func() {
		upCmd := exec.Command(
			"docker-compose", "up", "--detach", "--force-recreate",
		)
		session, err := gexec.Start(upCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session, 60).Should(gexec.Exit(0))

		for i := 1; i <= 10; i++ {
			curlCmd := exec.Command(
				"docker-compose", "exec", "-T", "resource",
				"curl",
				"-u", "admin:admin", "-sf", "-m", "10",
				"http://grafana:3000/api/health",
			)

			session, _ := gexec.Start(curlCmd, GinkgoWriter, GinkgoWriter)

			session.Wait(15 * time.Second)

			if session.ExitCode() == 0 {
				break
			}

			if i == 10 {
				Expect(session.ExitCode()).To(
					Equal(0), "Not healthy after 10 attempts",
				)
			}

			time.Sleep(1 * time.Second)
		}
	})

	It("should create an annotation", func() {
		beforeCreateCmd := exec.Command(
			"docker-compose", "exec", "-T", "resource",
			"curl",
			"-u", "admin:admin", "-sf", "-m", "10",
			"http://grafana:3000/api/annotations",
		)
		session, err := gexec.Start(beforeCreateCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(
			string(session.Wait(15 * time.Second).Out.Contents()),
		).To(Equal("[]"))

		runCreateOutCmd := exec.Command(
			"docker-compose", "exec", "-T", "-e", "BUILD_ID=12345", "resource",
			"/opt/resource/out", "/tmp",
		)
		runCreateOutCmd.Stdin = strings.NewReader(`{"source": {"url": "http://grafana:3000", "username": "admin", "password": "admin"}, "params": {}}`)
		session, err = gexec.Start(runCreateOutCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		session.Wait(15 * time.Second)

		fmt.Println(string(session.Out.Contents()))
		fmt.Println(string(session.Err.Contents()))

		Expect(
			string(session.Wait(15 * time.Second).Out.Contents()),
		).To(Equal(`{"version":{"id":"1"},"metadata":[{"name":"id","value":"1"},{"name":"tags","value":""},{"name":"text","value":"12345 nil/teams/nil/pipelines/nil/jobs/nil/builds/nil"}]}
`))

		afterCreateCmd := exec.Command(
			"docker-compose", "exec", "-T", "resource",
			"curl",
			"-u", "admin:admin", "-sf", "-m", "10",
			"http://grafana:3000/api/annotations",
		)
		session, err = gexec.Start(afterCreateCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		session.Wait(15 * time.Second)

		fmt.Println(string(session.Out.Contents()))
		fmt.Println(string(session.Err.Contents()))

		Expect(
			string(session.Out.Contents()),
		).To(SatisfyAll(
			Not(Equal("[]")),
			ContainSubstring(`"id":1`),
			ContainSubstring(`"tags":[]`),
			ContainSubstring(`"text":"12345 nil/teams/nil/pipelines/nil/jobs/nil/builds/nil"`),
			MatchRegexp(`"time":`),
		))

		runUpdateOutCmd := exec.Command(
			"docker-compose", "exec", "-T", "-e", "BUILD_ID=12345", "resource",
			"/opt/resource/out", "/",
		)
		runUpdateOutCmd.Stdin = strings.NewReader(`{"source": {"url": "http://grafana:3000", "username": "admin", "password": "admin"}, "params": {"path": "/tmp"}}`)
		session, err = gexec.Start(runUpdateOutCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		session.Wait(15 * time.Second)

		fmt.Println(string(session.Out.Contents()))
		fmt.Println(string(session.Err.Contents()))

		Expect(
			string(session.Wait(15 * time.Second).Out.Contents()),
		).To(Equal(`{"version":{"id":"1"},"metadata":[{"name":"id","value":"1"},{"name":"tags","value":""},{"name":"text","value":"12345 nil/teams/nil/pipelines/nil/jobs/nil/builds/nil"}]}
`))

		afterUpdateCmd := exec.Command(
			"docker-compose", "exec", "-T", "resource",
			"curl",
			"-u", "admin:admin", "-sf", "-m", "10",
			"http://grafana:3000/api/annotations",
		)
		session, err = gexec.Start(afterUpdateCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		session.Wait(15 * time.Second)

		fmt.Println(string(session.Out.Contents()))
		fmt.Println(string(session.Err.Contents()))

		Expect(
			string(session.Out.Contents()),
		).To(SatisfyAll(
			Not(Equal("[]")),
			ContainSubstring(`"id":1`),
			ContainSubstring(`"tags":[]`),
			ContainSubstring(`"text":"12345 nil/teams/nil/pipelines/nil/jobs/nil/builds/nil"`),
			MatchRegexp(`"time":`),
		))

		runCheckCmd := exec.Command(
			"docker-compose", "exec", "-T", "resource",
			"/opt/resource/check", "/",
		)
		runCheckCmd.Stdin = strings.NewReader(`{"source": {"url": "http://grafana:3000", "username": "admin", "password": "admin"}}`)
		session, err = gexec.Start(runCheckCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())

		session.Wait(15 * time.Second)

		fmt.Println(string(session.Out.Contents()))
		fmt.Println(string(session.Err.Contents()))

		Expect(
			string(session.Wait(15 * time.Second).Out.Contents()),
		).To(Equal(`[{"id":"1"}]
`))
	})

	AfterEach(func() {
		upCmd := exec.Command(
			"docker-compose", "down",
		)
		session, err := gexec.Start(upCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(session, 60).Should(gexec.Exit(0))
	})

	AfterSuite(func() {
		gexec.KillAndWait()
	})
})
