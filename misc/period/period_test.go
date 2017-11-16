package period

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo/extensions/table"
	"time"
)

var _ = Describe("Period", func() {

	const timeFormat = "150405.000"

	var (
		//order in the alphabet implies chronological order
		a, _ = time.Parse(timeFormat, "090000.000")
		b, _ = time.Parse(timeFormat, "100000.000")
		c, _ = time.Parse(timeFormat, "110000.000")
		d, _ = time.Parse(timeFormat, "120000.000")
		e, _ = time.Parse(timeFormat, "130000.000")
		f, _ = time.Parse(timeFormat, "140000.000")
		g, _ = time.Parse(timeFormat, "150000.000")
		h, _ = time.Parse(timeFormat, "160000.000")
		i, _ = time.Parse(timeFormat, "170000.000")
		j, _ = time.Parse(timeFormat, "180000.000")
		k, _ = time.Parse(timeFormat, "190000.000")
		l, _ = time.Parse(timeFormat, "200000.000")
		t    = []time.Time{a, b, c, d, e, f, g, h, i, j, k, l}
	)

	Describe("Given period set as d,i", func() {

		var period, _ = CreatePeriod(t[3], t[8])

		DescribeTable("When relationship check is done",
			func(period Period, startIncl, endExcl time.Time, expectedRelationship Relationship) {
				actualRelationship := period.Check(startIncl, endExcl)
				Expect(actualRelationship).To(Equal(expectedRelationship))
			},
			Entry("Will identify period as disparate and lower", period, t[0], t[1], DisparateAndLower),
			Entry("Will identify period is adjacent and lower", period, t[0], t[3], AdjacentAndLower),
			Entry("Will identify period as overlapping lower end", period, t[2], t[4], OverlappingLowerEnd),
			Entry("Will identify period as contained", period, t[5], t[7], Contained),
			Entry("Will identify period as same", period, t[3], t[8], Same),
			Entry("Will identify period as overlapping upper end", period, t[6], t[10], OverlappingUpperEnd),
			Entry("Will identify period as adjacent and higher", period, t[8], t[11], AdjacentAndHigher),
			Entry("Will identify period as disparate and higher", period, t[9], t[11], DisparateAndHigher),
			Entry("Will identify period as containing", period, t[0], t[11], Containing),
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

	Measure("the benchmark performance of Period.Check()", func(b Benchmarker) {

		const (
			timestampFormat    = "20060102 150405.000"
			nearlyASecondPrime = 999727999 * time.Nanosecond
		)
		periodStart, err := time.Parse(timestampFormat, "20150318 095214.522")
		Expect(err).ShouldNot(HaveOccurred())

		periodEnd := periodStart.Add(300 * nearlyASecondPrime)

		period, err := CreatePeriod(periodStart, periodEnd)
		Expect(err).ShouldNot(HaveOccurred())

		benchmarkStart := periodStart.Add(-120 * nearlyASecondPrime)
		benchmarkEnd := periodEnd.Add(120 * nearlyASecondPrime)

		counters := make(map[int]int, 10)
		runtime := b.Time("Period.Check() with scan", func() {
			for testStartTime := benchmarkStart; testStartTime.Before(benchmarkEnd); testStartTime = testStartTime.Add(nearlyASecondPrime) {
				for testEndTime := testStartTime; testEndTime.Before(benchmarkEnd); testEndTime = testEndTime.Add(nearlyASecondPrime) {
					periodRelationship := period.Check(testStartTime, testEndTime)
					counters[int(periodRelationship)]++
				}
			}
		})

		Expect(runtime.Nanoseconds()).Should(BeNumerically("<", (500 * time.Millisecond).Nanoseconds()))

		count := 0
		for _, counter := range counters {
			count += counter
		}
		b.RecordValue("Total invocations", float64(count))
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
