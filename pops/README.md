info for package "github.com/outo/temporal/pops"

pops - period(s) operations

Contains constructs:
* `Period`:

    an immutable representation of time range with the lower boundary being inclusive and upper being exclusive
    
    methods:
    * `Check` will identify the placement of specified time range with respect to the time range within this period
    * `Get` copy of startIncl and endIncl as multi return
    * `GetStartIncl` copy of startIncl
    * `GetEndExcl` copy of endIncl
    
     *Performance*
     
     each test performed on on an Intel® Core™ i7 4500U Processor with enough DDR3 SDRAM, 
     consisted of 100 samples, 
     each test based on checking the relationship between a target time range and a sliding and size-changing time range constructed to simulate each of the possible relationships 
     took on average nearly 100 nanoseconds per PeriodCheck (with extremums at 80 and 160 nanoseconds)
     
* `Periods`:

    an immutable representation of practically (memory, language constraints) unlimited number of periods

    methods:
    * `Subtract` will return a new `Periods` instance representing the time ranges within this object that are not within specified object
    * `Union` will return a new `Periods` instance representing the time ranges within this object and specified object
    
    both methods' result is:
    * defragmented - there aren't going to be two ranges that are overlapping, 
    * stitched - there aren't going to be two ranges that are adjacent,
    * sorted - the ranges are appearing from the earliest to most recent 
   
    *Note*:
     
     It is possible to create an instance of `Periods` with erratic time ranges 
     by passing such to the constructing function `CreatePeriods` 
     
     *Performance*
     
     each test performed on on an Intel® Core™ i7 4500U Processor with enough DDR3 SDRAM, 
     consisted of 1000 samples of two Periods, 
     each with about 30 distinct ranges (defragmented, stitched and sorted) 
     took on average:
     * Subtract - 80 microseconds (producing nearly 30 ranges)
     * Union - 76 microseconds (producing about 20 ranges)
     

     