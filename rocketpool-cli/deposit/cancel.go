package deposit

import (
    "errors"
    "fmt"

    "github.com/urfave/cli"

    "github.com/rocket-pool/smartnode/shared/services"
    "github.com/rocket-pool/smartnode/shared/utils/eth"
)


// Cancel a node deposit reservation
func cancelDeposit(c *cli.Context) error {

    // Initialise services
    p, err := services.NewProvider(c, services.ProviderOpts{
        AM: true,
        NodeContract: true,
        LoadContracts: []string{"rocketNodeAPI"},
        LoadAbis: []string{"rocketNodeContract"},
        WaitClientSync: true,
        WaitRocketStorage: true,
    })
    if err != nil { return err }
    defer p.Cleanup()

    // Check node has current deposit reservation
    hasReservation := new(bool)
    if err := p.NodeContract.Call(nil, hasReservation, "getHasDepositReservation"); err != nil {
        return errors.New("Error retrieving deposit reservation status: " + err.Error())
    } else if !*hasReservation {
        fmt.Fprintln(p.Output, "Node does not have a current deposit reservation")
        return nil
    }

    // Cancel deposit reservation
    if txor, err := p.AM.GetNodeAccountTransactor(); err != nil {
        return err
    } else {
        fmt.Fprintln(p.Output, "Canceling deposit reservation...")
        if _, err := eth.ExecuteContractTransaction(p.Client, txor, p.NodeContractAddress, p.CM.Abis["rocketNodeContract"], "depositReserveCancel"); err != nil {
            return errors.New("Error canceling deposit reservation: " + err.Error())
        }
    }

    // Log & return
    fmt.Fprintln(p.Output, "Deposit reservation cancelled successfully")
    return nil

}

