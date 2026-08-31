package seed

import (
	"time"

	"github.com/kauanpecanha/odsquiz-quiz/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type questionSeed struct {
	ID        string
	Number    int
	Text      string
	ODSNumber int
	Options   []optionSeed
}

type optionSeed struct {
	ID          string
	OptionValue int
	Label       string
	Score       int
}

func Questions(db *gorm.DB) error {
	now := time.Now()

	return db.Transaction(func(tx *gorm.DB) error {
		for _, item := range individualQuizQuestions {
			question := models.QuizQuestion{
				ID:        item.ID,
				Number:    item.Number,
				Text:      item.Text,
				ODSNumber: item.ODSNumber,
				Active:    true,
				CreatedAt: now,
				UpdatedAt: now,
			}

			err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"number",
					"text",
					"ods_number",
					"active",
					"updated_at",
				}),
			}).Create(&question).Error
			if err != nil {
				return err
			}

			for _, optionItem := range item.Options {
				option := models.QuizAnswerOption{
					ID:          optionItem.ID,
					QuestionID:  item.ID,
					OptionValue: optionItem.OptionValue,
					Label:       optionItem.Label,
					Score:       optionItem.Score,
					CreatedAt:   now,
					UpdatedAt:   now,
				}

				err := tx.Clauses(clause.OnConflict{
					Columns: []clause.Column{{Name: "id"}},
					DoUpdates: clause.AssignmentColumns([]string{
						"question_id",
						"option_value",
						"label",
						"score",
						"updated_at",
					}),
				}).Create(&option).Error
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

var individualQuizQuestions = []questionSeed{
	{
		ID:        "00000000-0000-4000-8000-000000000001",
		Number:    1,
		ODSNumber: 1,
		Text:      "Quantos salarios minimos compoem sua renda familiar bruta, contando com seus beneficios (exemplo: auxilios recebidos do governo)?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000001", OptionValue: 1, Label: "Ate 1 salario minimo.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000002", OptionValue: 2, Label: "Entre 1 e 2 salarios minimos.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000003", OptionValue: 3, Label: "Entre 3 e 5 salarios minimos.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000004", OptionValue: 4, Label: "Superior a 5 salarios minimos.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000005", OptionValue: 5, Label: "Nao gostaria de responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000002",
		Number:    2,
		ODSNumber: 3,
		Text:      "Atualmente, voce considera ter acesso a servicos publicos de saude de qualidade?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000006", OptionValue: 1, Label: "Nao ha servico publico de saude perto de minha residencia.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000007", OptionValue: 2, Label: "Ha servicos publicos de saude basicos, mas em condicoes extremamente precarias de preservacao e atendimento.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000008", OptionValue: 3, Label: "Ha servicos publicos de saude basicos, mas ainda ha a necessidade de me deslocar para outra regiao para melhores condicoes deste servico.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000009", OptionValue: 4, Label: "Conto com todos os tipos de servico publico de saude, de postos de saude a hospitais, em condicoes adequadas.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000010", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000003",
		Number:    3,
		ODSNumber: 4,
		Text:      "Defina a qualidade das instituicoes educacionais na comunidade na qual voce mora.",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000011", OptionValue: 1, Label: "Pessimo.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000012", OptionValue: 2, Label: "Regular.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000013", OptionValue: 3, Label: "Bom.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000014", OptionValue: 4, Label: "Muito bom.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000015", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000004",
		Number:    4,
		ODSNumber: 12,
		Text:      "Voce pratica a coleta seletiva dos residuos gerados no seu domicilio. Qual o grau de reciclagem que voce realiza?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000016", OptionValue: 1, Label: "Nao faco coleta seletiva, tenho lixo comum em casa.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000017", OptionValue: 2, Label: "Separo apenas alguns itens de reciclagem, quando da tempo.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000018", OptionValue: 3, Label: "Separo e encaminho para reciclagem e logistica reversa rotineiramente.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000019", OptionValue: 4, Label: "Pratico diariamente a separacao de residuos e encaminho adequadamente cada um deles.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000020", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000005",
		Number:    5,
		ODSNumber: 7,
		Text:      "Com que frequencia voce ve paineis solares ou carros eletricos na comunidade em que vive?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000021", OptionValue: 1, Label: "Nunca os vejo.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000022", OptionValue: 2, Label: "Pouquissimos.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000023", OptionValue: 3, Label: "Frequentemente os vejo.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000024", OptionValue: 4, Label: "Toda a comunidade em que vivo e energizada por energia solar e movida a carros eletricos.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000025", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000006",
		Number:    6,
		ODSNumber: 12,
		Text:      "Voce se considera uma pessoa mais consumista ou consciente?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000026", OptionValue: 1, Label: "Totalmente consumista, costumo comprar com frequencia sem pensar no impacto.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000027", OptionValue: 2, Label: "Parcialmente consumista, faco compras quase sempre que almejo algo.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000028", OptionValue: 3, Label: "Parcialmente consciente, sem deixar totalmente de comprar quando tenho vontade.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000029", OptionValue: 4, Label: "Totalmente consciente, compro somente quando ha necessidade e avalio o impacto ambiental.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000030", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000007",
		Number:    7,
		ODSNumber: 8,
		Text:      "Na regiao onde voce mora, como e caracterizado o comercio local?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000031", OptionValue: 1, Label: "Comumente devemos nos deslocar para outros locais ou bairros para encontrar oportunidades.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000032", OptionValue: 2, Label: "Ha comercio local composto por trabalho informal e comercios de pequeno porte.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000033", OptionValue: 3, Label: "Ha startups, microempresas, transnacionais e franquias famosas.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000034", OptionValue: 4, Label: "Comercio forte, com empresas relevantes, iniciativas empreendedoras e empregos formais.", Score: 400},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000008",
		Number:    8,
		ODSNumber: 5,
		Text:      "Quantas iniciativas de incentivo a igualdade de genero e insercao da mulher no mercado de trabalho voce conhece?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000035", OptionValue: 1, Label: "Nao conheco nenhuma iniciativa deste genero.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000036", OptionValue: 2, Label: "Entre uma e tres.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000037", OptionValue: 3, Label: "Entre tres e dez.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000038", OptionValue: 4, Label: "De dez para mais.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000039", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000009",
		Number:    9,
		ODSNumber: 6,
		Text:      "A comunidade em que vive possui acesso a agua potavel?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000040", OptionValue: 1, Label: "Nao ha acesso a agua potavel, seja por abastecimento municipal ou poco.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000041", OptionValue: 2, Label: "Poucos tem este acesso.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000042", OptionValue: 3, Label: "Uma parte significativa tem este acesso.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000043", OptionValue: 4, Label: "Toda a minha comunidade tem acesso a agua potavel.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000044", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000010",
		Number:    10,
		ODSNumber: 6,
		Text:      "A comunidade em que vive possui acesso a saneamento basico?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000045", OptionValue: 1, Label: "Nao ha rede de esgoto ou solucoes individuais adequadas.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000046", OptionValue: 2, Label: "Poucos tem acesso a saneamento basico.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000047", OptionValue: 3, Label: "Uma parte significativa tem acesso ao saneamento basico.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000048", OptionValue: 4, Label: "Toda a minha comunidade tem acesso a saneamento basico de qualidade.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000049", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000011",
		Number:    11,
		ODSNumber: 2,
		Text:      "Na comunidade em que voce vive, voce conhece pessoas que passam fome? E do seu conhecimento projetos locais voltados a ajuda-las?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000050", OptionValue: 1, Label: "Conheco pessoas em situacao de fome, e nao ha projetos voltados a ajuda-las.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000051", OptionValue: 2, Label: "Conheco pessoas em situacao de fome, entretanto ha um ou mais projetos voltados para sua ajuda.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000052", OptionValue: 3, Label: "Nao conheco pessoas em situacao de fome, e nao ha projetos voltados para elas onde moro.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000053", OptionValue: 4, Label: "Nao conheco pessoas em situacao de fome, e conheco projetos voltados para ajuda destas pessoas.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000054", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000012",
		Number:    12,
		ODSNumber: 14,
		Text:      "Como voce considera o nivel de qualidade dos ambientes aquaticos (mares e lagoas) proximos a sua comunidade?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000055", OptionValue: 1, Label: "Estao sempre sujos e poluidos.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000056", OptionValue: 2, Label: "Constantemente estao sujos e com lixo.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000057", OptionValue: 3, Label: "As vezes estao limpos.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000058", OptionValue: 4, Label: "Estao sempre limpos e despoluidos.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000059", OptionValue: 5, Label: "Nao tenho ambientes aquaticos proximos.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000013",
		Number:    13,
		ODSNumber: 14,
		Text:      "Voce conhece areas de preservacao marinha na sua regiao?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000060", OptionValue: 1, Label: "Nao conheco.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000061", OptionValue: 2, Label: "Conheco ao menos 1.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000062", OptionValue: 3, Label: "Conheco ao menos 3.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000063", OptionValue: 4, Label: "Conheco mais de 5.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000064", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000014",
		Number:    14,
		ODSNumber: 13,
		Text:      "Voce conhece areas de preservacao terrestre na sua regiao?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000065", OptionValue: 1, Label: "Nao conheco.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000066", OptionValue: 2, Label: "Conheco ao menos 1.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000067", OptionValue: 3, Label: "Conheco ao menos 3.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000068", OptionValue: 4, Label: "Conheco mais de 5.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000069", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000015",
		Number:    15,
		ODSNumber: 15,
		Text:      "As areas de reflorestamento sao essenciais para a preservacao da fauna e flora e recuperar areas degradadas. Voce conhece projetos de reflorestamento na sua comunidade?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000070", OptionValue: 1, Label: "Nao conheco nenhuma iniciativa deste genero.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000071", OptionValue: 2, Label: "Entre uma e tres.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000072", OptionValue: 3, Label: "Entre tres e cinco.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000073", OptionValue: 4, Label: "De cinco para mais.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000074", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000016",
		Number:    16,
		ODSNumber: 13,
		Text:      "Em relacao as catastrofes naturais, a sua comunidade passou por algum desses problemas nos ultimos meses?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000075", OptionValue: 1, Label: "Nao temos nenhuma catastrofe natural nos ultimos 3 anos.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000076", OptionValue: 2, Label: "Tivemos alguns momentos de catastrofe nos ultimos 2 anos.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000077", OptionValue: 3, Label: "Temos frequentemente sofrido com catastrofes anualmente.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000078", OptionValue: 4, Label: "Temos sofrido com algum tipo de catastrofe natural mais de uma vez no ano.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000079", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000017",
		Number:    17,
		ODSNumber: 13,
		Text:      "Algumas atividades humanas contribuem para a emissao de gases de efeito estufa. Assinale a opcao com que voce melhor se identifica.",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000080", OptionValue: 1, Label: "Nao mudo minhas atividades em funcao da emissao de gases de efeito estufa.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000081", OptionValue: 2, Label: "Raramente faco alguma mudanca de habito em funcao de GEE.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000082", OptionValue: 3, Label: "As vezes faco algumas alteracoes de habito para reduzir as emissoes de GEE.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000083", OptionValue: 4, Label: "Sempre busco opcoes na minha rotina para reduzir as emissoes de GEE.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000084", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000018",
		Number:    18,
		ODSNumber: 11,
		Text:      "Na sua comunidade existem pessoas vivendo em assentamentos ou habitacoes precarios, inadequados ou informais?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000085", OptionValue: 1, Label: "Nao conheco pessoas em situacao de habitacao precaria.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000086", OptionValue: 2, Label: "Conheco algumas pessoas em situacao de habitacao precaria.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000087", OptionValue: 3, Label: "Conheco muitas pessoas em situacao de habitacao precaria.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000088", OptionValue: 4, Label: "Todas as pessoas na minha comunidade vivem em situacao de habitacao precaria.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000089", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000019",
		Number:    19,
		ODSNumber: 11,
		Text:      "Existem na sua comunidade iniciativas de reciclagem e/ou compostagem?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000090", OptionValue: 1, Label: "Nao existem iniciativas ou nao conheco.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000091", OptionValue: 2, Label: "Existem algumas poucas iniciativas, mas nao participo.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000092", OptionValue: 3, Label: "Existem algumas iniciativas, e as vezes participo.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000093", OptionValue: 4, Label: "Existem varias iniciativas e participo com frequencia.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000094", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000020",
		Number:    20,
		ODSNumber: 10,
		Text:      "Voce conhece iniciativas que promovam aumento de renda para pessoas vivendo abaixo de 50% da mediana da renda na sua comunidade?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000095", OptionValue: 1, Label: "Nao conheco.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000096", OptionValue: 2, Label: "Conheco apenas uma.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000097", OptionValue: 3, Label: "Conheco entre uma e tres.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000098", OptionValue: 4, Label: "Conheco mais de tres.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000099", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000021",
		Number:    21,
		ODSNumber: 10,
		Text:      "Voce ou alguem que voce conhece ja foi discriminado em funcao de genero, idade, condicao fisica, etnia ou cidade/pais de origem?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000100", OptionValue: 1, Label: "Nunca sofri discriminacao de qualquer tipo.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000101", OptionValue: 2, Label: "Raramente sofro algum tipo de discriminacao.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000102", OptionValue: 3, Label: "As vezes sofro algum tipo de discriminacao.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000103", OptionValue: 4, Label: "Todos os dias sofro algum tipo de discriminacao.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000104", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
	{
		ID:        "00000000-0000-4000-8000-000000000022",
		Number:    22,
		ODSNumber: 16,
		Text:      "Voce ou alguem que voce conhece ja foi vitima de violencia fisica, psicologica ou sexual?",
		Options: []optionSeed{
			{ID: "10000000-0000-4000-8000-000000000105", OptionValue: 1, Label: "Nunca.", Score: 100},
			{ID: "10000000-0000-4000-8000-000000000106", OptionValue: 2, Label: "Raramente.", Score: 200},
			{ID: "10000000-0000-4000-8000-000000000107", OptionValue: 3, Label: "As vezes.", Score: 300},
			{ID: "10000000-0000-4000-8000-000000000108", OptionValue: 4, Label: "Todos os dias.", Score: 400},
			{ID: "10000000-0000-4000-8000-000000000109", OptionValue: 5, Label: "Nao sei responder.", Score: 0},
		},
	},
}
