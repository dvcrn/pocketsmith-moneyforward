package main

import (
	"flag"
	"fmt"
	"github.com/dvcrn/pocketsmith-moneyforward/sanitizer"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	moneyforward "github.com/dvcrn/moneyforward-go"
	"github.com/dvcrn/pocketsmith-go"
)

type Config struct {
	MoneyForwardCookie string
	PocketsmithToken   string
}

func getConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.MoneyForwardCookie, "mf-cookie", os.Getenv("MONEYFORWARD_COOKIE"), "MoneyForward cookie string")
	flag.StringVar(&config.PocketsmithToken, "pocketsmith-token", os.Getenv("POCKETSMITH_TOKEN"), "Pocketsmith API token")
	flag.Parse()

	if config.MoneyForwardCookie == "" {
		fmt.Println("Error: MoneyForward cookie is required. Set via -mf-cookie flag or MONEYFORWARD_COOKIE environment variable")
		os.Exit(1)
	}
	if config.PocketsmithToken == "" {
		fmt.Println("Error: Pocketsmith token is required. Set via -pocketsmith-token flag or POCKETSMITH_TOKEN environment variable")
		os.Exit(1)
	}

	return config
}

func findOrCreateAccount(ps *pocketsmith.Client, userID int, institutionName string, accountName string, accountType moneyforward.AccountType, currency string) (*pocketsmith.Account, error) {
	account, err := ps.FindAccountByName(userID, accountName)
	if err != nil {
		if err != pocketsmith.ErrNotFound {
			return nil, err
		}

		institution, err := ps.FindInstitutionByName(userID, institutionName)
		if err != nil {
			if err != pocketsmith.ErrNotFound {
				return nil, err
			}

			institution, err = ps.CreateInstitution(userID, institutionName, strings.ToLower(currency))
			if err != nil {
				return nil, err
			}
		}

		var psAccountType pocketsmith.AccountType
		switch accountType {
		case moneyforward.AccountTypeBank:
			psAccountType = pocketsmith.AccountTypeBank
		case moneyforward.AccountTypeCard:
			psAccountType = pocketsmith.AccountTypeCredits
		case moneyforward.AccountTypeSecurities:
			psAccountType = pocketsmith.AccountTypeStocks
		case moneyforward.AccountTypeCryptoFXPreciousMetals:
			psAccountType = pocketsmith.AccountTypeStocks
		case moneyforward.AccountTypeElectronicMoneyPrepaid:
			psAccountType = pocketsmith.AccountTypeBank
		case moneyforward.AccountTypePoints:
			psAccountType = pocketsmith.AccountTypeOtherAsset
		case moneyforward.AccountTypeMobilePhone:
			psAccountType = pocketsmith.AccountTypeOtherAsset
		case moneyforward.AccountTypeOnlineShopping:
			psAccountType = pocketsmith.AccountTypeOtherAsset
		case moneyforward.AccountTypeSupermarket:
			psAccountType = pocketsmith.AccountTypeOtherAsset
		default:
			psAccountType = pocketsmith.AccountTypeOtherAsset
		}

		fmt.Println("Would create account: ", institutionName, accountName)
		account, err := ps.CreateAccount(userID, institution.ID, accountName, "jpy", psAccountType)
		if err != nil {
			return nil, err
		}

		return account, nil
	}

	return account, nil
}

