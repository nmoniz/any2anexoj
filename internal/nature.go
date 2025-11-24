package internal

type Nature string

const (
	// NatureUnknown is the zero value of Nature type
	NatureUnknown Nature = ""

	// NatureG01 describes selling of stocks per table VII: Alienação onerosa de ações/partes sociais
	NatureG01 Nature = "G01"

	// NatureG20 describes selling units in investment funds (including ETFs) as per table VII:
	// Resgates ou alienação de unidades de participação ou liquidação de fundos de investimento
	NatureG20 Nature = "G20"
)
