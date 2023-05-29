# customers
Cadastro dos nossos clientes

# How to run lint locally

> golangci-lint run ./...

# Servidor em modo seguro com TLS

Com o objetivo de sustentar boas práticas no desenvolvimento de apps cloud natives, o nosso servidor `Customers` exigirá TLS.

Assim, será necessário configurar um certificado de serviço para o seu funcionamento, e para o ambiente de desenvolvimento, iremos usar um certificado auto-assinado.

A geração de certificados auto-assinados pode ser feita seguindo o artigo https://medium.com/opentelemetry/securing-your-opentelemetry-collector-1a4f9fa5bd6f

Como o artigo está em inglês e queremos fornecer um material que os brasileiros podem usar, o resumo para as instruções é o seguinte:

- crie uma pasta /tmp
- instale a ferramenta open-source [cfssl](https://github.com/cloudflare/cfssl) que nos permitirá a criação dos certificados
- crie um arquivo `ca-csr.json` necessário para a autoridade certificadora (CA). Use como exemplo o trecho abaixo
```
{
    "hosts": [
        "localhost",
        "127.0.0.1"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "O":  "Customers (Projeto Dose na Nuvem)"
        }
    ]
}
```
- execute o comando:
> cfssl genkey -initca ca-csr.json | cfssljson -bare ca
- crie um arquivo `cert-csr.json` necessário para o certificado do serviço. Use como exemplo o trecho abaixo
```
{
    "hosts": [
        "localhost",
        "127.0.0.1"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "O":  "Customers (Projeto Dose na Nuvem)"
        }
    ]
}
```
- execute o comando:
> cfssl gencert -ca ca.pem -ca-key ca-key.pem cert-csr.json | cfssljson -bare cert

Os arquivos importantes para execução no modo seguro são:
- cert.pem
- cert-key.pem

Use esses certificados para execução do servidor

1 - via arquivo de configuração. Acrescente o seguinte trecho no arquivo `config.yaml`.
```
tls:
  cert_file: tmp/cert.pem
  cert_key_file: tmp/cert-key.pem
  insecure: false
```
2 - via linha de comando
> go run . start --config config.yaml --server.tls.certfile tmp/cert.pem --server.tls.certkeyfile tmp/cert-key.pem

# Project's Toolkit

- https://golangci-lint.run