func main() {
	config := getConfig()

	// Initialize clients
	ps := pocketsmith.NewClient(config.PocketsmithToken)
	mf := moneyforward.NewClient(config.MoneyForwardCookie)

	// kick off mf sync
	if err := mf.ForceUpdate(); err != nil {
		fmt.Println("Error syncing MF: ", err)
	}

	// Get current Pocketsmith user
	currentUser, err := ps.GetCurrentUser()
	if err != nil {
		panic(err)
	}

	// Get MoneyForward accounts
	accounts, err := mf.GetAccountSummaries()
	if err != nil {
		panic(err)
	}

	// Process each account
	for _, account := range accounts.Accounts {
		fmt.Printf("Processing MoneyForward account: %s\n", account.Name)

		if account.Type != moneyforward.AccountTypeCard &&
			account.Type != moneyforward.AccountTypeBank &&
			//account.Type != moneyforward.AccountTypeSecurities &&
			//account.Type != moneyforward.AccountTypeCryptoFXPreciousMetals &&
			account.Type != moneyforward.AccountTypeElectronicMoneyPrepaid {
			fmt.Printf("Skipping account: %s\n", account.Name)
			continue
		}

		// Create or find corresponding Pocketsmith account
		fmt.Println("Finding or creating Pocketsmith account...")
		// psAccount, err := findOrCreateAccount(ps, currentUser.ID, account.ServiceType, account.Name, account.Type)
		// if err != nil {
		// 	fmt.Printf("Error creating/finding account: %v\n", err)
		// 	continue
		// }

		for _, subAccount := range account.SubAccounts {
			accDetail, err := mf.GetSubAccountDetail(account.AccountIDHash, subAccount.SubAccountIDHash)
			if err != nil {
				log.Fatal("Failed to get account detail:", err)
			}

			fmt.Println("Processing sub account: ", subAccount.SubName, "(", subAccount.SubType)

			fmt.Printf("\nAccount Detail: %+v\n", accDetail)

			keys := make([]string, 0, len(accDetail.AccountDetail.UserAssetDets))
			for k := range accDetail.AccountDetail.UserAssetDets {
				keys = append(keys, k)
			}

			// Get earliest login date from account detail
			earliestDate := accDetail.AccountDetail.UserAssetDets[keys[0]][0].Account.Account.FirstSucceededAt
			currency := accDetail.AccountDetail.UserAssetDets[keys[0]][0].Currency
			now := time.Now()

			institutionName := fmt.Sprintf("%s (MF)", account.Name)
			accountName := fmt.Sprintf("%s - %s (%s)", account.Name, subAccount.SubName, subAccount.SubType)
			psAccount, err := findOrCreateAccount(ps, currentUser.ID, institutionName, accountName, account.Type, currency)
			if err != nil {
				fmt.Printf("Error creating/finding account: %v\n", err)
				return
			}

			// set the initial balance on earliestDate to 0
			ps.UpdateTransactionAccount(psAccount.PrimaryTransactionAccount.ID, psAccount.PrimaryTransactionAccount.Institution.ID, 0, earliestDate.AddDate(0, 0, -1).Format("2006-01-02"))

			var allActs []*moneyforward.UserAssetAct
			currentDate := earliestDate

			// Iterate through months
			for currentDate.Before(now) {
				endDate := currentDate.AddDate(0, 10, 0)
				if endDate.After(now) {
					endDate = now
				}

				// Format dates as required by the API (YYYY-MM-DD)
				fromStr := currentDate.Format("2006-01-02")
				toStr := endDate.Format("2006-01-02")

				fmt.Printf("Fetching data for period %s to %s...\n", fromStr, toStr)

				// Get cash flow data for the current month
				flowData, err := mf.GetCashFlowTermData(subAccount.SubAccountIDHash, fromStr, toStr)
				if err != nil {
					log.Printf("Failed to get cash flow data for period %s to %s: %v", fromStr, toStr, err)
					return
				}

				// Extract and append transactions
				for _, act := range flowData.UserAssetActs {
					allActs = append(allActs, &act.UserAssetAct)
				}

				fmt.Println("Fetched data for period", fromStr, "to", toStr, "with", len(flowData.UserAssetActs), "transactions. Total so far:", len(allActs))

				// Move to next month
				currentDate = endDate
			}

			// Get transactions for this account

			repeatedFoundTransactions := 0

			// sort allActs by recognized date

			// remove duplicates by ID
			allActsMap := make(map[string]bool)
			filteredAllActs := []*moneyforward.UserAssetAct{}

			for _, act := range allActs {
				if _, exists := allActsMap[act.ID.String()]; !exists {
					allActsMap[act.ID.String()] = true
					filteredAllActs = append(filteredAllActs, act)
				}
			}

			allActs = filteredAllActs

			// Sort allActs by RecognizedAt, with newest coming first
			sort.Slice(allActs, func(i, j int) bool {
				// TODO: change me back to newest first
				// temporarily inverted it to backfill from the other side
				return allActs[i].RecognizedAt.Before(allActs[j].RecognizedAt)
			})

			for i, tx := range allActs {
				if repeatedFoundTransactions > 15 {
					fmt.Println("Too many repeated transactions found, likely everything processed already. Skipping...")
					break
				}

				convertedPayee := sanitizer.Sanitize(tx.Content)

				fmt.Printf("[%d/%d] Processing transaction: %s %s\n", i, len(allActs), tx.ID, convertedPayee)

				// Create Pocketsmith transaction
				mfidMemo := fmt.Sprintf("mfid=%s", tx.ID.String())
				psTx := &pocketsmith.CreateTransaction{
					Payee:        convertedPayee,
					Amount:       tx.Amount,
					Date:         tx.RecognizedAt.Format("2006-01-02"),
					IsTransfer:   tx.IsTransfer,
					NeedsReview:  false,
					Memo:         mfidMemo,
					ChequeNumber: tx.ID.String(),
				}

				// Check if transaction already exists
				searchRes, err := ps.SearchTransactionsByMemoContains(psAccount.PrimaryTransactionAccount.ID, tx.RecognizedAt, tx.ID.String())
				if err != nil {
					fmt.Printf("Error searching transactions: %v\n", err)
					continue
				}

				if len(searchRes) > 0 {
					fmt.Printf("Found existing transaction: %s\n", searchRes[0].Payee)
					repeatedFoundTransactions++
					continue
				}

				// Add new transaction
				_, err = ps.AddTransaction(psAccount.PrimaryTransactionAccount.ID, psTx)
				if err != nil {
					fmt.Printf("Error adding transaction: %v\n", err)
					continue
				}
			}
		}
	}
}
