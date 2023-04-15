# Curve Approximator
This tool helps with approximating a target curve from given datasets with a given precision.

## Features
- Written in **Go**.
- **Multithreading** support.
- Around **1000x** faster than an equivalent Python 3.10 script, measured on AMD Ryzen 7 5800X.
- CLI interface (web version in development!)

## Getting started
1. Fill in your data in the **Chart data** sheet of `Data.ods` spreadsheet (you can find it in the `examples/cli` folder).
2. Select all content of the **Data for analysis** sheet and paste it into `Data.csv` file. Optionally the decimal point can be represented by a comma. Make sure to save the file.
3. Install required dependencies using `go mod vendor`.
4. Run the software without compiling using `go run cmd/cli/main.go`. After the calculation is done, the result will appear in the command line.
5. Fill in the calculated values into the **Chart data** sheet in the row starting with **Optimal**.
6. Optionally: Export the results to `.pdf` by printing the spreadsheet.

## Compilation
To compile the software, run `go build -o bin/cli cmd/cli/main.go`. After that, the resulting executable will be available under `./bin/cli`.

## Flags
- `-file` - path to the file with data, the default is `-file examples/cli/Data.csv`.
- `-precision` - number of decimal places to calculate the result to, the default is `-precision 1`, which will calculate the optimal curve to a precision of 10%.
- `-separator` - decimal point separator used in data, the default is `-separator "."`.
- `-delimiter` - field delimiter used in data (one character only), the default is `-delimiter ","`. You can also specify the `\t` character, but no other escaped characters work.

Flags can be prefixed either with `-` and `--`. Bear in mind that you don't have to set the flag if the default value is what you want.

## Examples of usage:
```bash
# Run without compilation
go run cmd/cli/main.go -file Example/Data.csv -precision 3 -separator "." -delimiter ","

# Run compiled executable
./bin/cli --file Data.csv --precision 3 --separator "," --delimiter "\t"
```

## Expected execution time
The time complexity of the used algorithm is `O(10^(n*k))`, where `n` is the number of datasets and `k` is the given precision. This means that the expected execution time for 3 datasets and the precision of 3 decimal places is around 1 second (tested on AMD Ryzen 7 5800x). Setting the precision to 4 will result in execution for around 1000 seconds (around 17 minutes).
