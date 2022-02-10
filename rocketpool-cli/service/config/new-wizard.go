package config

import (
	"fmt"

	"github.com/rocket-pool/smartnode/shared"
	"github.com/rocket-pool/smartnode/shared/services/config"
)

type newUserWizard struct {
	md                     *mainDisplay
	welcomeModal           *page
	networkModal           *choiceModalLayout
	executionModeModal     *choiceModalLayout
	executionLocalModal    *choiceModalLayout
	executionExternalModal *textBoxModalLayout
	infuraModal            *textBoxModalLayout
	fallbackInfuraModal    *textBoxModalLayout
	fallbackExecutionModal *choiceModalLayout
	consensusModeModal     *choiceModalLayout
	consensusLocalModal    *choiceModalLayout
	consensusExternalModal *choiceModalLayout
	finishedModal          *page
}

func newNewUserWizard(md *mainDisplay) *newUserWizard {

	wiz := &newUserWizard{
		md: md,
	}

	wiz.createWelcomeModal()
	wiz.createNetworkModal()
	wiz.createExecutionModeModal()
	wiz.createLocalExecutionModal()
	wiz.createExternalExecutionModal()
	wiz.createInfuraModal()
	wiz.createFallbackExecutionModal()
	wiz.createFallbackInfuraModal()
	wiz.createConsensusModeModal()
	wiz.createLocalConsensusModal()
	//wiz.createExternalConsensusModal()
	wiz.createFinishedModal()

	return wiz

}

// ========================
// === 1: Welcome Modal ===
// ========================
func (wiz *newUserWizard) createWelcomeModal() {

	modal := newChoiceModalLayout(
		wiz.md.app,
		60,
		shared.Logo+"\n\n"+

			"Welcome to the Smartnode configuration wizard!\n\n"+
			"Since this is your first time configuring the Smartnode, we'll walk you through the basic setup.\n\n",
		[]string{"Next", "Quit"}, nil, DirectionalModalHorizontal,
	)
	modal.done = func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 0 {
			wiz.md.setPage(wiz.networkModal.page)
			wiz.networkModal.focus(0)
		} else if buttonIndex == 1 {
			wiz.md.app.Stop()
		}
	}

	page := newPage(nil, "new-user-welcome", "New User Wizard > [1/8] Welcome", "", modal.borderGrid)
	wiz.md.pages.AddPage(page.id, page.content, true, false)

	wiz.welcomeModal = page

}

// =========================
// === 2: Select Network ===
// =========================
func (wiz *newUserWizard) createNetworkModal() {

	// Create the button names and descriptions from the config
	networks := wiz.md.config.Smartnode.Network.Options
	networkNames := []string{}
	networkDescriptions := []string{}
	for _, network := range networks {
		networkNames = append(networkNames, network.Name)
		networkDescriptions = append(networkDescriptions, network.Description)
	}

	// Create the modal
	modal := newChoiceModalLayout(
		wiz.md.app,
		70,
		"Let's start by choosing which network you'd like to use.\n\n",
		networkNames,
		networkDescriptions,
		DirectionalModalVertical)

	// Set up the callbacks
	modal.done = func(buttonIndex int, buttonLabel string) {
		newNetwork := networks[buttonIndex].Value.(config.Network)
		wiz.md.config.ChangeNetwork(newNetwork)
		wiz.md.setPage(wiz.executionModeModal.page)
		wiz.executionModeModal.focus(0)
	}

	// Create the page
	wiz.networkModal = modal
	page := newPage(nil, "new-user-network", "New User Wizard > [2/8] Network", "", modal.borderGrid)
	wiz.md.pages.AddPage(page.id, page.content, true, false)
	modal.page = page
}

