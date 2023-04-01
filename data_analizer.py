import csv

class Distributor:
    def __init__(self):
        self.values = []
        self.target = []
        self.weighs = []
        self.distributions = []
        self.precision = 0

    def set_data(self, filename):
        values = []
        with open(filename) as file:
            reader = csv.reader(file, delimiter="\t")

            for i, row in enumerate(reader):
                if not i:
                    values = [list() for _ in range(len(row))]

                value = [float(j.replace(",", ".")) for j in row]
                self.target.append(value[0])
                self.weighs.append(value[1])

                for j, item in enumerate(value[2:]):
                    values[j].append(item)

        for i in values:
            self.values.append([])
            sum_value = 0
            for j in reversed(range(len(i))):
                self.values[-1].append(sum_value)
                sum_value += i[j]

            if not sum_value:
                self.values.pop()
            else:
                self.values[-1] = list(reversed(self.values[-1]))

        return True

    def set_precision(self, value):
        self.precision = int(1 / value)

    def generate_distributions(self):
        for i in range(self.precision ** len(self.values)):
            dist = []
            for _ in range(len(self.values)):
                dist.append(i % self.precision)
                i //= self.precision

            if sum(dist) == self.precision:
                self.distributions.append(dist)

    def calculate_deviation(self, dist):
        result = 0

        for i in range(len(self.target)):
            dist_values = [(coefficient * value[i]) for coefficient, value in zip(dist, self.values)]
            result += abs(self.target[i] - sum(dist_values) / self.precision) * self.weighs[i]

        return result / len(self.values)

    def get_best_distribution(self):
        self.generate_distributions()

        best_dist, best_value = self.distributions[0], self.calculate_deviation(self.distributions[0])
        for i in self.distributions:
            dist_deviation = self.calculate_deviation(i)

            if dist_deviation < best_value:
                best_dist, best_value = i, dist_deviation

        return [i / self.precision for i in best_dist]

dist = Distributor()
if dist.set_data(input("Enter file name: ")):
    print(f"Successfully gathered {len(dist.values)} datasets.")
else:
    exit(2)

dist.set_precision(float(input("Enter precision (decimal): ")))

print(dist.get_best_distribution())
