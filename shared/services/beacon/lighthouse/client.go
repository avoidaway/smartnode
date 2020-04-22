package lighthouse

import (
    "bytes"
    "encoding/hex"
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "strconv"

    "github.com/rocket-pool/smartnode/shared/services/beacon"
    hexutil "github.com/rocket-pool/smartnode/shared/utils/hex"
)


// Beacon config
const REQUEST_CONTENT_TYPE string = "application/json"

// Beacon endpoints
const REQUEST_ETH2_CONFIG_PATH string = "/spec"
const REQUEST_SLOTS_PER_EPOCH_PATH string = "/spec/slots_per_epoch"
const REQUEST_BEACON_HEAD_PATH string = "/beacon/head"
const REQUEST_VALIDATORS_PATH string = "/beacon/validators"


// Beacon request types
type ValidatorsRequest struct {
    Pubkeys []string                `json:"pubkeys"`
}

// Beacon response types
type Eth2ConfigResponse struct {
    GenesisForkVersion string       `json:"genesis_fork_version"`
    BLSWithdrawalPrefixByte string  `json:"bls_withdrawal_prefix_byte"`
    DomainBeaconProposer uint64     `json:"domain_beacon_proposer"`
    DomainBeaconAttester uint64     `json:"domain_beacon_attester"`
    DomainRandao uint64             `json:"domain_randao"`
    DomainDeposit uint64            `json:"domain_deposit"`
    DomainVoluntaryExit uint64      `json:"domain_voluntary_exit"`
}
type BeaconHeadResponse struct {
    Slot uint64                     `json:"slot"`
    FinalizedSlot uint64            `json:"finalized_slot"`
    JustifiedSlot uint64            `json:"justified_slot"`
}
type ValidatorResponse struct {
    Pubkey string                   `json:"pubkey"`
    Validator struct {
        WithdrawalCredentials string        `json:"withdrawal_credentials"`
        EffectiveBalance uint64             `json:"effective_balance"`
        Slashed bool                        `json:"slashed"`
        ActivationEligibilityEpoch uint64   `json:"activation_eligibility_epoch"`
        ActivationEpoch uint64              `json:"activation_epoch"`
        ExitEpoch uint64                    `json:"exit_epoch"`
        WithdrawableEpoch uint64            `json:"withdrawable_epoch"`
    }                               `json:"validator"`
}


// Client
type Client struct {
    providerUrl string
}


/**
 * Create client
 */
func NewClient(providerUrl string) *Client {
    return &Client{
        providerUrl: providerUrl,
    }
}


/**
 * Get the eth2 config
 */
func (c *Client) GetEth2Config() (*beacon.Eth2Config, error) {

    // Data channels
    configChannel := make(chan Eth2ConfigResponse)
    slotsPerEpochChannel := make(chan uint64)
    errorChannel := make(chan error)

    // Request eth2 config
    go (func() {
        var config Eth2ConfigResponse
        if responseBody, err := c.getRequest(REQUEST_ETH2_CONFIG_PATH); err != nil {
            errorChannel <- errors.New("Error retrieving eth2 config: " + err.Error())
        } else if err := json.Unmarshal(responseBody, &config); err != nil {
            errorChannel <- errors.New("Error unpacking eth2 config: " + err.Error())
        } else {
            configChannel <- config
        }
    })()

    // Request slots per epoch
    go (func() {
        if responseBody, err := c.getRequest(REQUEST_SLOTS_PER_EPOCH_PATH); err != nil {
            errorChannel <- errors.New("Error retrieving slots per epoch: " + err.Error())
        } else if slotsPerEpoch, err := strconv.Atoi(string(responseBody)); err != nil {
            errorChannel <- errors.New("Error unpacking slots per epoch: " + err.Error())
        } else {
            slotsPerEpochChannel <- uint64(slotsPerEpoch)
        }
    })()

    // Receive data
    var config Eth2ConfigResponse
    var slotsPerEpoch uint64
    for received := 0; received < 2; {
        select {
            case config = <-configChannel:
                received++
            case slotsPerEpoch = <-slotsPerEpochChannel:
                received++
            case err := <-errorChannel:
                return nil, err
        }
    }

    // Create response
    response := &beacon.Eth2Config{
        DomainBeaconProposer: config.DomainBeaconProposer,
        DomainBeaconAttester: config.DomainBeaconAttester,
        DomainRandao: config.DomainRandao,
        DomainDeposit: config.DomainDeposit,
        DomainVoluntaryExit: config.DomainVoluntaryExit,
        SlotsPerEpoch: slotsPerEpoch,
    }

    // Decode hex data and update
    if genesisForkVersionBytes, err := hex.DecodeString(hexutil.RemovePrefix(config.GenesisForkVersion)); err != nil {
        return nil, errors.New("Error decoding genesis fork version: " + err.Error())
    } else {
        response.GenesisForkVersion = genesisForkVersionBytes
    }
    if blsWithdrawalPrefixByteBytes, err := hex.DecodeString(hexutil.RemovePrefix(config.BLSWithdrawalPrefixByte)); err != nil {
        return nil, errors.New("Error decoding BLS withdrawal prefix byte: " + err.Error())
    } else {
        response.BLSWithdrawalPrefixByte = blsWithdrawalPrefixByteBytes
    }

    // Return
    return response, nil

}


