package track

import (
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// SaveDomain saves a new domain to the domains_tb table.
func (repo *Repository) SaveDomain(domain, key string) (int64, error) {
	result, err := repo.db.Exec("INSERT INTO domains_tb (domain, key) VALUES (?, ?)", domain, key)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetDomain returns the ID of the domain from the domains_tb table.
func (repo *Repository) GetDomain(domain string) (int, error) {
	var id int
	err := repo.db.QueryRow("SELECT id FROM domains_tb WHERE domain = ?", domain).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

type DomainKeyPair struct {
	Domain  string `db:"domain"`
	SiteKey string `db:"siteKey"`
}

// GetDomainKeyPair returns the key of the domain from the domains_tb table.
func (repo *Repository) GetDomainKeyPair(domain string) (DomainKeyPair, error) {
	var keyPair DomainKeyPair
	err := repo.db.QueryRow("SELECT domain, siteKey FROM domains_tb WHERE domain = ?", domain).Scan(&keyPair.Domain, &keyPair.SiteKey)
	if errors.Is(err, sql.ErrNoRows) {
		return DomainKeyPair{}, sql.ErrNoRows
	} else if err != nil {
		return DomainKeyPair{}, err
	}

	return keyPair, nil
}

// GetPage returns the ID of the page from the pages_tb table.
func (repo *Repository) GetPage(domainID int, pageURL string) (int, error) {
	var id int
	err := repo.db.QueryRow("SELECT id FROM pages_tb WHERE domain_id = ? AND page_url = ?", domainID, pageURL).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		} else {
			return 0, err
		}
	}

	return id, nil
}

// SavePage saves a new page to the pages_tb table.
func (repo *Repository) CreatePage(domainID int, pageURL string) (int64, error) {
	result, err := repo.db.Exec("INSERT INTO pages_tb (domain_id, page_url) VALUES (?, ?)", domainID, pageURL)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SaveIPAddress saves a new IP address to the ip_addresses_tb table.
func (repo *Repository) SaveIPAddress(ipAddress string) (int64, error) {
	result, err := repo.db.Exec("INSERT INTO ip_addresses_tb (ip_address) VALUES (?)", ipAddress)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SaveUTM saves a new UTM req to the utm_tb table.
func (repo *Repository) SaveUTM(pageID int, utmSource, utmMedium, utmCampaign string) (int64, error) {
	result, err := repo.db.Exec("INSERT INTO utm_tb (page_id, utm_source, utm_medium, utm_campaign) VALUES (?, ?, ?, ?)",
		pageID, utmSource, utmMedium, utmCampaign)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SaveClick saves a new click to the clicks_tb table.
func (repo *Repository) SaveClick(pageID, ipAddressID int, element, clickedURL string) (int64, error) {
	result, err := repo.db.Exec("INSERT INTO clicks_tb (page_id, ip_address_id, element, clicked_url) VALUES (?, ?, ?, ?)",
		pageID, ipAddressID, element, clickedURL)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
