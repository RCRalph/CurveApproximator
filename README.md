# Curve Approximator
This tool helps with approximating a target curve from given datasets with a given precision.

## Features
- Written in **Go 1.20.2**.
- Multithreading support.
- Around 1000x faster than an equivalent Python 3.10 script, measured on AMD Ryzen 7 5800X.

## Getting started
1. Fill in your data in the **Chart data** sheet of `Data.ods` spreadsheet (you can find it in the **Example** folder).
2. Select all content of the **Data for analysis** sheet and paste it into `Data.csv` file. Make sure the delimiter is the tabulation character (`\t`). Optionally the decimal point can be represented by a comma. Make sure to save the file.
3. Install require dependencies.
4. Run the software without compiling using `go run curve-approximator.go`. After the calculation is done, the result will appear in the command line.
5. Fill in the calculated values into the **Chart data** sheet in the row starting with **Optimal**.
6. Optionally: Export the results to `.pdf` by printing the spreadsheet.

## Compilation
To compile the software, run `go build curve-approximator.go`. After that, the resulting executable will be available under `./curve-approximator`.

## Flags
- `-file` - path to the file with data, the default is `-file Example/Data.csv`.
- `-precision` - precision to calculate the result to, the default is `-precision 0.01`, which will calculate the optimal curve to a precision of **1%**.
- `-separator` - decimal point separator used in data, the default is `-separator "."`.
- `-delimiter` - field delimiter used in data (one character only), the default is `-delimiter ","`. You can also specify the `\t` character, but no other escaped characters work.

Flags can be prefixed either with `-` and `--`. Bear in mind that you don't have to set the flag if the default value is what you want.

## Examples of usage:
```bash
# Run without compilation
go run curve-approximator.go -file Example/Data.csv -precision 0.001 -separator "." -delimiter ","

# Run compiled executable
./curve-approximator --file Data.csv --precision 0.001 --separator "," --delimiter "\t"
```
