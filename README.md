# any2anexoj

![Screenshot](https://i.ibb.co/TDDypq8X/Screenshot-2025-11-18-at-16-17-16.png)

This tool converts the statements from known brokers and exchanges into a format compatible with section 9 from the Portuguese IRS form: [Mod_3_anexo_j](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/modelos_formularios/irs/Documents/Mod_3_anexo_J.pdf)

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
cat statement.csv | any2anexoj-cli --platform=tranding212
```
