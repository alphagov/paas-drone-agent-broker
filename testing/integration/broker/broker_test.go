package broker_test

import (
	"code.cloudfoundry.org/lager"
	"github.com/alphagov/paas-go/provider"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	"strings"

	"encoding/json"
	"net/http"
	"os"

	"github.com/alphagov/paas-go/broker"
	brokertesting "github.com/alphagov/paas-go/testing/broker"
	"github.com/pivotal-cf/brokerapi"
	provideriface "github.com/richardTowers/paas-drone-agent-broker/provider"
	"net/http/httptest"
)

const (
	ASYNC_ALLOWED = true
)

type BindingResponse struct {
	Credentials map[string]interface{} `json:"credentials"`
}

var _ = Describe("Broker", func() {

	PIt("should manage the lifecycle of a drone agent", func() {
		var (
			instanceID = uuid.NewV4().String()
			binding1ID = uuid.NewV4().String()
			binding2ID = uuid.NewV4().String()
			serviceID  = "uuid-1"
			planID     = "uuid-2"
		)

		By("initialising")
		_, brokerTester := initialise()

		By("Provisioning")
		res := brokerTester.Provision(instanceID, brokertesting.RequestBody{
			ServiceID: serviceID,
			PlanID:    planID,
		}, ASYNC_ALLOWED)
		Expect(res.Code).To(Equal(http.StatusCreated))

		By("Binding an app")
		res = brokerTester.Bind(instanceID, binding1ID, brokertesting.RequestBody{
			ServiceID: serviceID,
			PlanID:    planID,
		}, ASYNC_ALLOWED)
		Expect(res.Code).To(Equal(http.StatusCreated))

		By("Asserting the credentials returned work")
		binding1Creds := extractCredentials(res)
		Expect(binding1Creds).To(Equal("todo: check that the credentials work"))

		By("Binding another app")
		res = brokerTester.Bind(instanceID, binding2ID, brokertesting.RequestBody{
			ServiceID: serviceID,
			PlanID:    planID,
		}, ASYNC_ALLOWED)
		Expect(res.Code).To(Equal(http.StatusCreated))

		By("Asserting the credentials returned work")
		binding2Creds := extractCredentials(res)
		Expect(binding2Creds).To(Equal("todo: check that the credentials work"))

		By("Asserting the first user's credentials still work")
		Expect("todo").To(Equal("check that the credentials work"))

		By("Unbinding the first app")
		res = brokerTester.Unbind(instanceID, serviceID, planID, binding1ID, ASYNC_ALLOWED)
		Expect(res.Code).To(Equal(http.StatusOK))

		By("Asserting the second user's credentials still work")
		Expect("todo").To(Equal("check that the credentials work"))

		By("Unbinding the second app")
		res = brokerTester.Unbind(instanceID, serviceID, planID, binding2ID, ASYNC_ALLOWED)
		Expect(res.Code).To(Equal(http.StatusOK))

		By("Deprovisioning")
		res = brokerTester.Deprovision(instanceID, serviceID, planID, ASYNC_ALLOWED)
		Expect(res.Code).To(Equal(http.StatusOK))

		By("Returning a 410 response when trying to delete a non-existent instance")
		res = brokerTester.Deprovision(instanceID, serviceID, planID, ASYNC_ALLOWED)
		Expect(res.Code).To(Equal(http.StatusGone))
	})
})

func initialise() (*provider.ServiceProvider, brokertesting.BrokerTester) {
	configFile := strings.NewReader(`{
		"basic_auth_username": "username",
		"basic_auth_password": "password",
		"catalog": {
			"services": [{
					"id": "uuid-1",
					"name": "Drone Agent",
					"description": "https://drone.io build agent",
					"bindable": true,
					"plan_updateable": false,
					"requires": [],
					"metadata": {},
					"plans": [{
							"id": "uuid-2",
							"name": "1x-small",
							"description": "A single drone agent running on a small virtual machine",
							"metadata": {}
					}]
			}]
		}
	}`)

	config, err := broker.NewConfig(configFile)
	Expect(err).ToNot(HaveOccurred())

	droneAgentProvider, err := provideriface.NewDroneAgentProvider(config.Provider)
	Expect(err).ToNot(HaveOccurred())

	logger := lager.NewLogger("drone-agent-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, config.API.LagerLogLevel))

	serviceBroker := broker.New(config, droneAgentProvider, logger)
	brokerAPI := broker.NewAPI(serviceBroker, logger, config)

	return &droneAgentProvider, brokertesting.New(brokerapi.BrokerCredentials{
		Username: "username",
		Password: "password",
	}, brokerAPI)
}

func extractCredentials(res *httptest.ResponseRecorder) interface{} {
	parsedResponse := BindingResponse{}
	err := json.NewDecoder(res.Body).Decode(&parsedResponse)
	Expect(err).ToNot(HaveOccurred())
	// Ensure returned credentials follow guidlines in https://docs.cloudfoundry.org/services/binding-credentials.html
	var str string
	creds := parsedResponse.Credentials
	Expect(creds).To(HaveKeyWithValue("password", BeAssignableToTypeOf(str)))
	return nil
}
