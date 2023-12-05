package icon

import (
	"github.com/hugobyte/dive-core/cli/common"
	"github.com/spf13/cobra"
)

const DefaultIconGenesisFile = "../../static-files/config/genesis-icon-0.zip"

var (
	serviceName    = ""
	ksPath         = ""
	ksPassword     = ""
	networkID      = ""
	nodeEndpoint   = ""
	genesis        = ""
	configFilePath = ""
)

var IconCmd = common.NewDiveCommandBuilder().
	SetUse("icon").
	SetShort("Build, initialize and start a icon node.").
	SetLong(`The command starts an Icon node, initiating the process of setting up and launching a local Icon network.
It establishes a connection to the Icon network and allows the node in executing smart contracts and maintaining the decentralized ledger.`).
	AddCommand(IconDecentralizeCmd).
	AddStringFlagWithShortHand(&genesis, "genesis", "g", "", "path to custom genesis file").
	AddStringFlagWithShortHand(&configFilePath, "config", "c", "", "path to custom config json file").
	AddBoolFlagP("decentralization", "d", false, "decentralize Icon Node").
	SetRun(icon).
	Build()

var IconDecentralizeCmd = common.NewDiveCommandBuilder().
	SetUse("decentralize").
	SetShort("Decentralize already running Icon Node").
	SetLong(`Decentralize Icon Node is necessary if you want to connect your local icon node to BTP network`).
	AddStringFlagWithShortHand(&serviceName, "serviceName", "s", "", "service name").
	AddStringFlagWithShortHand(&nodeEndpoint, "nodeEndpoint", "e", "", "endpoint address").
	AddStringFlagWithShortHand(&ksPath, "keystorePath", "k", "", "keystore path").
	AddStringFlagWithShortHand(&ksPassword, "keyPassword", "p", "", "keypassword").
	AddStringFlagWithShortHand(&networkID, "nid", "n", "", "NetworkId of Icon Node").
	MarkFlagsAsRequired([]string{"serviceName", "nodeEndpoint", "keystorePath", "keyPassword", "nid"}).
	SetRun(iconDecentralization).
	Build()

func icon(cmd *cobra.Command, args []string) {

	cliContext := common.GetCliWithKurtosisContext()

	err := common.ValidateArgs(args)
	if err != nil {
		cliContext.Logger().Fatal(common.CodeOf(err), err.Error())
	}

	decentralization, err := cmd.Flags().GetBool("decentralization")
	if err != nil {
		cliContext.Logger().Error(common.InvalidCommandError, err.Error())
	}

	var response = &common.DiveServiceResponse{}

	if decentralization {
		response, err = RunIconNode(cliContext)

		if err != nil {
			cliContext.Logger().Error(common.CodeOf(err), err.Error())
			cliContext.Spinner().Stop()
		}
		params := GetDecentralizeParams(response.ServiceName, response.PrivateEndpoint, response.KeystorePath, response.KeyPassword, response.NetworkId)

		err = RunDecentralization(cliContext, params)

		if err != nil {
			cliContext.Logger().Error(common.CodeOf(err), err.Error())
			cliContext.Spinner().Stop()
		}

	} else {
		response, err = RunIconNode(cliContext)

		if err != nil {
			cliContext.Logger().Error(common.CodeOf(err), err.Error())
			cliContext.Spinner().Stop()
		}

	}

	err = common.WriteServiceResponseData(response.ServiceName, *response, cliContext)
	if err != nil {
		cliContext.Spinner().Stop()
		cliContext.Logger().SetErrorToStderr()
		cliContext.Logger().Error(common.CodeOf(err), err.Error())

	}

	cliContext.Spinner().StopWithMessage("Icon Node Started. Please find service details in current working directory(services.json)")
}

func iconDecentralization(cmd *cobra.Command, args []string) {

	cliContext := common.GetCliWithKurtosisContext()

	err := common.ValidateArgs(args)

	if err != nil {
		cliContext.Logger().Fatal(common.CodeOf(err), err.Error())
	}

	cliContext.Spinner().StartWithMessage("Starting Icon Node Decentralization", "green")

	params := GetDecentralizeParams(serviceName, nodeEndpoint, ksPath, ksPassword, networkID)

	err = RunDecentralization(cliContext, params)

	if err != nil {
		cliContext.Logger().Error(common.KurtosisContextError, err.Error())

	}

	cliContext.Spinner().StopWithMessage("Decentralization Completed")
}
