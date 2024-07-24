

Mova o arquivo settings-example.json para settings.json com o seguinte comando:

``
  mv settings-example.json settings.json
``

Configure o arquivo settings.json de acordo com sua necessidade

### Objeto `source`:  Dados do servidores de email de origem. 

 - O campo `host` deve ser preenchido com o endereço do servidor de email (Dominio ou IP). *

 - O campo `domain` é o dominio que esta sendo migrado.  *

### Objeto `smartermail`: Dados para a conexão com o SmarterMail.
 - O campo `host` deve ser preenchido com o endereço do servidor SmarterMail. *
 - O campo `username` deve ser preenchido com o usuário admin do SmarterMail. *
 - O campo `password` deve ser preenchido com a senha do usuário admin do  SmarterMail. *
 - O campo `domain` deve ser preenchido com o dominio de destino. *


### Array `users`: Dados das contas para serem migradas
 - O campo `email` deve ser preenchido com a conta de email (Apenas a conta e não o dominio junto) *
 - O campo `password` deve ser preenchido com a senha da conta de email *
 - O campo `targetAccount` pode ser preenchido com a conta de destino (Opcional)
    - Default: Valor do campo email

PS: A senha deve ser a mesma que a senha do conta de email na origem, a mesma senha também sera utilizada na conta no SmarterMail

## Executando o script:

Com o arquivo settings.json configurado, baixe o script no mesmo caminho/pasta que o arquivo settings.json

Execute o script em seu terminal:

### Bash:
`chmod +x ./migrateSmarterMail`

`./migrateSmarterMail` 

### Powershell/CMD:
`.\migrateSmarterMail-windows.exe`







