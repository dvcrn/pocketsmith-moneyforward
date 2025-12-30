package mfmatch

import (
	"fmt"
	"strings"

	"github.com/dvcrn/pocketsmith-go"
	"github.com/dvcrn/pocketsmith-moneyforward/internal/accountmatch"
)

const MoneyForwardSuffix = " (MF)"

func normalizeAccountTitle(title string) string {
	return strings.ToLower(strings.Join(strings.Fields(title), " "))
}

func baseInstitutionName(institutionName string) string {
	return strings.TrimSuffix(strings.TrimSpace(institutionName), MoneyForwardSuffix)
}

func BuildBaseName(subName, subType string) string {
	name := strings.TrimSpace(subName)
	kind := strings.TrimSpace(subType)

	switch {
	case name != "" && kind != "":
		return fmt.Sprintf("%s (%s)", name, kind)
	case name != "":
		return name
	case kind != "":
		return kind
	default:
		return ""
	}
}

func BuildDisplayAccountName(institutionName, baseName string) string {
	displayInstitution := baseInstitutionName(institutionName)
	baseName = strings.TrimSpace(baseName)
	if baseName != "" && strings.Contains(displayInstitution, baseName) {
		baseName = ""
	}

	if baseName == "" {
		if displayInstitution == "" {
			return ""
		}
		return displayInstitution + MoneyForwardSuffix
	}

	displayName := accountmatch.BuildDisplayAccountName(displayInstitution, baseName)
	if strings.HasSuffix(displayName, MoneyForwardSuffix) {
		return displayName
	}

	return displayName + MoneyForwardSuffix
}

func LegacyAccountCandidates(institutionName, subName, subType string, includeInstitution bool) []string {
	base := baseInstitutionName(institutionName)
	raw := fmt.Sprintf("%s - %s (%s)", base, subName, subType)

	trimmedSubName := strings.TrimSpace(subName)
	trimmedSubType := strings.TrimSpace(subType)

	cleaned := raw
	switch {
	case trimmedSubName == "" && trimmedSubType == "":
		cleaned = base
	case trimmedSubName == "":
		cleaned = fmt.Sprintf("%s - %s", base, trimmedSubType)
	case trimmedSubType == "":
		cleaned = fmt.Sprintf("%s - %s", base, trimmedSubName)
	default:
		cleaned = fmt.Sprintf("%s - %s (%s)", base, trimmedSubName, trimmedSubType)
	}

	candidates := []string{raw, cleaned}
	if includeInstitution && base != "" {
		candidates = append(candidates, base, base+MoneyForwardSuffix)
	}

	return uniqueNonEmpty(candidates)
}

func FindMatchingAccount(accounts []*pocketsmith.Account, institutionName, baseName, displayName string, legacyCandidates []string) (*pocketsmith.Account, error) {
	candidates := filterAccountsByInstitution(accounts, institutionName)
	if len(candidates) == 0 {
		return nil, pocketsmith.ErrNotFound
	}

	exactCandidates := append([]string{displayName}, legacyCandidates...)
	for _, candidate := range exactCandidates {
		account, err := findExactMatch(candidates, candidate)
		if err != nil {
			return nil, err
		}
		if account != nil {
			return account, nil
		}
	}

	return nil, pocketsmith.ErrNotFound
}

func filterAccountsByInstitution(accounts []*pocketsmith.Account, institutionName string) []*pocketsmith.Account {
	normalizedInstitution := normalizeAccountTitle(institutionName)
	if normalizedInstitution == "" {
		return nil
	}

	var filtered []*pocketsmith.Account
	for _, account := range accounts {
		accountInstitution := normalizeAccountTitle(account.PrimaryTransactionAccount.Institution.Title)
		if accountInstitution == normalizedInstitution {
			filtered = append(filtered, account)
		}
	}

	return filtered
}

func findExactMatch(accounts []*pocketsmith.Account, candidate string) (*pocketsmith.Account, error) {
	normalizedCandidate := normalizeAccountTitle(candidate)
	if normalizedCandidate == "" {
		return nil, nil
	}

	var matches []*pocketsmith.Account
	for _, account := range accounts {
		if normalizeAccountTitle(account.Title) == normalizedCandidate {
			matches = append(matches, account)
		}
	}

	if len(matches) == 1 {
		return matches[0], nil
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple Pocketsmith accounts match %q", candidate)
	}

	return nil, nil
}

func uniqueNonEmpty(items []string) []string {
	seen := make(map[string]bool, len(items))
	unique := make([]string, 0, len(items))

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if seen[item] {
			continue
		}
		seen[item] = true
		unique = append(unique, item)
	}

	return unique
}