// ================================
// === 3: Select Execution Mode ===
// ================================
func (wiz *newUserWizard) createExecutionModeModal() {

	// Create the button names and descriptions from the config
	modes := wiz.md.config.ExecutionClientMode.Options
	modeNames := []string{}
	modeDescriptions := []string{}
	for _, mode := range modes {
		modeNames = append(modeNames, mode.Name)
		modeDescriptions = append(modeDescriptions, mode.Description)
	}

	// Create the modal
	modal := newChoiceModalLayout(
		wiz.md.app,
		76,
		"Now let's decide which mode you'd like to use for the Execution client (formerly eth1 client).\n\n"+
			"Would you like Rocket Pool to run and manage its own client, or would you like it to use an existing client you run and manage outside of Rocket Pool (also known as \"Hybrid Mode\")?",

		modeNames,
		modeDescriptions,
		DirectionalModalVertical)

	// Set up the callbacks
	modal.done = func(buttonIndex int, buttonLabel string) {
		wiz.md.config.ExecutionClientMode.Value = modes[buttonIndex].Value
		switch modes[buttonIndex].Value {
		case config.Mode_Local:
			wiz.md.setPage(wiz.executionLocalModal.page)
			wiz.executionLocalModal.focus(0)
		case config.Mode_External:
			wiz.md.setPage(wiz.executionExternalModal.page)
		default:
			panic(fmt.Sprintf("Unknown execution client mode %s", modes[buttonIndex].Value))
		}
	}

	// Create the page
	wiz.executionModeModal = modal
	page := newPage(nil, "new-user-execution-mode", "New User Wizard > [3/8] Execution Client Mode", "", modal.borderGrid)
	wiz.md.pages.AddPage(page.id, page.content, true, false)
	modal.page = page

}

// ========================================
// === 4a: Select Local Exection Client ===
// ========================================
func (wiz *newUserWizard) createLocalExecutionModal() {

	// Create the button names and descriptions from the config
	clients := wiz.md.config.ExecutionClient.Options
	clientNames := []string{}
	clientDescriptions := []string{}
	for _, client := range clients {
		clientNames = append(clientNames, client.Name)
		clientDescriptions = append(clientDescriptions, client.Description)
	}

	// Create the modal
	modal := newChoiceModalLayout(
		wiz.md.app,
		76,
		"Please select the Execution client you would like to use.\n\n"+
			"Highlight each one to see a brief description of it, or go to https://docs.rocketpool.net/guides/node/eth-clients.html#eth1-clients to learn more about them.",
		clientNames,
		clientDescriptions,
		DirectionalModalVertical)

	// Set up the callbacks
	modal.done = func(buttonIndex int, buttonLabel string) {
		selectedClient := clients[buttonIndex].Value.(config.ExecutionClient)
		wiz.md.config.ExecutionClient.Value = selectedClient
		switch selectedClient {
		case config.ExecutionClient_Geth:
			// Geth doesn't have any required parameters so move on
			wiz.md.setPage(wiz.fallbackExecutionModal.page)
			wiz.fallbackExecutionModal.focus(0)
		case config.ExecutionClient_Infura:
			// Switch to the Infura dialog
			wiz.md.setPage(wiz.infuraModal.page)
			wiz.infuraModal.focus()
		case config.ExecutionClient_Pocket:
			// Pocket doesn't have any required parameters so move on
			wiz.md.setPage(wiz.fallbackExecutionModal.page)
			wiz.fallbackExecutionModal.focus(0)
		}
	}

	// Create the page
	wiz.executionLocalModal = modal
	page := newPage(nil, "new-user-execution-local", "New User Wizard > [4/8] Execution Client > Selection", "", modal.borderGrid)
	wiz.md.pages.AddPage(page.id, page.content, true, false)
	modal.page = page

}

// ===========================================
// === 4b: Select External Exection Client ===
// ===========================================
func (wiz *newUserWizard) createExternalExecutionModal() {

	// Create the labels
	httpLabel := wiz.md.config.ExternalExecution.HttpUrl.Name
	wsLabel := wiz.md.config.ExternalExecution.WsUrl.Name

	// Create the modal
	modal := newTextBoxModalLayout(
		wiz.md.app,
		70,
		"Please enter the URL of the HTTP-based RPC API and the URL of the Websocket-based RPC API for your existing client.\n\n"+
			"For example: `http://192.168.1.45:8545` and `ws://192.168.1.45:8546`",
		[]string{httpLabel, wsLabel},
		[]string{})

	// Set up the callbacks
	modal.done = func(text map[string]string) {
		wiz.md.config.ExternalExecution.HttpUrl.Value = text[httpLabel]
		wiz.md.config.ExternalExecution.WsUrl.Value = text[wsLabel]
		wiz.md.setPage(wiz.fallbackExecutionModal.page)
		wiz.fallbackExecutionModal.focus(0)
	}

	// Create the page
	wiz.executionExternalModal = modal
	page := newPage(nil, "new-user-execution-external", "New User Wizard > [4/8] Execution Client (External)", "", modal.borderGrid)
	wiz.md.pages.AddPage(page.id, page.content, true, false)
	modal.page = page

}

