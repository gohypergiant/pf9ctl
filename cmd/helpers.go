package cmd

import (
	"errors"
	"fmt"

	"github.com/platform9/pf9ctl/pkg/color"
	"github.com/platform9/pf9ctl/pkg/keystone"
	"github.com/platform9/pf9ctl/pkg/pmk"
	"github.com/platform9/pf9ctl/pkg/util"
	"go.uber.org/zap"
)

// This function will validate the user credentials entered during config set and bail out if invalid
func validateUserCredentials(pmk.Config, pmk.Client) error {

	auth, err := c.Keystone.GetAuth(
		ctx.Username,
		ctx.Password,
		ctx.Tenant,
	)

	if err != nil {
		RegionInvalid = false
		return err
	}

	// To validate region.
	endpointURL, err1 := pmk.FetchRegionFQDN(ctx, auth)
	if endpointURL == "" || err1 != nil {
		RegionInvalid = true
		zap.S().Debug("Invalid Region")
		return errors.New("Invalid Region")
	}

	return nil
}

func configValidation(bool, int) error {

	if pmk.LoopCounter <= MaxLoopNoConfig-1 {

		//Check if we are setting config through pf9ctl config set command.
		if !SetConfig {
			// If Oldconfig exists and invalid credentials entered during config prompt
			if pmk.OldConfigExist {
				if pmk.InvalidExistingConfig {
					// If user enters invalid credentials during prompt of config (due to invalid config found after config loading).
					if RegionInvalid {
						fmt.Println("\n" + color.Red("x ") + "Invalid Region entered")
						zap.S().Debug("Invalid Region entered")
					} else {
						fmt.Println("\n" + color.Red("x ") + "Invalid credentials entered (Platform9 Account URL/Username/Password/Region/Tenant)")
						zap.S().Debug("Invalid credentials entered (Platform9 Account URL/Username/Password/Region/Tenant)")
					}

				} else if pmk.OldConfigExist && pmk.LoopCounter == 0 {
					// If invalid credentials are found during config loading
					if RegionInvalid {
						fmt.Println("\n" + color.Red("x ") + "Invalid Region found")
						zap.S().Debug("Invalid Region found")
					} else {
						fmt.Println("\n" + color.Red("x ") + "Invalid credentials found (Platform9 Account URL/Username/Password/Region/Tenant)")
						zap.S().Debug("Invalid credentials found (Platform9 Account URL/Username/Password/Region/Tenant)")
					}
				}

			} else {
				// If user enters invalid credentials during new config promput (due to config not found)
				if RegionInvalid {
					fmt.Println("\n" + color.Red("x ") + "Invalid Region entered")
					zap.S().Debug("Invalid Region entered")
				} else {
					fmt.Println("\n" + color.Red("x ") + "Invalid credentials entered (Platform9 Account URL/Username/Password/Region/Tenant)")
					zap.S().Debug("Invalid credentials entered (Platform9 Account URL/Username/Password/Region/Tenant)")
				}

			}
		} else {
			// If user enters invalid credentials during config set through "pf9ctl config set"
			if RegionInvalid {
				fmt.Println("\n" + color.Red("x ") + "Invalid Region entered")
				zap.S().Debug("Invalid Region entered")
			} else {
				fmt.Println("\n" + color.Red("x ") + "Invalid credentials entered (Platform9 Account URL/Username/Password/Region/Tenant)")
				zap.S().Debug("Invalid credentials entered (Platform9 Account URL/Username/Password/Region/Tenant)")
			}

		}
	}
	// If existing initial config is Invalid
	if (pmk.LoopCounter == 0) && (pmk.OldConfigExist) {
		pmk.InvalidExistingConfig = true
		pmk.LoopCounter += 1
	} else {
		// If user enteres invalid credentials during new config pormpt.
		pmk.LoopCounter += 1
	}

	// If any invalid credentials extered multiple times in new config prompt then to bail out the recursive loop (thrice)
	if pmk.LoopCounter >= MaxLoopNoConfig && !(pmk.InvalidExistingConfig) {
		zap.S().Fatalf("Invalid credentials entered multiple times (Platform9 Account URL/Username/Password/Region/Tenant)")
	} else if pmk.LoopCounter >= MaxLoopNoConfig+1 && pmk.InvalidExistingConfig {
		if RegionInvalid {
			fmt.Println(color.Red("x ") + "Invalid Region entered")
		} else {
			fmt.Println(color.Red("x ") + "Invalid credentials entered (Platform9 Account URL/Username/Password/Region/Tenant)")
		}
		zap.S().Fatalf("Invalid credentials entered multiple times (Platform9 Account URL/Username/Password/Region/Tenant)")
	}
	return nil
}

func loadCredentials() (pmk.Config, pmk.Client) {
	var (
		ctx pmk.Config
		err error
	)
	// This flag is used to loop back if user enters invalid credentials during config set.
	credentialFlag = true
	// To bail out if loop runs recursively more than thrice
	pmk.LoopCounter = 0

	for credentialFlag {

		ctx, err = pmk.LoadConfig(util.Pf9DBLoc)
		if err != nil {
			zap.S().Fatalf("Unable to load the context: %s\n", err.Error())
		}

		executor, err := getExecutor()
		if err != nil {
			zap.S().Debug("Error connecting to host %s", err.Error())
			zap.S().Fatalf(" Invalid (Username/Password/IP)")
		}

		c, err = pmk.NewClient(ctx.Fqdn, executor, ctx.AllowInsecure, false)
		if err != nil {
			zap.S().Fatalf("Unable to load clients needed for the Cmd. Error: %s", err.Error())
		}

		// Validate the user credentials entered during config set and will loop back again if invalid
		if err := validateUserCredentials(ctx, c); err != nil {
			clearContext(&pmk.Context)
			//Check if no or invalid config exists, then bail out if asked for correct config for maxLoop times.
			err = configValidation(RegionInvalid, pmk.LoopCounter)
		} else {
			// We will store the set config if its set for first time using check-node
			if pmk.IsNewConfig {
				if err := pmk.StoreConfig(ctx, util.Pf9DBLoc); err != nil {
					zap.S().Errorf("Failed to store config: %s", err.Error())
				} else {
					pmk.IsNewConfig = false
				}
			}
			credentialFlag = false
		}
	}
	return ctx, c
}

func getAuthConfig(c pmk.Client) (keystone.KeystoneAuth, error) {
	config, _ := loadCredentials()
	auth, err := c.Keystone.GetAuth(config.Username, config.Password, config.Tenant)
	if err != nil {
		zap.S().Debug("Failed to get keystone %s", err.Error())
	}
	return auth, err
}
