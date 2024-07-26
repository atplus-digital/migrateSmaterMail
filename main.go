package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
)

func main() {

	// Puxa as configurações do arquivo settings.json
	setting := initConfig("./settings.json")

	// Inicia uma instancia do SmarterMail
	sm, err := InitSmarterMail(setting.SmarterMailConfig)
	check(err)

	var numWokerPool int
	val, present := os.LookupEnv("NUM_WORKER_POOL")
	if present {
		workerPoolEnv, err := strconv.Atoi(val)
		if err != nil {
			// fmt.Println("Erro ao utilizar a variavel NUM_WORKER_POOL, utilizando o valor 5 para workerPool")
			numWokerPool = 5
		} else {
			numWokerPool = workerPoolEnv
		}

	} else {
		numWokerPool = 5
	}

	resultsTestMail := make(chan EmailAuthResult)
	var numEmailAuthError int

	color.White("Verificando autenticação das contas... \n\n")
	go testEmailAuthentication(setting.Users, setting.ServerAddress, resultsTestMail)

	for v := range resultsTestMail {
		if v.AuthError != nil {
			// fmt.Println(v.Email, "Error", v.AuthError)
			fmt.Printf("%v: %v - %v\n", color.RedString("Error"), v.Email, v.AuthError)
			numEmailAuthError++
			continue
		}
		fmt.Printf("%v: %v usuário autenticado com sucesso\n", color.GreenString("Success"), v.Email)
	}

	if numEmailAuthError > 0 {
		log.Fatalln(color.RedString("Número de autenticação com erro: %v", numEmailAuthError))
	}
	fmt.Println("")
	color.Green("Validação de autenticação - Success\n\n")

	// Cria as contas no SmarterMail - WorkerPool

	color.White("Criando contas no SmarterMail... \n\n")

	resultCreateAccount := make(chan EmailCreateResult)
	InMailAccountChannel := make(chan InMailAccount)
	var numCreateMailError int

	go func(InMailAccountChannel chan InMailAccount) {
		for _, v := range setting.Users {
			InMailAccountChannel <- InMailAccount{
				Email:         v.Username,
				TargetAccount: v.TargetAccount,
				Password:      v.Password,
				FullName:      v.FullName,
				JobTitle:      v.JobTitle,
				Domain:        setting.SmarterMailConfig.Domain}
		}
		close(InMailAccountChannel)

	}(InMailAccountChannel)

	go sm.CreateAccountsSmarterMail(numWokerPool, InMailAccountChannel, resultCreateAccount)

	for v := range resultCreateAccount {
		if v.CreateError != nil {
			fmt.Printf("%v: %v - %v\n", color.RedString("Error"), v.Email, v.CreateError)
			numCreateMailError++
			continue
		}
		fmt.Printf("%v: %v@%v usuário criado com sucesso\n", color.GreenString("Success"), v.Email, setting.SmarterMailConfig.Domain)
	}

	if numCreateMailError > 0 {
		log.Fatalln(color.RedString("Número de contas que deram erro ao serem criadas: %v", numCreateMailError))
	}
	fmt.Println("")
	color.Green("Criação de usuários - Success\n\n")

	// Inicia a Migração das contas - WorkerPool
	color.White("Iniciando a migração das contas...\n\n")

	InMigrateMailboxChannel := make(chan InMailAccount)
	resultMigrateMailboxChannel := make(chan EmailMigrateResult)
	var numMigrateError int

	go func(InMigrateMailboxChannel chan InMailAccount) {
		for _, v := range setting.Users {
			InMigrateMailboxChannel <- InMailAccount{
				Email:         v.Username,
				TargetAccount: v.TargetAccount,
				Password:      v.Password,
				Domain:        setting.SmarterMailConfig.Domain}
		}
		close(InMigrateMailboxChannel)

	}(InMigrateMailboxChannel)

	go sm.MigrateAccountsSmarterMail(numWokerPool, InMigrateMailboxChannel, resultMigrateMailboxChannel, setting.ServerAddress)

	for v := range resultMigrateMailboxChannel {
		if v.Error != nil {
			fmt.Printf("%v: %v - %v\n", color.RedString("Error"), v.Email, v.Error)
			numMigrateError++
			continue
		}
		fmt.Printf("%v: %v Tarefa de migração criada com sucesso\n", color.GreenString("Success"), v.Email)

	}

	if numMigrateError > 0 {
		log.Fatalln(color.RedString("Número de contas que deram erro ao serem migradas: %v", numMigrateError))
	}
	fmt.Println("")
	color.Green("Migração das contas - Success\n\n")

	// Expire as senhas das contas caso o campo IsPasswordExpired seja true
	if sm.SmarterMailConfig.IsPasswordExpired {
		color.White("Expirando as senhas...\n\n")
		err := sm.ExpireUsersPassword(setting.Users)
		check(err)

		fmt.Printf("%v: Senhas expiradas com sucesso\n\n", color.GreenString("Success"))

	}

}

func initConfig(pathJsonfile string) InputCredencialsFileDTO {
	readFile, err := os.ReadFile(pathJsonfile)
	if err != nil {
		log.Fatalf("Não foi possivel ler o arquivo: %v", pathJsonfile)
	}

	var InputCredencials InputCredencialsFileDTO

	err = json.Unmarshal(readFile, &InputCredencials)
	check(err)

	return InputCredencials
}

func check(e error) {
	if e != nil {
		ErrorTxt := fmt.Sprintf("%v: %v", color.RedString("Error:"), e)
		log.Fatal(ErrorTxt)
	}
}

func getfullEmail(email string, domain string) string {
	return fmt.Sprintf("%v@%v", email, domain)
}