// ========================
// === 4c: Local Infura ===
// ========================
func (wiz *newUserWizard) createInfuraModal() {

	// Create the labels
	projectIdLabel := wiz.md.config.Infura.ProjectID.Name

	// Create the modal
	modal := newTextBoxModalLayout(
		wiz.md.app,
		70,
		"Please enter the Project ID for your Infura Ethereum project. You can find this on the Infura website, in your Ethereum project settings.",
		[]string{projectIdLabel},
		[]string{})

	// Set up the callbacks
	modal.done = func(text map[string]string) {
		wiz.md.config.Infura.ProjectID.Value = text[projectIdLabel]
		wiz.md.setPage(wiz.fallbackExecutionModal.page)
		wiz.fallbackExecutionModal.focus(0)
	}

	// Create the page
	wiz.infuraModal = modal
	page := newPage(nil, "new-user-execution-infura", "New User Wizard > [4/8] Execution Client > Infura", "", modal.borderGrid)
	wiz.md.pages.AddPage(page.id, page.content, true, false)
	modal.page = page

}

// =====================================
// === 5a: Fallback Execution Client ===
// =====================================
func (wiz *newUserWizard) createFallbackExecutionModal() {

	// Create the button names and descriptions from the config
	clients := wiz.md.config.FallbackExecutionClient.Options
	clientNames := []string{"None"}
	clientDescriptions := []string{"Do not use a fallback client."}
	for _, client := range clients {
		clientNames = append(clientNames, client.Name)
		clientDescriptions = append(clientDescriptions, client.Description)
	}

	// Create the modal
	modal := newChoiceModalLayout(
		wiz.md.app,
		70,
		"If you would like to add a fallback Execution client, please choose it below.\n\nThe Smartnode will temporarily use this instead of your main Execution client if the main client ever fails.\nIt will switch back to the main client when it starts working again.",
		clientNames,
		clientDescriptions,
		DirectionalModalVertical,
	)

	// Set up the callbacks
	modal.done = func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 0 {
			wiz.md.config.UseFallbackExecutionClient.Value = false
		} else {
			wiz.md.config.UseFallbackExecutionClient.Value = true
			selectedClient := clients[buttonIndex].Value.(config.ExecutionClient)
			wiz.md.config.ExecutionClient.Value = selectedClient
			switch selectedClient {
			case config.ExecutionClient_Infura:
				// Switch to the Infura dialog
				wiz.md.setPage(wiz.fallbackInfuraModal.page)
				wiz.fallbackInfuraModal.focus()
			}
		}
	}

	// Create the page
	wiz.fallbackExecutionModal = modal
	page := newPage(nil, "new-user-fallback-execution", "New User Wizard > [5/8] Fallback Execution Client", "", modal.borderGrid)
	wiz.md.pages.AddPage(page.id, page.content, true, false)
	modal.page = page

}

// ===========================
// === 5b: Fallback Infura ===
// ===========================
func (wiz *newUserWizard) createFallbackInfuraModal() {

	// Create the labels
	projectIdLabel := wiz.md.config.FallbackInfura.ProjectID.Name

	// Create the modal
	modal := newTextBoxModalLayout(
		wiz.md.app,
		70,
		"Please enter the Project ID for your Infura Ethereum project. You can find this on the Infura website, in your Ethereum project settings.",
		[]string{projectIdLabel},
		[]string{})

	// Set up the callbacks
	modal.done = func(text map[string]string) {
		wiz.md.config.FallbackInfura.ProjectID.Value = text[projectIdLabel]
		wiz.md.setPage(wiz.fallbackExecutionModal.page)
	}

	// Create the page
	wiz.fallbackInfuraModal = modal
	page := newPage(nil, "new-user-fallback-execution-infura", "New User Wizard > [5/8] Fallback Execution Client > Infura", "", modal.borderGrid)
	wiz.md.pages.AddPage(page.id, page.content, true, false)
	modal.page = page

}

