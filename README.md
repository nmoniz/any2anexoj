# any2anexoj

This tool converts the statements from brokers and exchanges into a format compatible with the Portuguese IRS form: [Mod_3_anexo_j](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/modelos_formularios/irs/Documents/Mod_3_anexo_J.pdf)

> [!WARNING]
> Although I made significant efforts to ensure the correctness of the calculations you should verify any outputs produced by this tool on your own or with a certified accountant.

> [!NOTE]
> This tool is in early stages of development. Use at your own risk!

## Install

```bash
go install github.com/nmoniz/any2anexoj/cmd/any2anexoj-cli@latest
```

## Usage

```bash
cat statement.csv | broker2anexoj --platform=tranding212
```
