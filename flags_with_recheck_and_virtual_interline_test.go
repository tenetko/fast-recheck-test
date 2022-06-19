package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Тестируем на xml-файлах, в которых условный партнёр прислал и тег baggageRecheck, и тег virtualInterline.
//Тег baggageRecheck отвечает за флаг recheck_baggage в массиве flight_legs в ответе дельты,
//тег virtualInterline отвечает за флаг is_virtual_interline в массиве transfer_terms в ответе дельты.

//Мы в ответе дельты должны расставить флаги recheck_baggage и is_virtual_interline в одинаковом порядке.

//Что должны проверить тесты:
//
//- что в формате дельты флаги recheck_baggage в массиве flight_legs расставляются в правильном порядке,
//  при двух разных значениях ключа recheck_baggage_after в конфиге (true|false);
//- что в формате дельты флаги is_virtual_interline в массиве transfer_terms расставляются в том же порядке,
//  что флаги recheck_baggage в массиве flight_legs, также при двух разных значениях ключа recheck_baggage_after в конфиге (true|false).

//---------------------------

//recheckBaggageAfter = false
//virtualInterlineAfter = true

//recheckBaggageAfter = true
//virtualInterlineAfter = false

//Такие комбинации проверять смысла нет, т.к. такая расстановка флагов нелогична.
//Мы ожидаем, что партнёры будут расставлять теги baggageRecheck и virtualInterline парами, вместе.
//Проще говоря, если в конкретном флайте один тег true, то второй тоже true.

//-------------------------
//recheckBaggageAfter = false
//virtualInterlineAfter = false

//Эта комбинация флагов - идеальная конфигурация, которую мы хотим добиться от партнёров (и после этого отрефакторить все эти костыли)

//--------------------------

//# Case 1
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: false}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: false}]]

func TestParseTwoFlightsWithoutRBAndVIAndFlagsFalseFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_vi_rb/false-false.xml", false, false)
	assert.Equal(t, false, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)

	assert.Equal(t, false, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, 1, len(transferTerms[0]))
}

//# Case 2
//Перелеты Партнера: [{RecheckBaggage: true}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: true}]]

func TestParseTwoFlightsWithRBAndVIAndFlagsFalseFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_vi_rb/true-false.xml", false, false)
	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, 1, len(transferTerms[0]))
}

//# Case 3
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: false}, {IsVirtualInterline: true}]]

func TestParseThreeFlightsWithRBAndVIAndFlagsFalseFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_vi_rb/false-true-false.xml", false, false)
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

func TestParseFourFlightsWithRBAndVIAndFlagsFalseFalse(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_vi_rb/true-false-true-false.xml", false, false)
	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)
	assert.Equal(t, true, legs[2].RecheckBaggage)
	assert.Equal(t, false, legs[3].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, false, transferTerms[0][1].IsVirtualInterline)
	assert.Equal(t, true, transferTerms[0][2].IsVirtualInterline)

	assert.Equal(t, 3, len(transferTerms[0]))
}

//-------------------------
//recheckBaggageAfter = true
//virtualInterlineAfter = true

//Для этой пачки тестов включим в условном конфиге интеграции и ключ recheck_baggage_after, и ключ virtual_interline_after.
//Предполагается, что это конфигурация на время переезда для партнёра, который отдаёт оба тега.

//--------------------------

//# Case 1
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: false}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: false}]]

func TestParseTwoFlightsWithoutRBAndVIAndFlagsTrueTrue(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_vi_rb/false-false.xml", true, true)

	assert.Equal(t, false, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)

	assert.Equal(t, false, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, 1, len(transferTerms[0]))
}

//# Case 2
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}]
//Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: true}]]

func TestParseTwoFlightsWithRBAndVIAndFlagsTrueTrue(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_vi_rb/false-true.xml", true, true)

	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, 1, len(transferTerms[0]))
}

//# Case 3
//Перелеты Партнера: [{RecheckBaggage: false}, {RecheckBaggage: true}, {RecheckBaggage: false}]
//Перелеты в Дельте: [{RecheckBaggage: true}, {RecheckBaggage: false}, {RecheckBaggage: false}]
//Условия пересадки в transferTerms: [[{IsVirtualInterline: true}, {IsVirtualInterline: false}]]

func TestParseThreeFlightsWithRBAndVIAndFlagsTrueTrue(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_vi_rb/false-true-false.xml", true, true)

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

func TestParseFourFlightsWithRBAndVIAndFlagsTrueTrue(t *testing.T) {
	// второй и третий аргументы Parse - ключи конфига recheckBaggageAfter и virtualInterlineAfter
	legs, transferTerms := Parse("xml_vi_rb/false-true-false-true.xml", true, true)

	assert.Equal(t, true, legs[0].RecheckBaggage)
	assert.Equal(t, false, legs[1].RecheckBaggage)
	assert.Equal(t, true, legs[2].RecheckBaggage)
	assert.Equal(t, false, legs[3].RecheckBaggage)

	assert.Equal(t, true, transferTerms[0][0].IsVirtualInterline)
	assert.Equal(t, false, transferTerms[0][1].IsVirtualInterline)
	assert.Equal(t, true, transferTerms[0][2].IsVirtualInterline)

	assert.Equal(t, 3, len(transferTerms[0]))
}