/**
 * Get the beacon head
 */
func (c *Client) GetBeaconHead() (*beacon.BeaconHead, error) {

    // Data channels
    headChannel := make(chan BeaconHeadResponse)
    slotsPerEpochChannel := make(chan uint64)
    errorChannel := make(chan error)

    // Request beacon head
    go (func() {
        var head BeaconHeadResponse
        if responseBody, err := c.getRequest(REQUEST_BEACON_HEAD_PATH); err != nil {
            errorChannel <- errors.New("Error retrieving beacon head: " + err.Error())
        } else if err := json.Unmarshal(responseBody, &head); err != nil {
            errorChannel <- errors.New("Error unpacking beacon head: " + err.Error())
        } else {
            headChannel <- head
        }
    })()

    // Request slots per epoch
    go (func() {
        if responseBody, err := c.getRequest(REQUEST_SLOTS_PER_EPOCH_PATH); err != nil {
            errorChannel <- errors.New("Error retrieving slots per epoch: " + err.Error())
        } else if slotsPerEpoch, err := strconv.Atoi(string(responseBody)); err != nil {
            errorChannel <- errors.New("Error unpacking slots per epoch: " + err.Error())
        } else {
            slotsPerEpochChannel <- uint64(slotsPerEpoch)
        }
    })()

    // Receive data
    var head BeaconHeadResponse
    var slotsPerEpoch uint64
    for received := 0; received < 2; {
        select {
            case head = <-headChannel:
                received++
            case slotsPerEpoch = <-slotsPerEpochChannel:
                received++
            case err := <-errorChannel:
                return nil, err
        }
    }

    // Return response
    return &beacon.BeaconHead{
        Epoch: head.Slot / slotsPerEpoch,
        FinalizedEpoch: head.FinalizedSlot / slotsPerEpoch,
        JustifiedEpoch: head.JustifiedSlot / slotsPerEpoch,
    }, nil

}


/**
 * Get a validator's status
 */
func (c *Client) GetValidatorStatus(pubkey string) (*beacon.ValidatorStatus, error) {

    // Request
    responseBody, err := c.postRequest(REQUEST_VALIDATORS_PATH, ValidatorsRequest{Pubkeys: []string{pubkey}})
    if err != nil {
        return nil, errors.New("Error retrieving validator status: " + err.Error())
    }

    // Unmarshal response
    var validators []ValidatorResponse
    if err := json.Unmarshal(responseBody, &validators); err != nil {
        return nil, errors.New("Error unpacking validator status: " + err.Error())
    }
    validator := validators[0]

    // Create response
    response := &beacon.ValidatorStatus{
        EffectiveBalance: validator.Validator.EffectiveBalance,
        Slashed: validator.Validator.Slashed,
        ActivationEligibilityEpoch: validator.Validator.ActivationEligibilityEpoch,
        ActivationEpoch: validator.Validator.ActivationEpoch,
        ExitEpoch: validator.Validator.ExitEpoch,
        WithdrawableEpoch: validator.Validator.WithdrawableEpoch,
        Exists: (validator.Validator.ActivationEpoch != 0), // Activation epoch is 0 only if validator is null in JSON response
    }

    // Decode hex data and update
    if pubkeyBytes, err := hex.DecodeString(hexutil.RemovePrefix(validator.Pubkey)); err != nil {
        return nil, errors.New("Error decoding validator pubkey: " + err.Error())
    } else {
        response.Pubkey = pubkeyBytes
    }
    if withdrawalCredentialsBytes, err := hex.DecodeString(hexutil.RemovePrefix(validator.Validator.WithdrawalCredentials)); err != nil {
        return nil, errors.New("Error decoding validator withdrawal credentials: " + err.Error())
    } else {
        response.WithdrawalCredentials = withdrawalCredentialsBytes
    }

    // Return
    return response, nil

}


/**
 * Make GET request to beacon server
 */
func (c *Client) getRequest(requestPath string) ([]byte, error) {

    // Send request
    response, err := http.Get(c.providerUrl + requestPath)
    if err != nil { return nil, err }
    defer response.Body.Close()

    // Get response
    body, err := ioutil.ReadAll(response.Body)
    if err != nil { return nil, err }

    // Return
    return body, nil

}


/**
 * Make POST request to beacon server
 */
func (c *Client) postRequest(requestPath string, requestBody interface{}) ([]byte, error) {

    // Get request body
    requestBodyBytes, err := json.Marshal(requestBody)
    if err != nil { return nil, err }
    requestBodyReader := bytes.NewReader(requestBodyBytes)

    // Send request
    response, err := http.Post(c.providerUrl + requestPath, REQUEST_CONTENT_TYPE, requestBodyReader)
    if err != nil { return nil, err }
    defer response.Body.Close()

    // Get response
    body, err := ioutil.ReadAll(response.Body)
    if err != nil { return nil, err }

    // Return
    return body, nil

}

