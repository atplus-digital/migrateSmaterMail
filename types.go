package main

type InputCredencialsFileDTO struct {
	Users             []UsersSctruct       `json:"users"`
	ServerAddress     SourceAddressDTO     `json:"source"`
	SmarterMailConfig SmarterMailConfigDTO `json:"smartermail"`
}

type SmarterMailConfigDTO struct {
	Host              string `json:"host"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Domain            string `json:"domain"`
	IsPasswordExpired bool   `json:"passwordExpired"`
}

type SourceAddressDTO struct {
	Address string `json:"address"`
	Domain  string `json:"domain"`
	TLS     bool   `json:"tls/ssl"`
}

type UsersSctruct struct {
	Username      string `json:"email"`
	TargetAccount string `json:"targetAccount"`
	Password      string `json:"password"`
}

type EmailAuthResult struct {
	Email     string
	AuthError error
}

type EmailCreateResult struct {
	Email       string
	CreateError error
}

type EmailMigrateResult struct {
	Email string
	Error error
}

type InMailAccount struct {
	Email         string
	TargetAccount string
	Password      string
	Domain        string
}

type CreateUserInputDTO struct {
	UserData `json:"userData"`
}

type UserData struct {
	UserName          string `json:"username"`
	FullName          string `json:"fullName"`
	Password          string `json:"password"`
	IsPasswordExpired bool   `json:"isPasswordExpired"`
	SecurityFlags     `json:"securityFlags"`
}

type SecurityFlags struct {
	AuthType                    int  `json:"authType"`
	AuthenticatingWindowsDomain any  `json:"authenticatingWindowsDomain"`
	IsDomainAdmin               bool `json:"isDomainAdmin"`
}

type MigrateMailboxStruct struct {
	ImapAccount `json:"imapAccount"`
}

type ImapAccount struct {
	ServerAddress                string `json:"serverAddress"`
	Username                     string `json:"username"`
	Password                     string `json:"password"`
	ServerPort                   int    `json:"serverPort"`
	UseSsl                       bool   `json:"useSsl"`
	EnableSpamFilter             bool   `json:"enableSpamFilter"`
	IsManualRetrieval            bool   `json:"isManualRetrieval"`
	AccountType                  string `json:"accountType"`
	UseOnlyOnce                  bool   `json:"useOnlyOnce"`
	UserDisplayed                bool   `json:"userDisplayed"`
	AccountTypeDescription       string `json:"accountTypeDescription"`
	ItemsToImport                int    `json:"itemsToImport"`
	IsMailboxMigration           bool   `json:"isMailboxMigration"`
	DeleteEverythingBeforeImport bool   `json:"deleteEverythingBeforeImport"`
}

type ExpireUsersPasswordDTO struct {
	EmailAddresses []string `json:"input"`
}
