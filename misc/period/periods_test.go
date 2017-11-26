package period_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/outo/go-shared/misc/period"
	"time"
)

var _ = Describe("Periods", func() {

	const timeFormat = "150405"

	timeZero, _ := time.Parse(timeFormat, "090000.000")

	convertStringToPeriods := func(timeZero time.Time, s string) Periods {
		var ps []Period
		var startTime, endTime time.Time
		t := timeZero
		periodStarted := false
		for pos := 0; pos < len(s); pos++ {

			if s[pos] == '1' {
				if periodStarted {
				} else {
					startTime = t
				}
				periodStarted = true
			} else {
				if periodStarted {
					p, _ := CreatePeriod(startTime, endTime)
					ps = append(ps, p)
				}
				periodStarted = false
			}
			t = t.Add(time.Minute)
			endTime = t
		}
		if periodStarted {
			p, _ := CreatePeriod(startTime, endTime)
			ps = append(ps, p)
		}
		return CreatePeriods(ps)
	}

	Describe("Helper function converting string to periods", func() {
		Context("Given time zero as 9am", func() {

			newPeriod := func(startIncl, endExcl string) (period Period) {
				period, err := CreatePeriod(ParseShortTimePanicOnError(startIncl), ParseShortTimePanicOnError(endExcl))
				if err != nil {
					panic("helper function newPeriod failed")
				}
				return
			}

			DescribeTable("Periods in string notation",
				func(periodsAsString string, expectedPeriods Periods) {
					actualPeriods := convertStringToPeriods(timeZero, periodsAsString)
					Expect(actualPeriods).To(Equal(expectedPeriods))
				},
				Entry("no periods", "0", CreatePeriods([]Period{
				})),
				Entry("no periods, longer input", "0000000000", CreatePeriods([]Period{
				})),
				Entry("single short period", "1", CreatePeriods([]Period{
					newPeriod("090000", "090100"),
				})),
				Entry("single long period", "1111111111", CreatePeriods([]Period{
					newPeriod("090000", "091000"),
				})),
				Entry("single period surrounded by zeroes", "0001111000", CreatePeriods([]Period{
					newPeriod("090300", "090700"),
				})),
				Entry("multiple periods a", "110011", CreatePeriods([]Period{
					newPeriod("090000", "090200"),
					newPeriod("090400", "090600"),
				})),
				Entry("multiple periods b", "0111100111", CreatePeriods([]Period{
					newPeriod("090100", "090500"),
					newPeriod("090700", "091000"),
				})),
				Entry("multiple periods c", "1111001100", CreatePeriods([]Period{
					newPeriod("090000", "090400"),
					newPeriod("090600", "090800"),
				})),
				Entry("multiple periods d", "1011001101", CreatePeriods([]Period{
					newPeriod("090000", "090100"),
					newPeriod("090200", "090400"),
					newPeriod("090600", "090800"),
					newPeriod("090900", "091000"),
				})),
			)
		})
	})

	Describe("Given input parameters", func() {
		period1, err := CreatePeriod(time.Now(), time.Now().Add(time.Second))
		period2, err := CreatePeriod(time.Now(), time.Now().Add(time.Minute))
		period3, err := CreatePeriod(time.Now(), time.Now().Add(time.Hour))
		expectedPeriodsAsSlice := []Period{
			period1, period2, period3,
		}
		periods := CreatePeriods(expectedPeriodsAsSlice)

		It("Will not fail", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("Will return them upon call to Periods.AsSlice()", func() {
			actualPeriodsAsSlice := periods.AsSlice()
			Expect(actualPeriodsAsSlice).To(BeEquivalentTo(expectedPeriodsAsSlice))
		})
	})

	Describe("Given parametrised tests", func() {

		tableEntries := []TableEntry{
			Entry("disparate and higher", "dummy",
				"0000001111000000", //periodsA
				"0000000000001111", //periodsB
				"0000000000001111", //expectedPeriodsBSubtractA
				"0000001111001111"), //expectedPeriodsAUnionB
			Entry("adjacent and higher", "dummy",
				"0000001111000000",
				"0000000000111111",
				"0000000000111111",
				"0000001111111111"),
			Entry("overlapping upper end", "dummy",
				"0000001111000000",
				"0000000001111111",
				"0000000000111111",
				"0000001111111111"),
			Entry("containing", "dummy",
				"0000001111000000",
				"0000001111111111",
				"0000000000111111",
				"0000001111111111"),
			Entry("containing 2", "dummy",
				"0000001111000000",
				"0000011111111111",
				"0000010000111111",
				"0000011111111111"),
			Entry("containing 3", "dummy",
				"0000001111000000",
				"0000011111000000",
				"0000010000000000",
				"0000011111000000"),
			Entry("overlapping lower end", "dummy",
				"0000001111000000",
				"0000011110000000",
				"0000010000000000",
				"0000011111000000"),
			Entry("overlapping lower end", "dummy",
				"0000001111000000",
				"0001111000000000",
				"0001110000000000",
				"0001111111000000"),
			Entry("adjacent and lower", "dummy",
				"0000001111000000",
				"0001110000000000",
				"0001110000000000",
				"0001111111000000"),
			Entry("disparate and lower", "dummy",
				"0000001111000000",
				"0011100000000000",
				"0011100000000000",
				"0011101111000000"),
			Entry("containing 4", "dummy",
				"0000001111000000",
				"0011111111111000",
				"0011110000111000",
				"0011111111111000"),
			Entry("contained", "dummy",
				"0000001111000000",
				"0000000110000000",
				"0000000000000000",
				"0000001111000000"),
			Entry("second empty", "dummy",
				"0000001111000000",
				"0000000000000000",
				"0000000000000000",
				"0000001111000000"),
			Entry("first empty", "dummy",
				"0000000000000000",
				"0000001111000000",
				"0000001111000000",
				"0000001111000000"),
			//multiple periods
			Entry("multiple periods 1", "dummy",
				"1111111111111111",
				"0000001000110000",
				"0000000000000000",
				"1111111111111111"),
			Entry("multiple periods 2", "dummy",
				"0011001111000000",
				"0001111110000000",
				"0000110000000000",
				"0011111111000000"),
			Entry("multiple periods 3", "dummy",
				"0011001111000000",
				"0111111111111100",
				"0100110000111100",
				"0111111111111100"),
			Entry("multiple periods 4", "dummy",
				"1011001111011100",
				"0110101001110010",
				"0100100000100010",
				"1111101111111110"),
		}

		DescribeTable("periods union",
			func(dummyForFormattingOnly, periodsA, periodsB, expectedPeriodsBSubtractA, expectedPeriodsAUnionB string) {
				startPeriods := convertStringToPeriods(timeZero, periodsA)
				requestedPeriods := convertStringToPeriods(timeZero, periodsB)
				expectedEndPeriods := convertStringToPeriods(timeZero, expectedPeriodsAUnionB)

				actualEndPeriods, err := requestedPeriods.Union(startPeriods)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(actualEndPeriods).To(Equal(expectedEndPeriods))
			}, tableEntries...,
		)

		DescribeTable("periods subtract",
			func(dummyForFormattingOnly, periodsA, periodsB, expectedPeriodsBSubtractA, expectedPeriodsAUnionB string) {
				startPeriods := convertStringToPeriods(timeZero, periodsA)
				requestedPeriods := convertStringToPeriods(timeZero, periodsB)
				expectedApprovedPeriods := convertStringToPeriods(timeZero, expectedPeriodsBSubtractA)

				actualApprovedPeriods, err := requestedPeriods.Subtract(startPeriods)
				Expect(err).ShouldNot(HaveOccurred())

				Expect(actualApprovedPeriods).To(Equal(expectedApprovedPeriods))
			}, tableEntries...,
		)
	})

	Measure("the benchmark performance of Periods.Union()", func(b Benchmarker) {

		periodsA := convertStringToPeriods(timeZero,
			"000101111010110101111100010110101011011010010110001010001001010100100001111111010010010010100101110101010111")
		periodsB := convertStringToPeriods(timeZero,
			"11011100110101001010101001000101010010111101010101100010101011101011001001010100000001111101010010110")
		expectedPeriodsAUnionB := convertStringToPeriods(timeZero,
			"110111111111110111111110010111111111111111010111011010101011111110110011111111010010011111110101111101010111")

		runtime := b.Time("Periods.Union()", func() {
			actualPeriodsAUnionB, err := periodsA.Union(periodsB)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(actualPeriodsAUnionB).To(Equal(expectedPeriodsAUnionB))
		})

		Expect(runtime.Nanoseconds()).Should(BeNumerically("<", (100 * time.Millisecond).Nanoseconds()))

		b.RecordValue("Execution time in microseconds", float64(runtime.Nanoseconds()/1000))
	},
		1000)

	Measure("the benchmark performance of Periods.Subtract()", func(b Benchmarker) {


		periodsA := convertStringToPeriods(timeZero,
			"000101111010110101111100010110101011011010010110001010001001010100100001111111010010010010100101110101010111")
		periodsB := convertStringToPeriods(timeZero,
			"11011100110101001010101001000101010010111101010101100010101011101011001001010100000001111101010010110")
		expectedPeriodsASubtractB := convertStringToPeriods(timeZero,
			"000000110010100101010100000110101011010000000010000010000001000100000001101010010010000000100001010001010111")

		runtime := b.Time("Periods.Subtract()", func() {
			actualPeriodsASubtractB, err := periodsA.Subtract(periodsB)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(actualPeriodsASubtractB).To(Equal(expectedPeriodsASubtractB))
		})

		Expect(runtime.Nanoseconds()).Should(BeNumerically("<", (100 * time.Millisecond).Nanoseconds()))

		b.RecordValue("Execution time in microseconds", float64(runtime.Nanoseconds()/1000))
	},
		1000)

})
