package config

type SeedConfig struct {
	AdminFirstName string
	AdminLastName  string
	AdminUsername  string
	AdminEmail     string
	AdminPhone     string
	AdminPassword  string
}

func LoadSeedConfig() (SeedConfig, error) {
	adminFirstName, err := extractEnv("SEED_ADMIN_FIRST_NAME")
	if err != nil {
		return SeedConfig{}, err
	}

	adminLastName, err := extractEnv("SEED_ADMIN_LAST_NAME")
	if err != nil {
		return SeedConfig{}, err
	}

	adminUsername, err := extractEnv("SEED_ADMIN_USERNAME")
	if err != nil {
		return SeedConfig{}, err
	}

	adminEmail, err := extractEnv("SEED_ADMIN_EMAIL")
	if err != nil {
		return SeedConfig{}, err
	}

	adminPhone, err := extractEnv("SEED_ADMIN_PHONE")
	if err != nil {
		return SeedConfig{}, err
	}

	adminPassword, err := extractEnv("SEED_ADMIN_PASSWORD")
	if err != nil {
		return SeedConfig{}, err
	}

	return SeedConfig{
		AdminFirstName: adminFirstName,
		AdminLastName:  adminLastName,
		AdminUsername:  adminUsername,
		AdminEmail:     adminEmail,
		AdminPhone:     adminPhone,
		AdminPassword:  adminPassword,
	}, nil
}
