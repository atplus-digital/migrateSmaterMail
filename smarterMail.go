package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

func (c *SmarterMailClient) CreateAccountsSmarterMail(WorkerPool int, InMailAccountChannel chan InMailAccount, resultChannel chan EmailCreateResult) {
	var wg sync.WaitGroup

	for i := 1; i <= WorkerPool; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for v := range InMailAccountChannel {
				err := c.CreateUserSmarterMail(v)
				if err != nil {
					// fmt.Printf("O Worker %v esta trabalhando com a conta %v\n", x, v.Email )
					// fmt.Println("Criando conta", v.Email)
					resultChannel <- EmailCreateResult{Email: v.Email, CreateError: err}
					return
				}
				resultChannel <- EmailCreateResult{Email: v.Email, CreateError: nil}
			}
		}(i)
	}
	wg.Wait()
	close(resultChannel)
}
func (c *SmarterMailClient) MigrateAccountsSmarterMail(WorkerPool int, InMailAccountChannel chan InMailAccount, resultChannel chan EmailMigrateResult, ServerAddress string) {
	var wg sync.WaitGroup

	for i := 1; i <= WorkerPool; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for v := range InMailAccountChannel {
				fmt.Println("Migrando a conta", v.Email)
				err := c.MigrateMailboxSmarterMail(v, ServerAddress)
				if err != nil {
					resultChannel <- EmailMigrateResult{Email: v.Email, Error: err}
					return
				}
				resultChannel <- EmailMigrateResult{Email: v.Email, Error: nil}
			}
		}(i)
	}
	wg.Wait()
	close(resultChannel)
}

func (c *SmarterMailClient) CreateUserSmarterMail(u InMailAccount) error {
	UserExist, err := c.CheckUserExist(fmt.Sprintf("%v@%v", u.Email, u.Domain))
	if err != nil {
		return err
	}

	if UserExist {
		return nil
	}

	err = c.CreateUser(u.Email, u.Password, u.Domain)
	if err != nil {
		return err
	}
	return nil
}
func (c *SmarterMailClient) MigrateMailboxSmarterMail(u InMailAccount, ServerAddress string) error {

	username := fmt.Sprintf("%v@%v", u.Email, u.Domain)

	ServerAddressSplit := strings.Split(ServerAddress, ":")
	ServerAddressName := ServerAddressSplit[0]
	ServerAddressPort, err := strconv.Atoi(ServerAddressSplit[1])
	if err != nil {
		ServerAddressPort = 993
	}

	UserSmarterMailConfig := SmarterMailConfigDTO{
		Host:     c.SmarterMailConfig.Host,
		Username: username,
		Password: u.Password,
	}

	CommonUserSm, err := InitSmarterMail(UserSmarterMailConfig)
	if err != nil {
		return err
	}

	MigrateMailboxInputDTO := MigrateMailboxStruct{
		ImapAccount{
			ServerAddress:                ServerAddressName,
			Username:                     username,
			Password:                     u.Password,
			ServerPort:                   ServerAddressPort,
			UseSsl:                       true,
			EnableSpamFilter:             false,
			IsManualRetrieval:            true,
			AccountType:                  "IMAP",
			UseOnlyOnce:                  true,
			UserDisplayed:                false,
			AccountTypeDescription:       "Other",
			ItemsToImport:                1,
			IsMailboxMigration:           true,
			DeleteEverythingBeforeImport: true,
		},
	}

	MigrateMailboxJson, err := json.Marshal(MigrateMailboxInputDTO)
	if err != nil {
		return err
	}

	MigrateMailboxBuf := bytes.NewBuffer(MigrateMailboxJson)

	MigrateMailboxResp, err := CommonUserSm.Post("/settings/imap-migration", MigrateMailboxBuf, nil)
	if err != nil {
		return err
	}
	defer MigrateMailboxResp.Body.Close()

	MigrateMailboxOutputDTO := struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}{}

	MigrateMailboxRespBytes, err := io.ReadAll(MigrateMailboxResp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(MigrateMailboxRespBytes, &MigrateMailboxOutputDTO)
	if err != nil {
		return err
	}

	if !MigrateMailboxOutputDTO.Success {
		return fmt.Errorf(MigrateMailboxOutputDTO.Message)
	}

	return nil
}

func (c *SmarterMailClient) CheckUserExist(mailAccount string) (bool, error) {
	MailAccountInputDTO := struct {
		Email string `json:"email"`
	}{
		Email: mailAccount,
	}

	MailAccountJsonPayload, err := json.Marshal(MailAccountInputDTO)
	if err != nil {
		return false, err
	}

	MailAccountJsonPayloadBuf := bytes.NewBuffer(MailAccountJsonPayload)

	resp, err := c.Post("/settings/sysadmin/user-exists", MailAccountJsonPayloadBuf, nil)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	ResponseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	RespMailAccountExist := struct {
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(ResponseBodyBytes, &RespMailAccountExist)
	if err != nil {
		return false, err
	}

	if RespMailAccountExist.Message == "ERRORS.USER_NOT_FOUND" {
		return false, nil
	}

	return true, nil
}
func (c *SmarterMailClient) CreateUser(mailAccount string, password string, domain string) error {

	CreateUserInputDTO := CreateUserInputDTO{
		UserData{
			UserName:          mailAccount,
			FullName:          mailAccount,
			Password:          password,
			IsPasswordExpired: false,
			SecurityFlags: SecurityFlags{
				AuthType:                    0,
				AuthenticatingWindowsDomain: nil,
				IsDomainAdmin:               false,
			},
		},
	}

	CreateUserInputJson, err := json.Marshal(CreateUserInputDTO)
	if err != nil {
		return err
	}

	CreateUserInputBuf := bytes.NewBuffer(CreateUserInputJson)

	SmarterMailDomainHeader := map[string]string{"X-SmarterMailDomain": domain}

	ResponseCreateUser, err := c.Post("/settings/domain/user-put", CreateUserInputBuf, SmarterMailDomainHeader)
	if err != nil {
		return err
	}

	defer ResponseCreateUser.Body.Close()

	var CreateUserOutputDTO struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	CreateUserOutputBuf, err := io.ReadAll(ResponseCreateUser.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(CreateUserOutputBuf, &CreateUserOutputDTO)
	if err != nil {
		return err
	}

	if !CreateUserOutputDTO.Success {
		return fmt.Errorf(CreateUserOutputDTO.Message)
	}

	return nil
}
