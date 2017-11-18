package period

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo/extensions/table"
	"time"
)

func ParseTimePanicOnError(timeAsString string) (result time.Time) {
	result, err := time.Parse("150405", timeAsString)
	if err != nil {
		panic("helper function parseTime failed")
	}
	return
}

var _ = Describe("Period", func() {

	var (
		//order in the slice is chronological
		t = []time.Time{
			ParseTimePanicOnError("090000"),
			ParseTimePanicOnError("100000"),
			ParseTimePanicOnError("110000"),
			ParseTimePanicOnError("120000"),
			ParseTimePanicOnError("130000"),
			ParseTimePanicOnError("140000"),
			ParseTimePanicOnError("150000"),
			ParseTimePanicOnError("160000"),
			ParseTimePanicOnError("170000"),
			ParseTimePanicOnError("180000"),
			ParseTimePanicOnError("190000"),
			ParseTimePanicOnError("200000"),
		}
	)

	Describe("Given period from 4 to 10", func() {
		period, err := CreatePeriod(t[4], t[10])

		It("Will not fail", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("Will return time point 4 upon calling Period.GetStartIncl()", func() {
			Expect(period.GetStartIncl()).To(Equal(t[4]))
		})
		It("Will return time point 10 upon calling Period.GetEndExcl()", func() {
			Expect(period.GetEndExcl()).To(Equal(t[10]))
		})
		It("Will return time points 4 and 10 upon calling Period.Get()", func() {
			p1, p2 := period.Get()
			Expect(p1).To(Equal(t[4]))
			Expect(p2).To(Equal(t[10]))
		})
	})

	Describe("Given period from 3 (incl) to 8 (excl)", func() {

		var period, _ = CreatePeriod(t[3], t[8])

		DescribeTable("When relationship check is done",
			func(period Period, startIncl, endExcl time.Time, expectedRelationship Relationship) {
				actualRelationship := period.Check(startIncl, endExcl)
				Expect(actualRelationship).To(Equal(expectedRelationship))
			},
			Entry("Will identify period from 0 to 1 as disparate and lower", period, t[0], t[1], DisparateAndLower),
			Entry("Will identify period from 0 to 3 as adjacent and lower", period, t[0], t[3], AdjacentAndLower),
			Entry("Will identify period from 2 to 4 as overlapping lower end", period, t[2], t[4], OverlappingLowerEnd),
			Entry("Will identify period from 5 to 7 as contained", period, t[5], t[7], Contained),
			Entry("Will identify period from 3 to 8 as same", period, t[3], t[8], Same),
			Entry("Will identify period from 6 to 10 as overlapping upper end", period, t[6], t[10], OverlappingUpperEnd),
			Entry("Will identify period from 8 to 11 as adjacent and higher", period, t[8], t[11], AdjacentAndHigher),
			Entry("Will identify period from 9 to 11 as disparate and higher", period, t[9], t[11], DisparateAndHigher),
			Entry("Will identify period from 0 to 11 as containing", period, t[0], t[11], Containing),
		)

	})

	Describe("Failure scenarios", func() {
		It("Will fail to create a period where startIncl is equal endExcl", func() {
			now := time.Now()
			_, err := CreatePeriod(now, now)
			Expect(err).Should(HaveOccurred())
		})

		It("Will fail to create a period where startIncl is after endExcl", func() {
			now := time.Now()
			_, err := CreatePeriod(now.Add(time.Nanosecond), now)
			Expect(err).Should(HaveOccurred())
		})

	})

	Measure("The benchmark performance of Period.Check()", func(b Benchmarker) {

		const (
			timestampFormat = "20060102 150405.000"
			nearlyASecond   = 987654321 * time.Nanosecond //to exercise the fractions of seconds during tests
		)
		periodStart, err := time.Parse(timestampFormat, "20150318 095214.522")
		Expect(err).ShouldNot(HaveOccurred())

		periodEnd := periodStart.Add(300 * nearlyASecond)

		period, err := CreatePeriod(periodStart, periodEnd)
		Expect(err).ShouldNot(HaveOccurred())

		//benchmark boundary expands sideways from each of the period's boundaries as a way to ensure the adjacency and equality relationship occur
		benchmarkStart := periodStart.Add(-120 * nearlyASecond)
		benchmarkEnd := periodEnd.Add(120 * nearlyASecond)

		counters := make(map[int]int, 10)
		runtime := b.Time("Period.Check() with scan", func() {
			for testStartTime := benchmarkStart; testStartTime.Before(benchmarkEnd); testStartTime = testStartTime.Add(nearlyASecond) {
				for testEndTime := testStartTime; testEndTime.Before(benchmarkEnd); testEndTime = testEndTime.Add(nearlyASecond) {
					periodRelationship := period.Check(testStartTime, testEndTime)
					counters[int(periodRelationship)]++
				}
			}
		})

		Expect(runtime.Nanoseconds()).Should(BeNumerically("<", (500 * time.Millisecond).Nanoseconds()))

		totalInvocations := 0
		for _, counter := range counters {
			totalInvocations += counter
		}
		b.RecordValue("Total invocations", float64(totalInvocations))
		b.RecordValue("Averaged approximation of time (in nanoseconds) it takes to invoke single Period.Check()", float64(int64(runtime.Nanoseconds())/int64(totalInvocations)))
		b.RecordValue("Unknown", float64(counters[int(Unknown)]))
		b.RecordValue("DisparateAndLower", float64(counters[int(DisparateAndLower)]))
		b.RecordValue("AdjacentAndLower", float64(counters[int(AdjacentAndLower)]))
		b.RecordValue("OverlappingLowerEnd", float64(counters[int(OverlappingLowerEnd)]))
		b.RecordValue("Contained", float64(counters[int(Contained)]))
		b.RecordValue("Same", float64(counters[int(Same)]))
		b.RecordValue("OverlappingUpperEnd", float64(counters[int(OverlappingUpperEnd)]))
		b.RecordValue("AdjacentAndHigher", float64(counters[int(AdjacentAndHigher)]))
		b.RecordValue("DisparateAndHigher", float64(counters[int(DisparateAndHigher)]))
		b.RecordValue("Containing", float64(counters[int(Containing)]))
	},
		100)

})
