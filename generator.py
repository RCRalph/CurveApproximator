import random
import csv

def generate_numbers(n):
    # Initialize an empty list to store the random numbers
    nums = []
    # Generate n-1 random numbers between 0 and 1
    for _ in range(n-1):
        nums.append(random.uniform(0, 1))
    # Sort the numbers in ascending order
    nums.sort()
    # Calculate the differences between adjacent numbers
    diffs = [nums[0]] + [nums[i+1]-nums[i] for i in range(n-2)] + [1-nums[-1]]
    # Return the differences as the n random numbers
    return diffs

set_count, column_count = 3, 40
# Generate four sets of column_count numbers between 0 and 1 that include 0 and 1
sets = []
for i in range(set_count):
    s = [0, 1]
    while len(s) < column_count:
        x = round(random.uniform(0, 1), 3)
        if x not in s:
            s.append(x)
    sets.append(sorted(s))

coefficients = generate_numbers(set_count)
print(coefficients)

target_set = []
for i in range(column_count):
    target_set.append(round(sum([sets[j][i] * coefficients[j] for j in range(set_count)]), 3))

sets.insert(0, target_set)

# Export the sets to a CSV file
with open('sets.csv', 'w', newline='') as csvfile:
    writer = csv.writer(csvfile)
    for row in zip(*sets):
        writer.writerow(row)
