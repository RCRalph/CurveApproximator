# Curve Approximator
This tool helps with approximating a target curve from given datasets with a given precision.

## Features
- Written in Go.
- Multithreading support.
- Around 1000x faster than an equivalent Python 3.10 script, measured on AMD Ryzen 7 5800X.

## Getting started
1. Fill in your data in the **Chart data** sheet of `Data.ods` spreadsheet (you can find it in the **Example** folder).
2. Select all content of the **Data for analysis** sheet and paste it into `Data.csv` file. Make sure the delimiter is the tabulation character (`\t`). Optionally the decimal point can be represented by a comma. Make sure to save the file.
3. Run the software without compiling using `go run curve-approximator.go`. After the calculation is done, the result will appear in the command line.
4. Fill in the calculated values into the **Chart data** sheet in the row starting with **Optimal**.
5. Optionally: Export the results to `.pdf` by printing the spreadsheet.

## Compilation
To compile the software, run `go build curve-approximator.go`. After that, the resulting executable will be available under `./curve-approximator`.

## Flags
- `-file` - path to the file with data, ex. `-file Example/Data.csv`.
- `-precision` - precision to calculate the result to, ex. `-precision 0.01` will calculate the optimal curve to a precision of **1%**.

Flags can be prefixed either with `-` and `--`.

## Examples of usage:
```sh
# Run without compilation
go run curve-approximator.go -file Example/Data.csv -precision 0.001

# Run compiled executable
./curve-approximator --file Data.csv --precision 0.001
```
