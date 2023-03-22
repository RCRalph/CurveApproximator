import csv

class Distributor:
    def __init__(self):
        self.values = []
        self.distributions = []
        self.materials = 0
        self.precision = 0

    def set_data(self, filename):
        with open(filename) as file:
            reader = csv.reader(file)
            for i in reader:
                self.values.append([float(j.replace(",", ".")) for j in i])

        if len(self.values):
            self.materials = len(self.values[0]) - 1

        return True

    def set_precision(self, value):
        self.precision = int(1 / value)

    def generate_distributions(self):
        for i in range(self.precision ** self.materials):
            dist = []
            for _ in range(self.materials):
                dist.append(i % self.precision)
                i //= self.precision

            if sum(dist) == self.precision:
                self.distributions.append(dist)

    def calculate_deviation(self, dist):
        result = 0

        for i in self.values:
            dist_values = [(coefficient * value) for coefficient, value in zip(dist, i[1:])]
            result += abs(i[0] - sum(dist_values) / self.precision)

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
    print("Successfully gathered data.")
else:
    exit(2)

dist.set_precision(float(input("Enter precision (decimal): ")))

print(dist.get_best_distribution())
