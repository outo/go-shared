package pops_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo/extensions/table"
	"time"
	p "github.com/outo/temporal/pops"
)

var _ = Describe("Period", func() {

	now := time.Now()

	tm := func(point int) time.Time {
		return now.Add(time.Duration(point) * time.Hour)
	}

	Describe("Given period from 4 to 10", func() {
		period, err := p.NewPeriod(tm(4), tm(10))

		It("Will not fail", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("Will return time point 4 upon calling Period.GetStartIncl()", func() {
			Expect(period.GetStartIncl()).To(Equal(tm(4)))
		})
		It("Will return time point 10 upon calling Period.GetEndExcl()", func() {
			Expect(period.GetEndExcl()).To(Equal(tm(10)))
		})
		It("Will return time points 4 and 10 upon calling Period.Get()", func() {
			p1, p2 := period.Get()
			Expect(p1).To(Equal(tm(4)))
			Expect(p2).To(Equal(tm(10)))
		})
	})

	Describe("Given period from 3 (incl) to 8 (excl)", func() {

		var period, _ = p.NewPeriod(tm(3), tm(8))

		DescribeTable("When relationship check is done",
			func(p p.Period, startIncl, endExcl time.Time, expected p.Relationship) {
				actualRelationship := p.Check(startIncl, endExcl)
				Expect(actualRelationship).To(Equal(expected))
			},
			Entry("Range from 0 to 1 is disparate and lower", period,
				tm(0), tm(1), p.DisparateAndLower),
			Entry("Range from 0 to 3 is adjacent and lower", period,
				tm(0), tm(3), p.AdjacentAndLower),
			Entry("Range from 2 to 4 is overlapping lower end", period,
				tm(2), tm(4), p.OverlappingLowerEnd),
			Entry("Range from 5 to 7 is contained", period,
				tm(5), tm(7), p.Contained),
			Entry("Range from 3 to 8 is same", period,
				tm(3), tm(8), p.Same),
			Entry("Range from 6 to 10 is overlapping upper end", period,
				tm(6), tm(10), p.OverlappingUpperEnd),
			Entry("Range from 8 to 11 is adjacent and higher", period,
				tm(8), tm(11), p.AdjacentAndHigher),
			Entry("Range from 9 to 11 is disparate and higher", period,
				tm(9), tm(11), p.DisparateAndHigher),
			Entry("Range from 0 to 11 is containing", period,
				tm(0), tm(11), p.Containing),
		)
	})

	Describe("Failure scenarios", func() {
		It("Will fail to create a period where startIncl is equal endExcl", func() {
			now := time.Now()
			_, err := p.NewPeriod(now, now)
			Expect(err).Should(HaveOccurred())
		})

		It("Will fail to create a period where startIncl is after endExcl", func() {
			now := time.Now()
			_, err := p.NewPeriod(now.Add(time.Nanosecond), now)
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

		period, err := p.NewPeriod(periodStart, periodEnd)
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
		b.RecordValue("Unknown", float64(counters[int(p.Unknown)]))
		b.RecordValue("DisparateAndLower", float64(counters[int(p.DisparateAndLower)]))
		b.RecordValue("AdjacentAndLower", float64(counters[int(p.AdjacentAndLower)]))
		b.RecordValue("OverlappingLowerEnd", float64(counters[int(p.OverlappingLowerEnd)]))
		b.RecordValue("Contained", float64(counters[int(p.Contained)]))
		b.RecordValue("Same", float64(counters[int(p.Same)]))
		b.RecordValue("OverlappingUpperEnd", float64(counters[int(p.OverlappingUpperEnd)]))
		b.RecordValue("AdjacentAndHigher", float64(counters[int(p.AdjacentAndHigher)]))
		b.RecordValue("DisparateAndHigher", float64(counters[int(p.DisparateAndHigher)]))
		b.RecordValue("Containing", float64(counters[int(p.Containing)]))
	},
		100)
})