// ================================
// === 6: Select Consensus Mode ===
// ================================
func (wiz *newUserWizard) createConsensusModeModal() {

	// Create the button names and descriptions from the config
	modes := wiz.md.config.ConsensusClientMode.Options
	modeNames := []string{}
	modeDescriptions := []string{}
	for _, mode := range modes {
		modeNames = append(modeNames, mode.Name)
		modeDescriptions = append(modeDescriptions, mode.Description)
	}

	// Create the modal
	modal := newChoiceModalLayout(
		wiz.md.app,
		76,
		"Next, let's decide which mode you'd like to use for the Consensus client (formerly eth2 client).\n\n"+
			"Would you like Rocket Pool to run and manage its own client, or would you like it to use an existing client you run and manage outside of Rocket Pool (also known as \"Hybrid Mode\")?",

		modeNames,
		modeDescriptions,
		DirectionalModalVertical)

	// Set up the callbacks
	modal.done = func(buttonIndex int, buttonLabel string) {
		wiz.md.config.ConsensusClientMode.Value = modes[buttonIndex].Value
		switch modes[buttonIndex].Value {
		case config.Mode_Local:
			wiz.md.setPage(wiz.consensusLocalModal.page)
			wiz.consensusLocalModal.focus(0)
		case config.Mode_External:
			wiz.md.setPage(wiz.consensusExternalModal.page)
		default:
			panic(fmt.Sprintf("Unknown execution client mode %s", modes[buttonIndex].Value))
		}
	}

	// Create the page
	wiz.consensusModeModal = modal
	page := newPage(nil, "new-user-consensus-mode", "New User Wizard > [6/8] Consensus Client Mode", "", modal.borderGrid)
	wiz.md.pages.AddPage(page.id, page.content, true, false)
	modal.page = page

}

// =========================================
// === 7a: Select Local Consensus Client ===
// =========================================
func (wiz *newUserWizard) createLocalConsensusModal() {

	// Create the button names and descriptions from the config
	clients := wiz.md.config.ConsensusClient.Options
	clientNames := []string{}
	clientDescriptions := []string{}
	for _, client := range clients {
		clientNames = append(clientNames, client.Name)
		clientDescriptions = append(clientDescriptions, client.Description)
	}

	// Create the modal
	modal := newChoiceModalLayout(
		wiz.md.app,
		76,
		"Please select the Consensus client you would like to use.\n\n"+
			"Highlight each one to see a brief description of it, or go to https://docs.rocketpool.net/guides/node/eth-clients.html#eth2-clients to learn more about them.",
		clientNames,
		clientDescriptions,
		DirectionalModalVertical)

	// Set up the callbacks
	modal.done = func(buttonIndex int, buttonLabel string) {
		selectedClient := clients[buttonIndex].Value.(config.ConsensusClient)
		wiz.md.config.ConsensusClient.Value = selectedClient
		switch selectedClient {
		/*
			case config.ExecutionClient_Geth:
				// Geth doesn't have any required parameters so move on
				wiz.md.setPage(wiz.fallbackExecutionModal.page)
			case config.ExecutionClient_Infura:
				// Switch to the Infura dialog
				wiz.md.setPage(wiz.infuraModal.page)
				wiz.infuraModal.focus()
			case config.ExecutionClient_Pocket:
				// Pocket doesn't have any required parameters so move on
				wiz.md.setPage(wiz.fallbackExecutionModal.page)
		*/
		}
	}

	// Create the page
	wiz.consensusLocalModal = modal
	page := newPage(nil, "new-user-consensus-local", "New User Wizard > [7/8] Consensus Client > Selection", "", modal.borderGrid)
	wiz.md.pages.AddPage(page.id, page.content, true, false)
	modal.page = page

}

// Create the finished modal
func (wiz *newUserWizard) createFinishedModal() {

	modal := newChoiceModalLayout(
		wiz.md.app,
		40,
		"All done! You're ready to run.\n\n"+
			"If you'd like, you can review and change all of the Smartnode and client settings next or just save and exit.",
		[]string{
			"Review All Settings",
			"Save and Exit",
		},
		nil,
		DirectionalModalVertical)
	modal.done = func(buttonIndex int, buttonLabel string) {
		if buttonIndex == 0 {
			wiz.md.setPage(wiz.md.settingsHome.homePage)
		} else {
			wiz.md.app.Stop()
		}
	}

	page := newPage(nil, "new-user-finished", "New User Wizard > [8/8] Finished", "", modal.borderGrid)
	wiz.md.pages.AddPage(page.id, page.content, true, false)

	wiz.finishedModal = page

}
