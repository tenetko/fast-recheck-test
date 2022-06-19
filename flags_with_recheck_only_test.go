package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Тестируем на xml-файлах, в которых условный партнёр прислал тег baggageRecheck, но не прислал тег virtualInterline.
//Тег baggageRecheck отвечает за флаг recheck_baggage в массиве flight_legs в ответе дельты,
//тег virtualInterline отвечает за флаг is_virtual_interline в массиве transfer_terms в ответе дельты.

//Так как в тестовом ответе партнёра нет тега virtualInterline, мы в ответе дельты должны расставить флаги is_virtual_interline
//в том же порядке, в котором они определены тегами baggageRecheck.

//Что должны проверить тесты:
//
//- что в формате дельты флаги recheck_baggage в массиве flight_legs расставляются в правильном порядке,
//  при двух разных значениях ключа recheck_baggage_after в конфиге (true|false);
//- что в формате дельты флаги is_virtual_interline в массиве transfer_terms расставляются в том же порядке,
//  что флаги recheck_baggage в массиве flight_legs, также при двух разных значениях ключа recheck_baggage_after в конфиге (true|false).

//---------------------------

//recheckBaggageAfter = false
//virtualInterlineAfter = true

//Такую комбинацию проверять смысла нет, т.к. такая расстановка флагов нелогична.
//Мы ожидаем, что партнёры будут расставлять теги baggageRecheck и virtualInterline парами, вместе.
//Проще говоря, если в конкретном флайте один тег true, то второй тоже true.

//--------------------------

//-------------------------
//recheckBaggageAfter = false
//virtualInterlineAfter = false

//Эта комбинация флагов - идеальная конфигурация, которую мы хотим добиться от партнёров (и после этого отрефакторить все эти костыли)

//---------------------------

//# Case 1
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: false}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: false}]]

func TestParseTwoFlightsWithoutRBAndFlagsFalseFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/false-false.xml", false, false)
	assert.Equal(t, false, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)

	assert.Equal(t, false, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, 1, len(transferTerms[0]))
}

//# Case 2
//Перелеты Партнера: [{RecheckBaggage: true}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: true}]]

func TestParseTwoFlightsWithRBAndFlagsFalseFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/true-false.xml", false, false)
	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, 1, len(transferTerms[0]))
}

//# Case 3
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: false}, {IsVirtualInterline: true}]]

func TestParseThreeFlightsWithRBAndFlagsFalseFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/false-true-false.xml", false, false)
	assert.Equal(t, false, legs[0].RecheckBaggage)
	assert.Equal(t, true, legs[1].RecheckBaggage)
	assert.Equal(t, false, legs[2].RecheckBaggage)

	assert.Equal(t, false, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, true, transferTerms[0][1].IsVirtualInterline)
	assert.Equal(t, 2, len(transferTerms[0]))
}

//# Case 4
//Перелеты Партнера: [{RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: true}, {IsVirtualInterline: false}, {IsVirtualInterline: true}]]

func TestParseFourFlightsWithRBAndFlagsFalseFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/true-false-true-false.xml", false, false)
	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)
	assert.Equal(t, true, legs[2].RecheckBaggage)
	assert.Equal(t, false, legs[3].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, false, transferTerms[0][1].IsVirtualInterline)
	assert.Equal(t, true, transferTerms[0][2].IsVirtualInterline)

	assert.Equal(t, 3, len(transferTerms[0]))
}

//recheckBaggageAfter = true
//virtualInterlineAfter = false

//Для этой пачки тестов включим в условном конфиге интеграции только ключ recheck_baggage_after, но не virtual_interline_after.
//Тестируем такой кейс, потому что сейчас так работают партнёры, которые не отдают нам тег virtualInterline,
//и хочется понимать, есть ли подводные камни.

//--------------------------

//# Case 1
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: false}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: false}]]

func TestParseTwoFlightsWithoutRBAndFlagsTrueFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/false-false.xml", true, false)
	assert.Equal(t, false, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)

	assert.Equal(t, false, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, 1, len(transferTerms[0]))
}

//# Case 2
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}]
//Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: true}]]

func TestParseTwoFlightsWithRBAndFlagsTrueFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/false-true.xml", true, false)
	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, 1, len(transferTerms[0]))
}

//# Case 3
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: true}, {IsVirtualInterline: false}]]

func TestParseThreeFlightsWithRBAndFlagsTrueFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/false-true-false.xml", true, false)
	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)
	assert.Equal(t, false, legs[2].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, false, transferTerms[0][1].IsVirtualInterline)
	assert.Equal(t, 2, len(transferTerms[0]))
}

//# Case 4
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: true}]
//Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: true}, {IsVirtualInterline: false}, {IsVirtualInterline: true}]]

func TestParseFourFlightsWithoutRBAndFlagsTrueFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/false-true-false-true.xml", true, false)
	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)
	assert.Equal(t, true, legs[2].RecheckBaggage)
	assert.Equal(t, false, legs[3].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, false, transferTerms[0][1].IsVirtualInterline)
	assert.Equal(t, true, transferTerms[0][2].IsVirtualInterline)

	assert.Equal(t, 3, len(transferTerms[0]))
}

//--------------------------

//recheckBaggageAfter = true
//virtualInterlineAfter = true

//Для этой пачки тестов включим в условном конфиге интеграции и ключ recheck_baggage_after, и ключ virtual_interline_after.
//Условный партнёр в этом тесте не отдаёт нам теги virtualInterline,
//но хочется понимать, что будет, если мы включим ключ virtual_interline_after.

//--------------------------

//# Case 1
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: false}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: false}]]

func TestParseTwoFlightsWithoutRBAndFlagsTrueTrue(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/false-false.xml", true, true)

	assert.Equal(t, false, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)

	assert.Equal(t, false, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, 1, len(transferTerms[0]))
}

//# Case 2
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}]
//Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: true}]]

func TestParseTwoFlightsWithRBAndFlagsTrueTrue(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/false-true.xml", true, true)

	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, 1, len(transferTerms[0]))
}

//# Case 3
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: true}, {IsVirtualInterline: false}]]

func TestParseThreeFlightsWithRBAndFlagsTrueTrue(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/false-true-false.xml", true, true)

	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)
	assert.Equal(t, false, legs[2].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, false, transferTerms[0][1].IsVirtualInterline)
	assert.Equal(t, 2, len(transferTerms[0]))
}

//# Case 4
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: true}]
//Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: true}, {IsVirtualInterline: false}, {IsVirtualInterline: true}]]

func TestParseFourFlightsWithRBAndFlagsTrueTrue(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_rb/false-true-false-true.xml", true, true)

	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)
	assert.Equal(t, true, legs[2].RecheckBaggage)
	assert.Equal(t, false, legs[3].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, false, transferTerms[0][1].IsVirtualInterline)
	assert.Equal(t, true, transferTerms[0][2].IsVirtualInterline)

	assert.Equal(t, 3, len(transferTerms[0]))
}
