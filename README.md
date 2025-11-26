# any2anexoj

[![Go Report Card](https://goreportcard.com/badge/github.com/nmoniz/any2anexoj)](https://goreportcard.com/report/github.com/nmoniz/any2anexoj)
[![Coverage Status](https://coveralls.io/repos/github/nmoniz/any2anexoj/badge.svg?branch=coveralls-badge)](https://coveralls.io/github/nmoniz/any2anexoj?branch=coveralls-badge)

<p align="center">
  <img src="https://i.ibb.co/0yRtwq2C/0-FBA40-FD-D97-A-4-AFB-8618-49582-DB98-F3-C.png" alt="Screenshot" with="600" border="0">
</p>

This tool converts the statements from known brokers and exchanges into a format compatible with section 9 from the Portuguese IRS form: [Mod_3_anexo_j](https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/modelos_formularios/irs/Documents/Mod_3_anexo_J.pdf)

> [!WARNING]
> Although I made significant efforts to ensure the correctness of the calculations you should verify any outputs produced by this tool on your own or with a certified accountant.

## Install

```bash
go install github.com/nmoniz/any2anexoj/cmd/any2anexoj-cli@latest
```

## Usage

```bash
cat statement.csv | any2anexoj-cli --platform=tranding212
```

## Rounding

All Euro values are rounded to cents (2 decimal places) but internal calculations use the statement values with full precision.
There are no explicit rules or details about how to round Euro values in Anexo J.
This application rounds according to `Portaria n.º 1180/2001, art. 2.º, alínea c) e d)` (Ministerial Order / Government Order) examples, which imply we should round to the 2nd decimal place by rounding up (ceiling) or down (floor) depending on whether the third decimal place is ≥ 5 or < 5, respectively.
